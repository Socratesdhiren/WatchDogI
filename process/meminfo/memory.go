package mem

import (
	"encoding/json"
	"practice/process/common"
)

var invoke common.Invoker = common.Invoke{}

// Memory usage statistics. Total, Available and Used contain numbers of bytes
// for human consumption.
//
// The other fields in this struct contain kernel specific values.
type VirtualMemoryStat struct {
	// Total amount of RAM on this system
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	//
	// This value is computed from the kernel specific values.
	Available uint64 `json:"available"`

	// RAM used by programs
	//
	// This value is computed from the kernel specific values.
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	//
	// This value is computed from the kernel specific values.
	UsedPercent float64 `json:"usedPercent"`

	Free uint64 `json:"free"`


	Active   uint64 `json:"active"`
	Inactive uint64 `json:"inactive"`
	Wired    uint64 `json:"wired"`

	Buffers      uint64 `json:"buffers"`
	Cached       uint64 `json:"cached"`
	Writeback    uint64 `json:"writeback"`
	Dirty        uint64 `json:"dirty"`
	WritebackTmp uint64 `json:"writebacktmp"`
	Shared       uint64 `json:"shared"`
	Slab         uint64 `json:"slab"`
	PageTables   uint64 `json:"pagetables"`
	SwapCached   uint64 `json:"swapcached"`
	CommitLimit  uint64 `json:"commitlimit"`
	CommittedAS  uint64 `json:"committedas"`
}

type SwapMemoryStat struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
	Sin         uint64  `json:"sin"`
	Sout        uint64  `json:"sout"`
}

func (m VirtualMemoryStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

func (m SwapMemoryStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}
