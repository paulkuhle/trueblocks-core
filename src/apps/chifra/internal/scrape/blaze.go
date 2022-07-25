package scrapePkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpcClient"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"
)

// ScrapedData combines the block data, trace data, and log data into a single structure
type ScrapedData struct {
	blockNumber int
	traces      rpcClient.Traces
	logs        rpcClient.Logs
}

type BlazeOptions struct {
	Chain         string                     `json:"chain"`
	NChannels     uint64                     `json:"nChannels"`
	NProcessed    uint64                     `json:"nProcessed"`
	StartBlock    uint64                     `json:"startBlock"`
	BlockCount    uint64                     `json:"blockCnt"`
	RipeBlock     uint64                     `json:"ripeBlock"`
	UnripeDist    uint64                     `json:"unripe"`
	RpcProvider   string                     `json:"rpcProvider"`
	AppearanceMap index.AddressAppearanceMap `json:"-"`
	TsArray       []tslib.Timestamp          `json:"-"`
	BlockWg       sync.WaitGroup             `json:"-"`
	AppearanceWg  sync.WaitGroup             `json:"-"`
	TsWg          sync.WaitGroup             `json:"-"`
}

// HandleBlaze does the actual scraping, walking through block_cnt blocks and querying traces and logs
// and then extracting addresses and timestamps from those data structures.
func (opts *BlazeOptions) HandleBlaze(meta *rpcClient.MetaData) (ok bool, err error) {

	// Prepare three channels to process first blocks, then appearances and timestamps
	blockChannel := make(chan int)
	appearanceChannel := make(chan ScrapedData)
	tsChannel := make(chan tslib.Timestamp)

	opts.BlockWg.Add(int(opts.NChannels))
	for i := 0; i < int(opts.NChannels); i++ {
		go opts.BlazeProcessBlocks(meta, blockChannel, appearanceChannel, tsChannel)
	}

	opts.AppearanceWg.Add(int(opts.NChannels))
	for i := 0; i < int(opts.NChannels); i++ {
		// TODO: BOGUS - HOW DOES ONE HANDLE ERRORS INSIDE OF GO ROUTINES?
		go opts.BlazeProcessAppearances(meta, appearanceChannel)
	}

	opts.TsWg.Add(int(opts.NChannels))
	for i := 0; i < int(opts.NChannels); i++ {
		// TODO: BOGUS - HOW DOES ONE HANDLE ERRORS INSIDE OF GO ROUTINES?
		go opts.BlazeProcessTimestamps(tsChannel)
	}

	for block := int(opts.StartBlock); block < int(opts.StartBlock+opts.BlockCount); block++ {
		// TODO: BOGUS - HOW DOES ONE HANDLE ERRORS INSIDE OF GO ROUTINES?
		blockChannel <- block
	}

	close(blockChannel)
	opts.BlockWg.Wait()

	close(appearanceChannel)
	opts.AppearanceWg.Wait()

	close(tsChannel)
	opts.TsWg.Wait()

	return true, nil
}

// BlazeProcessBlocks Processes the block channel and for each block query the node for both traces and logs. Send results down appearanceChannel.
func (opts *BlazeOptions) BlazeProcessBlocks(meta *rpcClient.MetaData, blockChannel chan int, appearanceChannel chan ScrapedData, tsChannel chan tslib.Timestamp) (err error) {
	defer opts.BlockWg.Done()

	for blockNum := range blockChannel {

		// RPCPayload is used during to make calls to the RPC.
		var traces rpcClient.Traces
		tracePayload := rpcClient.RPCPayload{
			Jsonrpc:   "2.0",
			Method:    "trace_block",
			RPCParams: rpcClient.RPCParams{fmt.Sprintf("0x%x", blockNum)},
			ID:        1002,
		}
		err = rpcClient.FromRpc(opts.RpcProvider, &tracePayload, &traces)
		if err != nil {
			return fmt.Errorf("call to trace_block returned error: %s", err)
		}

		var logs rpcClient.Logs
		logsPayload := rpcClient.RPCPayload{
			Jsonrpc:   "2.0",
			Method:    "eth_getLogs",
			RPCParams: rpcClient.RPCParams{rpcClient.LogFilter{Fromblock: fmt.Sprintf("0x%x", blockNum), Toblock: fmt.Sprintf("0x%x", blockNum)}},
			ID:        1003,
		}
		err = rpcClient.FromRpc(opts.RpcProvider, &logsPayload, &logs)
		if err != nil {
			return fmt.Errorf("call to eth_getLogs returned error: %s", err)
		}

		appearanceChannel <- ScrapedData{
			blockNumber: blockNum,
			traces:      traces,
			logs:        logs,
		}

		tsChannel <- tslib.Timestamp{
			Ts: uint32(rpcClient.GetBlockTimestamp(opts.RpcProvider, uint64(blockNum))),
			Bn: uint32(blockNum),
		}
	}

	return
}

