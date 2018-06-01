package process

import (
	"os"
	"strconv"
	"encoding/json"
	"practice/process/common"
	"context"
	"fmt"

	"io/ioutil"
	"strings"

	"practice/process/cpuinfo"


	"practice/process/host"
	"path/filepath"
	"bytes"
)

var PageSize = uint64(os.Getpagesize())

const (

	ClockTicks  = 100 // C.sysconf(C._SC_CLK_TCK)
)

// MemoryInfoExStat is different between OSes
type MemoryInfoExStat struct {
	RSS    uint64 `json:"rss"`    // bytes
	VMS    uint64 `json:"vms"`    // bytes
	Shared uint64 `json:"shared"` // bytes
	Text   uint64 `json:"text"`   // bytes
	Lib    uint64 `json:"lib"`    // bytes
	Data   uint64 `json:"data"`   // bytes
	Dirty  uint64 `json:"dirty"`  // bytes
}

func (m MemoryInfoExStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

type MemoryMapsStat struct {
	Path         string `json:"path"`
	Rss          uint64 `json:"rss"`
	Size         uint64 `json:"size"`
	Pss          uint64 `json:"pss"`
	SharedClean  uint64 `json:"sharedClean"`
	SharedDirty  uint64 `json:"sharedDirty"`
	PrivateClean uint64 `json:"privateClean"`
	PrivateDirty uint64 `json:"privateDirty"`
	Referenced   uint64 `json:"referenced"`
	Anonymous    uint64 `json:"anonymous"`
	Swap         uint64 `json:"swap"`
}

// String returns JSON value of the process.
func (m MemoryMapsStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

// NewProcess creates a new Process instance, it only stores the pid and
// checks that the process exists. Other method on Process can be used
// to get more information about the process. An error will be returned
// if the process does not exist.
func NewProcess(pid int32) (*Process, error) {
	p := &Process{
		Pid: int32(pid),
	}
	//file, err := os.Open(common.HostProc(strconv.Itoa(int(p.Pid))))
	file , err := os.Open(common.HostProc(strconv.Itoa(int(p.Pid))))
	defer file.Close()
	return p, err
}

// Name returns name of the process.
func (p *Process) Name() (string, error) {
	return p.NameWithContext(context.Background())
}

func (p *Process) NameWithContext(ctx context.Context) (string, error) {
	if p.name == "" {
		if err := p.fillFromStatus(); err != nil {
			return "", err
		}

		}
		return p.name, nil
	}


// Tgid returns tgid, a Linux-synonym for user-space Pid
func (p *Process) Tgid() (int32, error) {
	if p.tgid == 0 {
		if err := p.fillFromStatus(); err != nil {
			return 0, err
		}
	}
	return p.tgid, nil
}


// Cmdline returns the command line arguments of the process as a string with
// each argument separated by 0x20 ascii character.
func (p *Process) Cmdline() (string, error) {
	return p.CmdlineWithContext(context.Background())
}

func (p *Process) CmdlineWithContext(ctx context.Context) (string, error) {
	return p.fillFromCmdline()
}

func (p *Process) CreateTime() (int64, error) {
	return p.CreateTimeWithContext(context.Background())
}

func (p *Process) CreateTimeWithContext(ctx context.Context) (int64, error) {
	_, _, _, createTime, err := p.fillFromStat()
	if err != nil {
		return 0, err
	}
	return createTime, nil
}



// Cwd returns current working directory of the process.
func (p *Process) Cwd() (string, error) {
	return p.CwdWithContext(context.Background())
}

func (p *Process) CwdWithContext(ctx context.Context) (string, error) {
	return p.fillFromCwd()
}

// Parent returns parent Process of the process.
func (p *Process) Parent() (*Process, error) {
	return p.ParentWithContext(context.Background())
}

func (p *Process) ParentWithContext(ctx context.Context) (*Process, error) {
	err := p.fillFromStatus()
	if err != nil {
		return nil, err
	}
	if p.parent == 0 {
		return nil, fmt.Errorf("wrong number of parents")
	}
	return NewProcess(p.parent)
}

// Status returns the process status.
// Return value could be one of these.
// R: Running S: Sleep T: Stop I: Idle
// Z: Zombie W: Wait L: Lock
// The charactor is same within all supported platforms.
func (p *Process) Status() (string, error) {
	return p.StatusWithContext(context.Background())
}

func (p *Process) StatusWithContext(ctx context.Context) (string, error) {
	err := p.fillFromStatus()
	if err != nil {
		return "", err
	}
	return p.status, nil
}
func (p *Process) CmdlineSlice() ([]string, error) {
	return p.CmdlineSliceWithContext(context.Background())
}

func (p *Process) CmdlineSliceWithContext(ctx context.Context) ([]string, error) {
	return p.fillSliceFromCmdline()
}

// Uids returns user ids of the process as a slice of the int
func (p *Process) Uids() ([]int32, error) {
	return p.UidsWithContext(context.Background())
}

func (p *Process) UidsWithContext(ctx context.Context) ([]int32, error) {
	err := p.fillFromStatus()
	if err != nil {
		return []int32{}, err
	}
	return p.uids, nil
}

// Gids returns group ids of the process as a slice of the int
func (p *Process) Gids() ([]int32, error) {
	return p.GidsWithContext(context.Background())
}

func (p *Process) GidsWithContext(ctx context.Context) ([]int32, error) {
	err := p.fillFromStatus()
	if err != nil {
		return []int32{}, err
	}
	return p.gids, nil
}


// IOnice returns process I/O nice value (priority).
func (p *Process) IOnice() (int32, error) {
	return p.IOniceWithContext(context.Background())
}

func (p *Process) IOniceWithContext(ctx context.Context) (int32, error) {
	return 0, common.ErrNotImplementedError
}



// NumThreads returns the number of threads used by the process.
func (p *Process) NumThreads() (int32, error) {
	return p.NumThreadsWithContext(context.Background())
}

func (p *Process) NumThreadsWithContext(ctx context.Context) (int32, error) {
	err := p.fillFromStatus()
	if err != nil {
		return 0, err
	}
	return p.numThreads, nil
}
func (p *Process) Times() (*cpu.TimesStat, error) {
	return p.TimesWithContext(context.Background())
}

func (p *Process) TimesWithContext(ctx context.Context) (*cpu.TimesStat, error) {
	_, _, cpuTimes, _, err := p.fillFromStat()
	if err != nil {
		return nil, err
	}
	return cpuTimes, nil
}



// MemoryInfo returns platform in-dependend memory information, such as RSS, VMS and Swap
func (p *Process) MemoryInfo() (*MemoryInfoStat, error) {
	return p.MemoryInfoWithContext(context.Background())
}

func (p *Process) MemoryInfoWithContext(ctx context.Context) (*MemoryInfoStat, error) {
	meminfo, _, err := p.fillFromStatm()
	if err != nil {
		return nil, err
	}
	return meminfo, nil
}


// Children returns a slice of Process of the process.
func (p *Process) Children() ([]*Process, error) {
	return p.ChildrenWithContext(context.Background())
}

func (p *Process) ChildrenWithContext(ctx context.Context) ([]*Process, error) {
	pids, err := common.CallPgrepWithContext(ctx, invoke, p.Pid)
	if err != nil {
		if pids == nil || len(pids) == 0 {
			return nil, ErrorNoChildren
		}
		return nil, err
	}
	ret := make([]*Process, 0, len(pids))
	for _, pid := range pids {
		np, err := NewProcess(pid)
		if err != nil {
			return nil, err
		}
		ret = append(ret, np)
	}
	return ret, nil
}



// MemoryMaps get memory maps from /proc/(pid)/smaps
func (p *Process) MemoryMaps(grouped bool) (*[]MemoryMapsStat, error) {
	return p.MemoryMapsWithContext(context.Background(), grouped)
}

func (p *Process) MemoryMapsWithContext(ctx context.Context, grouped bool) (*[]MemoryMapsStat, error) {
	pid := p.Pid
	var ret []MemoryMapsStat
	smapsPath := common.HostProc(strconv.Itoa(int(pid)), "smaps")
	contents, err := ioutil.ReadFile(smapsPath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(contents), "\n")

	// function of parsing a block
	getBlock := func(first_line []string, block []string) (MemoryMapsStat, error) {
		m := MemoryMapsStat{}
		m.Path = first_line[len(first_line)-1]

		for _, line := range block {
			if strings.Contains(line, "VmFlags") {
				continue
			}
			field := strings.Split(line, ":")
			if len(field) < 2 {
				continue
			}
			v := strings.Trim(field[1], " kB") // remove last "kB"
			t, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return m, err
			}

			switch field[0] {
			case "Size":
				m.Size = t
			case "Rss":
				m.Rss = t
			case "Pss":
				m.Pss = t
			case "Shared_Clean":
				m.SharedClean = t
			case "Shared_Dirty":
				m.SharedDirty = t
			case "Private_Clean":
				m.PrivateClean = t
			case "Private_Dirty":
				m.PrivateDirty = t
			case "Referenced":
				m.Referenced = t
			case "Anonymous":
				m.Anonymous = t
			case "Swap":
				m.Swap = t
			}
		}
		return m, nil
	}

	blocks := make([]string, 16)
	for _, line := range lines {
		field := strings.Split(line, " ")
		if strings.HasSuffix(field[0], ":") == false {
			// new block section
			if len(blocks) > 0 {
				g, err := getBlock(field, blocks)
				if err != nil {
					return &ret, err
				}
				ret = append(ret, g)
			}
			// starts new block
			blocks = make([]string, 16)
		} else {
			blocks = append(blocks, line)
		}
	}

	return &ret, nil
}






// Get cwd from /proc/(pid)/cwd
func (p *Process) fillFromCwd() (string, error) {
	return p.fillFromCwdWithContext(context.Background())
}

func (p *Process) fillFromCwdWithContext(ctx context.Context) (string, error) {
	pid := p.Pid
	cwdPath := common.HostProc(strconv.Itoa(int(pid)), "cwd")
	cwd, err := os.Readlink(cwdPath)
	if err != nil {
		return "", err
	}
	return string(cwd), nil
}



// Get cmdline from /proc/(pid)/cmdline
func (p *Process) fillFromCmdline() (string, error) {
	return p.fillFromCmdlineWithContext(context.Background())
}

func (p *Process) fillFromCmdlineWithContext(ctx context.Context) (string, error) {
	pid := p.Pid
	cmdPath := common.HostProc(strconv.Itoa(int(pid)), "cmdline")
	cmdline, err := ioutil.ReadFile(cmdPath)
	if err != nil {
		return "", err
	}
	ret := strings.FieldsFunc(string(cmdline), func(r rune) bool {
		if r == '\u0000' {
			return true
		}
		return false
	})

	return strings.Join(ret, " "), nil
}


// Get memory info from /proc/(pid)/statm
func (p *Process) fillFromStatm() (*MemoryInfoStat, *MemoryInfoExStat, error) {
	return p.fillFromStatmWithContext(context.Background())
}

func (p *Process) fillFromStatmWithContext(ctx context.Context) (*MemoryInfoStat, *MemoryInfoExStat, error) {
	pid := p.Pid
	memPath := common.HostProc(strconv.Itoa(int(pid)), "statm")
	contents, err := ioutil.ReadFile(memPath)
	if err != nil {
		return nil, nil, err
	}
	fields := strings.Split(string(contents), " ")

	vms, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	rss, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	memInfo := &MemoryInfoStat{
		RSS: rss * PageSize,
		VMS: vms * PageSize,
	}

	shared, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	text, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	lib, err := strconv.ParseUint(fields[4], 10, 64)
	if err != nil {
		return nil, nil, err
	}
	dirty, err := strconv.ParseUint(fields[5], 10, 64)
	if err != nil {
		return nil, nil, err
	}

	memInfoEx := &MemoryInfoExStat{
		RSS:    rss * PageSize,
		VMS:    vms * PageSize,
		Shared: shared * PageSize,
		Text:   text * PageSize,
		Lib:    lib * PageSize,
		Dirty:  dirty * PageSize,
	}

	return memInfo, memInfoEx, nil
}

func (p *Process) fillSliceFromCmdline() ([]string, error) {
	return p.fillSliceFromCmdlineWithContext(context.Background())
}

func (p *Process) fillSliceFromCmdlineWithContext(ctx context.Context) ([]string, error) {
	pid := p.Pid
	cmdPath := common.HostProc(strconv.Itoa(int(pid)), "cmdline")
	cmdline, err := ioutil.ReadFile(cmdPath)
	if err != nil {
		return nil, err
	}
	if len(cmdline) == 0 {
		return nil, nil
	}
	if cmdline[len(cmdline)-1] == 0 {
		cmdline = cmdline[:len(cmdline)-1]
	}
	parts := bytes.Split(cmdline, []byte{0})
	var strParts []string
	for _, p := range parts {
		strParts = append(strParts, string(p))
	}

	return strParts, nil
}

// Get various status from /proc/(pid)/status
func (p *Process) fillFromStatus() error {
	return p.fillFromStatusWithContext(context.Background())
}

func (p *Process) fillFromStatusWithContext(ctx context.Context) error {
	pid := p.Pid
	statPath := common.HostProc(strconv.Itoa(int(pid)), "status")
	contents, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(contents), "\n")

	p.memInfo = &MemoryInfoStat{}
	//p.sigInfo = &SignalInfoStat{}
	for _, line := range lines {
		tabParts := strings.SplitN(line, "\t", 2)
		if len(tabParts) < 2 {
			continue
		}
		value := tabParts[1]
		switch strings.TrimRight(tabParts[0], ":") {

		case "Name":
			p.name = strings.Trim(value, " \t")
			if len(p.name) >= 15 {
				cmdlineSlice, err := p.CmdlineSlice()
				if err != nil {
					return err
				}
				if len(cmdlineSlice) > 0 {
					extendedName := filepath.Base(cmdlineSlice[0])
					if strings.HasPrefix(extendedName, p.name) {
						p.name = extendedName
					}
				}
			}

		case "State":
			p.status = value[0:1]

		case "Tgid":
			pval, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			p.tgid = int32(pval)

		case "PPid", "Ppid":
			pval, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			p.parent = int32(pval)



		case "Uid":
			p.uids = make([]int32, 0, 4)
			for _, i := range strings.Split(value, "\t") {
				v, err := strconv.ParseInt(i, 10, 32)
				if err != nil {
					return err
				}
				p.uids = append(p.uids, int32(v))
			}
		case "Gid":
			p.gids = make([]int32, 0, 4)
			for _, i := range strings.Split(value, "\t") {
				v, err := strconv.ParseInt(i, 10, 32)
				if err != nil {
					return err
				}
				p.gids = append(p.gids, int32(v))
			}
		case "Threads":
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}
			p.numThreads = int32(v)

		case "VmRSS":
			value := strings.Trim(value, " kB") // remove last "kB"
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			p.memInfo.RSS = v *1024

		case "VmSize":
			value := strings.Trim(value, " kB") // remove last "kB"
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			p.memInfo.VMS = v  *1024
		case "VmSwap":
			value := strings.Trim(value, " kB") // remove last "kB"
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			}
			p.memInfo.Swap = v * 1024

		}

	}
	return nil
}


