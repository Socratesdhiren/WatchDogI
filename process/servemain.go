package main

import (
	"practice/process/processinfo"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
)




type ProcInfo struct {
	Pid  int32 `json:"pid"`
	//Pid          *process.Process    `json:"tgid""`
	CPUPercent    float64    `json:"cpu_percent"`
	MemoryPercent float32     `json:"memory_percent"`
	Status        string          `json:"status"`
	Name        string  `json:"name"`
	MemInfo       *process.MemoryInfoStat     `json:"mem_info"`


}

func main() {


	r := mux.NewRouter()


	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// var p= Ping{Message:"Hello"}


		w.Header().Set("Content-Type", "text/html; charset=ascii")
			w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
			//w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.WriteHeader(http.StatusOK)

		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}


		if processes , err := process.Processes();err!=nil {
			// if processes ,err := process.Processes(); err!=nil {
			log.Fatal(err)
		} else {
			//processes, _ := process.Processes();
			var procinfos  []ProcInfo
			for _, p := range processes {
				//pid :=p.Pid
				pid, _ := p.Tgid()
				n, _ := p.CPUPercent()
				a, _ := p.MemoryPercent()
				status, _ := p.Status()
				memory, _ := p.MemoryInfo()
				//vmem, _ := p.MemoryInfo()
				name, _ := p.Name()

				procinfos = append(procinfos, ProcInfo{Pid: pid, CPUPercent: n, MemoryPercent: a, Status: status, MemInfo: memory, Name: name},)

			}

			b,err :=json.Marshal(procinfos)
			if err!= nil {
				fmt.Println("json err:",err)
			}

			fmt.Fprintf(w,string(b))
		}

			})


	http.ListenAndServe(":8000", r)

}

