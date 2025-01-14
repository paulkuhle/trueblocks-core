package types

import "io"

// EXISTING_CODE
// EXISTING_CODE

type RawMonitor interface {
	// EXISTING_CODE
	// EXISTING_CODE
}

type SimpleMonitor struct {
	Address     string      `json:"address"`
	NRecords    int         `json:"nRecords"`
	FileSize    int64       `json:"fileSize"`
	LastScanned uint32      `json:"lastScanned"`
	raw         *RawMonitor `json:"-"`
	// EXISTING_CODE
	// EXISTING_CODE
}

func (s *SimpleMonitor) Raw() *RawMonitor {
	return s.raw
}

func (s *SimpleMonitor) SetRaw(raw *RawMonitor) {
	s.raw = raw
}

func (s *SimpleMonitor) Model(showHidden bool, format string, extraOptions map[string]any) Model {
	var model = map[string]interface{}{}
	var order = []string{}

	// EXISTING_CODE
	model = map[string]any{
		"address":     s.Address,
		"nRecords":    s.NRecords,
		"fileSize":    s.FileSize,
		"lastScanned": s.LastScanned,
	}
	order = []string{
		"address",
		"nRecords",
		"fileSize",
		"lastScanned",
	}
	// EXISTING_CODE

	return Model{
		Data:  model,
		Order: order,
	}
}

func (s *SimpleMonitor) WriteTo(w io.Writer) (n int64, err error) {
	// EXISTING_CODE
	// EXISTING_CODE
	return 0, nil
}

func (s *SimpleMonitor) ReadFrom(r io.Reader) (n int64, err error) {
	// EXISTING_CODE
	// EXISTING_CODE
	return 0, nil
}

// EXISTING_CODE
// EXISTING_CODE
