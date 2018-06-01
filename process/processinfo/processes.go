package process

import (
	"context"
	"practice/process/common"
	"errors"
	"time"
	"encoding/json"
	"runtime"
	"practice/process/cpuinfo"
	"practice/process/meminfo"
)

var (
	invoke          common.Invoker = common.Invoke{}
	ErrorNoChildren                = errors.New("process does not have children")
)

type Process struct {
	Pid            int32 `json:"pid"`
	name           string   `json:"name"`
	status         string   `json:"status"`
	parent         int32

	uids           []int32
	gids           []int32
	numThreads     int32
	memInfo        *MemoryInfoStat  
	sigInfo        *SignalInfoStat
    lastCPUTimes  *cpu.TimesStat
	lastCPUTime  time.Time

	tgid               int32
}

type OpenFilesStat struct {
	Path string `json:"path"`
	Fd   uint64 `json:"fd"`
}

type MemoryInfoStat struct {
	RSS    uint64 `json:"rss"`    // bytes
	VMS    uint64 `json:"vms"`    // bytes
   // Data   uint64 `json:"data"`   // bytes
	//Stack  uint64 `json:"stack"`  // bytes
	//Locked uint64 `json:"locked"` // bytes
	Swap   uint64 `json:"swap"`   // bytes
}

type SignalInfoStat struct {
	PendingProcess uint64 `json:"pending_process"`
	PendingThread  uint64 `json:"pending_thread"`
	Blocked        uint64 `json:"blocked"`
	Ignored        uint64 `json:"ignored"`
	Caught         uint64 `json:"caught"`
}




func (p Process) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

func (o OpenFilesStat) String() string {
	s, _ := json.Marshal(o)
	return string(s)
}

func (m MemoryInfoStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}







func PidExists(pid int32) (bool, error) {
	return PidExistsWithContext(context.Background(), pid)
}

func PidExistsWithContext(ctx context.Context, pid int32) (bool, error) {
	pids, err := Pids()
	if err != nil {
		return false, err
	}

	for _, i := range pids {
		if i == pid {
			return true, err
		}
	}

	return false, err
}

// If interval is 0, return difference from last call(non-blocking).
// If interval > 0, wait interval sec and return diffrence between start and end.
func (p *Process) Percent(interval time.Duration) (float64, error) {
	return p.PercentWithContext(context.Background(), interval)
}

func (p *Process) PercentWithContext(ctx context.Context, interval time.Duration) (float64, error) {
	cpuTimes, err := p.Times()
	if err != nil {
		return 0, err
	}
	now := time.Now()

	if interval > 0 {
		//p.lastCPUTimes = cpuTimes
		p.lastCPUTimes = cpuTimes
		p.lastCPUTime = now
		time.Sleep(interval)
		cpuTimes, err = p.Times()
		now = time.Now()
		if err != nil {
			return 0, err
		}
	} else {
		if p.lastCPUTimes == nil {
			// invoked first time
			p.lastCPUTimes = cpuTimes
			p.lastCPUTime = now
			return 0, nil
		}
	}

	numcpu := runtime.NumCPU()
	delta := (now.Sub(p.lastCPUTime).Seconds()) * float64(numcpu)
	ret := calculatePercent(p.lastCPUTimes, cpuTimes, delta, numcpu)
	p.lastCPUTimes = cpuTimes
	p.lastCPUTime = now
	return ret, nil
}

func calculatePercent(t1, t2 *cpu.TimesStat, delta float64, numcpu int) float64 {
	if delta == 0 {
		return 0
	}
	delta_proc := t2.Total() - t1.Total()
	overall_percent := ((delta_proc / delta) * 100) * float64(numcpu)
	return overall_percent
}

// MemoryPercent returns how many percent of the total RAM this process uses
func (p *Process) MemoryPercent() (float32, error) {
	return p.MemoryPercentWithContext(context.Background())
}

func (p *Process) MemoryPercentWithContext(ctx context.Context) (float32, error) {

	machineMemory, err := mem.VirtualMemory()
	//machineMemory, err := mem.VirtualMemory()

	if err != nil {
		return 0, err
	}
	total := machineMemory.Total

	processMemory, err := p.MemoryInfo()
	if err != nil {
		return 0, err
	}
	used := processMemory.RSS

	return (100 * float32(used) / float32(total)), nil
}

// CPU_Percent returns how many percent of the CPU time this process uses
func (p *Process) CPUPercent() (float64, error) {
	return p.CPUPercentWithContext(context.Background())
}

func (p *Process) CPUPercentWithContext(ctx context.Context) (float64, error) {
	crt_time, err := p.CreateTime()
	if err != nil {
		return 0, err
	}

	cput, err := p.Times()
	if err != nil {
		return 0, err
	}

	created := time.Unix(0, crt_time*int64(time.Millisecond))
	totalTime := time.Since(created).Seconds()
	if totalTime <= 0 {
		return 0, nil
	}

	return (100 * cput.Total() / totalTime), nil
}