var blazeMutex sync.Mutex

// BlazeProcessAppearances processes ScrapedData objects shoved down the appearanceChannel
func (opts *BlazeOptions) BlazeProcessAppearances(meta *rpcClient.MetaData, appearanceChannel chan ScrapedData) (err error) {
	defer opts.AppearanceWg.Done()
	for sData := range appearanceChannel {
		addressMap := make(map[string]bool)
		err = opts.BlazeExtractFromTraces(sData.blockNumber, &sData.traces, addressMap)
		if err != nil {
			return
		}

		err = opts.BlazeExtractFromLogs(sData.blockNumber, &sData.logs, addressMap)
		if err != nil {
			return
		}

		err = opts.WriteAppearances(meta, sData.blockNumber, addressMap)
		if err != nil {
			return
		}
	}
	return
}

// BlazeProcessTimestamps processes timestamp data (currently by printing to a temporary file)
func (opts *BlazeOptions) BlazeProcessTimestamps(tsChannel chan tslib.Timestamp) (err error) {
	defer opts.TsWg.Done()
	for ts := range tsChannel {
		blazeMutex.Lock()
		opts.TsArray = append(opts.TsArray, ts)
		blazeMutex.Unlock()
	}
	return
}

var mapSync sync.Mutex

func (opts *BlazeOptions) AddToMaps(address string, bn, txid int, addressMap map[string]bool) {
	mapSync.Lock()
	key := fmt.Sprintf("%s\t%09d\t%05d", address, bn, txid)
	addressMap[key] = true
	opts.AppearanceMap[address] = append(opts.AppearanceMap[address], index.AppearanceRecord{
		BlockNumber:   uint32(bn),
		TransactionId: uint32(txid),
	})

	mapSync.Unlock()
}

