// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package monitorsPkg

import (
	"context"
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/monitor"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// HandleList
func (opts *MonitorsOptions) HandleList() error {
	chain := opts.Globals.Chain
	monitorMap, monArray := monitor.GetMonitorMap(chain)

	errors := make([]error, 0)
	addrMap := map[base.Address]bool{}
	for _, addr := range opts.Addrs {
		a := base.HexToAddress(addr)
		addrMap[a] = true
		if monitorMap[a] == nil {
			errors = append(errors, fmt.Errorf("address %s is not being monitored", addr))
		}
	}

	ctx := context.Background()
	fetchData := func(modelChan chan types.Modeler[types.RawMonitor], errorChan chan error) {
		for _, e := range errors {
			errorChan <- e
		}

		for _, mon := range monArray {
			if len(addrMap) == 0 || addrMap[mon.Address] {
				s := types.SimpleMonitor{
					Address:     mon.Address.Hex(),
					NRecords:    int(mon.Count()),
					FileSize:    file.FileSize(mon.Path()),
					LastScanned: mon.Header.LastScanned,
				}
				modelChan <- &s
			}
		}
	}

	return output.StreamMany(ctx, fetchData, opts.Globals.OutputOpts())
}