func (p *Process) fillFromTIDStat(tid int32) (uint64, int32, *cpu.TimesStat, int64,  error) {
	return p.fillFromTIDStatWithContext(context.Background(), tid)
}

func (p *Process) fillFromTIDStatWithContext(ctx context.Context, tid int32) (uint64, int32, *cpu.TimesStat, int64,  error) {
	pid := p.Pid
	var statPath string

	if tid == -1 {
		statPath = common.HostProc(strconv.Itoa(int(pid)), "stat")
	} else {
		statPath = common.HostProc(strconv.Itoa(int(pid)), "task", strconv.Itoa(int(tid)), "stat")
	}

	contents, err := ioutil.ReadFile(statPath)
	if err != nil {
		return 0, 0, nil, 0,err
	}
	fields := strings.Fields(string(contents))

	i := 1
	for !strings.HasSuffix(fields[i], ")") {
		i++
	}

	terminal, err := strconv.ParseUint(fields[i+5], 10, 64)
	if err != nil {
		return 0, 0, nil, 0, err
	}
	ppid, err := strconv.ParseInt(fields[i+2], 10, 32)
	if err != nil {
		return 0, 0, nil,  0, err
	}
	utime, err := strconv.ParseFloat(fields[i+12], 64)
	if err != nil {
		return 0, 0, nil, 0, err
	}
	stime, err := strconv.ParseFloat(fields[i+13], 64)
	if err != nil {
		return 0, 0, nil, 0, err
	}
	cpuTimes := &cpu.TimesStat{
		CPU:    "cpu",
		User:   float64(utime / ClockTicks),
		System: float64(stime / ClockTicks),
	}

	//bootTime, _ := host.BootTime()
	bootTime,_ := host.BootTime()
	t, err := strconv.ParseUint(fields[i+20], 10, 64)
	if err != nil {
		return 0, 0, nil, 0,  err
	}
	ctime := (t / uint64(ClockTicks)) + uint64(bootTime)
	createTime := int64(ctime * 1000)





	return terminal, int32(ppid), cpuTimes ,createTime,nil
}
func (p *Process) fillFromStat() (uint64, int32, *cpu.TimesStat, int64,  error) {
	return p.fillFromStatWithContext(context.Background())
}

func (p *Process) fillFromStatWithContext(ctx context.Context) (uint64, int32, *cpu.TimesStat, int64, error) {
	return p.fillFromTIDStat(-1)
}


// Pids returns a slice of process ID list which are running now.
func  Pids() ([]int32, error) {
	return PidsWithContext(context.Background())
}

func PidsWithContext( context.Context) ([]int32, error) {
	return readPidsFromDir(common.HostProc())
}

// Process returns a slice of pointers to Process structs for all
// currently running processes.
func Processes() ([]*Process, error) {
	return ProcessesWithContext(context.Background())
}

func ProcessesWithContext( context.Context) ([]*Process, error) {
	out := []*Process{}

	pids, err := Pids()
	if err != nil {
		return out, err
	}

	for _, pid := range pids {
		p, err := NewProcess(pid)
		if err != nil {
			continue
		}
		out = append(out, p)
	}

	return out, nil
}

func readPidsFromDir(path string) ([]int32, error) {
	var ret []int32

	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	fnames, err := d.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	for _, fname := range fnames {
		pid, err := strconv.ParseInt(fname, 10, 32)
		if err != nil {
			// if not numeric name, just skip
			continue
		}
		ret = append(ret, int32(pid))
	}

	return ret, nil
}