func (opts *BlazeOptions) BlazeExtractFromTraces(bn int, traces *rpcClient.Traces, addressMap map[string]bool) (err error) {
	if traces.Result == nil || len(traces.Result) == 0 {
		return
	}

	for i := 0; i < len(traces.Result); i++ {
		txid := traces.Result[i].TransactionPosition

		if traces.Result[i].Type == "call" {
			// If it's a call, get the to and from
			from := traces.Result[i].Action.From
			if isAddress(from) {
				opts.AddToMaps(from, bn, txid, addressMap)
			}
			to := traces.Result[i].Action.To
			if isAddress(to) {
				opts.AddToMaps(to, bn, txid, addressMap)
			}

		} else if traces.Result[i].Type == "reward" {
			if traces.Result[i].Action.RewardType == "block" {
				author := traces.Result[i].Action.Author
				if validate.IsZeroAddress(author) {
					// Early clients allowed misconfigured miner settings with address
					// 0x0 (reward got burned). We enter a false record with a false tx_id
					// to account for this.
					author = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
					opts.AddToMaps(author, bn, 99997, addressMap)

				} else {
					if isAddress(author) {
						opts.AddToMaps(author, bn, 99999, addressMap)
					}
				}

			} else if traces.Result[i].Action.RewardType == "uncle" {
				author := traces.Result[i].Action.Author
				if validate.IsZeroAddress(author) {
					// Early clients allowed misconfigured miner settings with address
					// 0x0 (reward got burned). We enter a false record with a false tx_id
					// to account for this.
					author = "0xdeaddeaddeaddeaddeaddeaddeaddeaddeaddead"
					opts.AddToMaps(author, bn, 99998, addressMap)

				} else {
					if isAddress(author) {
						opts.AddToMaps(author, bn, 99998, addressMap)
					}
				}

			} else if traces.Result[i].Action.RewardType == "external" {
				// This only happens in xDai as far as we know...
				author := traces.Result[i].Action.Author
				if isAddress(author) {
					opts.AddToMaps(author, bn, 99996, addressMap)
				}

			} else {
				logger.Log(logger.Error, "Unknown reward type:", traces.Result[i].Action.RewardType)
			}

		} else if traces.Result[i].Type == "suicide" {
			// add the contract that died, and where it sent it's money
			address := traces.Result[i].Action.Address
			if isAddress(address) {
				opts.AddToMaps(address, bn, txid, addressMap)
			}
			refundAddress := traces.Result[i].Action.RefundAddress
			if isAddress(refundAddress) {
				opts.AddToMaps(refundAddress, bn, txid, addressMap)
			}

		} else if traces.Result[i].Type == "create" {
			// add the creator, and the new address name
			from := traces.Result[i].Action.From
			if isAddress(from) {
				opts.AddToMaps(from, bn, txid, addressMap)
			}
			address := traces.Result[i].Result.Address
			if isAddress(address) {
				opts.AddToMaps(address, bn, txid, addressMap)
			}

			// If it's a top level trace, then the call data is the init,
			// so to match with TrueBlocks, we just parse init
			if len(traces.Result[i].TraceAddress) == 0 {
				if len(traces.Result[i].Action.Init) > 10 {
					initData := traces.Result[i].Action.Init[10:]
					for i := 0; i < len(initData)/64; i++ {
						addr := string(initData[i*64 : (i+1)*64])
						if isImplicitAddress(addr) {
							opts.AddToMaps(addr, bn, txid, addressMap)
						}
					}
				}
			}

			// Handle contract creations that may have errored out
			if traces.Result[i].Action.To == "" {
				if traces.Result[i].Result.Address == "" {
					if traces.Result[i].Error != "" {
						var receipt rpcClient.Receipt
						var txReceiptPl = rpcClient.RPCPayload{
							Jsonrpc:   "2.0",
							Method:    "eth_getTransactionReceipt",
							RPCParams: rpcClient.RPCParams{traces.Result[i].TransactionHash},
							ID:        1005,
						}
						err = rpcClient.FromRpc(opts.RpcProvider, &txReceiptPl, &receipt)
						if err != nil {
							return fmt.Errorf("call to eth_getTransactionReceipt returned error: %s", err)
						}
						addr := receipt.Result.ContractAddress
						if isAddress(addr) {
							opts.AddToMaps(addr, bn, txid, addressMap)
						}
					}
				}
			}

		} else {
			logger.Log(logger.Error, "Unknown trace type:", traces.Result[i].Type)
		}

		// Try to get addresses from the input data
		if len(traces.Result[i].Action.Input) > 10 {
			inputData := traces.Result[i].Action.Input[10:]
			//fmt.Println("Input data:", inputData, len(inputData))
			for i := 0; i < len(inputData)/64; i++ {
				addr := string(inputData[i*64 : (i+1)*64])
				if isImplicitAddress(addr) {
					opts.AddToMaps(addr, bn, txid, addressMap)
				}
			}
		}

		// Parse output of trace
		if len(traces.Result[i].Result.Output) > 2 {
			outputData := traces.Result[i].Result.Output[2:]
			for i := 0; i < len(outputData)/64; i++ {
				addr := string(outputData[i*64 : (i+1)*64])
				if isImplicitAddress(addr) {
					opts.AddToMaps(addr, bn, txid, addressMap)
				}
			}
		}
	}

	return
}

