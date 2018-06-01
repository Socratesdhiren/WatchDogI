package main

import (
	"time"
	"log"
	"fmt"
	"practice/process/processinfo"
)

type stats struct {
	startTime time.Time
	// procMemUsed process.MemoryInfoStat

	ProcUptime float64 //seconds

	ProcMemUsedPct float64
	Status         string
}

func NewStats() *stats {
	s := stats{}
	s.startTime = time.Now()
	//	s.procMemUsed=process.MemoryInfoStat{}
	return &s
}

/*
type ProcInfoII struct {
	Tgid          int32      `json:"Pid"`
	CPUPercent    float64    `json:"cpu_percent"`
	MemoryPercent float32     `json:"memory_percent"`
	Status        string          `json:"status"`
	MemInfo       *process.MemoryInfoStat     `json:"mem_info"`
	VMemInfo      *process.MemoryInfoStat       `json:"v_mem_info"`
	Name          string                    `json:"name"`
}
*/

type ProcInfoII struct {
	Pid  int32 `json:"pid"`
	//Tgid          int32
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float32 `json:"memory_percent"`
	Status        string `json:"status"`
	MemInfo       *process.MemoryInfoStat `json:"mem_info"`
	VMemInfo      *process.MemoryInfoStat `json:"v_mem_info"`
	Name          string `json:"name"`
}

func main() {

	//fmt.Printf("Pid  %%CPUused %%MemoryUsed status\t Res(B) \t Vmem(B) PName \n ")

	if processes , err := process.Processes();err!=nil {
 	// if processes ,err := process.Processes(); err!=nil {
		 log.Fatal(err)
	 } else {
		//processes, _ := process.Processes();
		var procinfos []ProcInfoII
		for _, p := range processes {
			//pid :=p.Pid
			pid, _ := p.Tgid()
			n, _ := p.CPUPercent()
			a, _ := p.MemoryPercent()
			status, _ := p.Status()
			memory, _ := p.MemoryInfo()
			vmem, _ := p.MemoryInfo()
			name, _ := p.Name()

			procinfos = append(procinfos, ProcInfoII{pid, n, a, status, memory, vmem, name})
		}

		fmt.Printf(" Pid  %%CPUused %%MemoryUsed status\t Res(B) \tVmem(B) PName \n ")
		for _, p := range procinfos[:]{

			fmt.Printf(" %d  \t %3.1f \t %3.1f \t %s \t \t \t %d \t %d \t \t %s \n ",p.Pid,p.CPUPercent,p.MemoryPercent,p.Status,p.MemInfo.RSS,p.MemInfo.VMS,p.Name)
		}



		fmt.Println("*************************************************************************************************************")
	}


}