// extractFromLogs Extracts addresses from any part of the log data.
func (opts *BlazeOptions) BlazeExtractFromLogs(bn int, logs *rpcClient.Logs, addressMap map[string]bool) (err error) {
	if logs.Result == nil || len(logs.Result) == 0 {
		return
	}

	for i := 0; i < len(logs.Result); i++ {
		txid, _ := strconv.ParseInt(logs.Result[i].TransactionIndex, 0, 32)
		for j := 0; j < len(logs.Result[i].Topics); j++ {
			addr := string(logs.Result[i].Topics[j][2:])
			if isImplicitAddress(addr) {
				opts.AddToMaps(addr, bn, int(txid), addressMap)
			}
		}

		if len(logs.Result[i].Data) > 2 {
			inputData := logs.Result[i].Data[2:]
			for i := 0; i < len(inputData)/64; i++ {
				addr := string(inputData[i*64 : (i+1)*64])
				if isImplicitAddress(addr) {
					opts.AddToMaps(addr, bn, int(txid), addressMap)
				}
			}
		}
	}

	return
}

func (opts *BlazeOptions) WriteAppearances(meta *rpcClient.MetaData, bn int, addressMap map[string]bool) (err error) {
	if len(addressMap) > 0 {
		appearanceArray := make([]string, 0, len(addressMap))
		for record := range addressMap {
			appearanceArray = append(appearanceArray, record)
		}
		sort.Strings(appearanceArray)

		blockNumStr := utils.PadLeft(strconv.Itoa(bn), 9)
		fileName := config.GetPathToIndex(opts.Chain) + "ripe/" + blockNumStr + ".txt"
		if bn > int(opts.RipeBlock) {
			fileName = config.GetPathToIndex(opts.Chain) + "unripe/" + blockNumStr + ".txt"
		}

		toWrite := []byte(strings.Join(appearanceArray[:], "\n") + "\n")
		err = ioutil.WriteFile(fileName, toWrite, 0744)
		if err != nil {
			return fmt.Errorf("call to WriteFile returned error: %s", err)
		}
	}

	// TODO: BOGUS - TESTING SCRAPING
	if !utils.DebuggingOn {
		// TODO: THIS IS A PERFORMANCE ISSUE PRINTING EVERY BLOCK
		step := uint64(7)
		if opts.NProcessed%step == 0 {
			dist := uint64(0)
			if opts.RipeBlock > uint64(bn) {
				dist = (opts.RipeBlock - uint64(bn))
			}
			f := "-------- ( ------)- <PROG>  : Scraping %-04d of %-04d at block %d of %d (%d blocks from head)\r"
			fmt.Fprintf(os.Stderr, f, opts.NProcessed, opts.BlockCount, bn, opts.RipeBlock, dist)
		}
	}
	opts.NProcessed++
	return
}

// isAddress Returns true if the address is not a precompile and not the zero address
func isAddress(addr string) bool {
	// As per EIP 1352, all addresses less or equal to the following value are reserved for pre-compiles.
	// We don't index precompiles. https://eips.ethereum.org/EIPS/eip-1352
	return addr > "0x000000000000000000000000000000000000ffff"
}

// isImplicitAddress processes a transaction's 'input' data and 'output' data or an event's data field.
// Anything with 12 bytes of leading zeros but not more than 19 leading zeros (24 and 38 characters
// respectively).
func isImplicitAddress(addr string) bool {
	// Any 32-byte value smaller than this number is assumed to be a 'value'. We call them baddresses.
	// While this may seem like a lot of addresses being labeled as baddresses, it's not very many:
	// ---> 2 out of every 10000000000000000000000000000000000000000000000 are baddresses.
	small := "00000000000000000000000000000000000000ffffffffffffffffffffffffff"
	//        -------+-------+-------+-------+-------+-------+-------+-------+
	if addr <= small {
		return false
	}

	// Any 32-byte value with less than this many leading zeros is not an address (address are 20-bytes and
	// zero padded to the left)
	largePrefix := "000000000000000000000000"
	//              -------+-------+-------+
	if !strings.HasPrefix(addr, largePrefix) {
		return false
	}

	// Of the valid addresses, we assume any ending with this many trailing zeros is also a baddress.
	if strings.HasSuffix(addr, "00000000") {
		return false
	}

	// extract the potential address
	return isAddress("0x" + string(addr[24:]))
}
