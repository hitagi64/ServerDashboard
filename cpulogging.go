package main

import (
	"sync"
	"time"
)

type CPULog struct {
	mu sync.Mutex
	maxlen int
	waittime time.Duration
	log [][]float64
}
var cl1 CPULog

func (cl *CPULog) getCPULog() [][]float64 {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.log
}

func (cl *CPULog) getCPULogMaxLength() (max int) {
	return cl.maxlen
}

func (cl *CPULog) cpuLoggingTask() {
	for {
		cpuusage := getCPUUsageAll()
		cl.mu.Lock()
		for i := 0; i < len(cpuusage); i++ {
			if (len(cl.log) <= i) {
				cl.log = append(cl.log, []float64{})
				for ii := 0; ii < cl.maxlen; ii++ {
					cl.log[i] = append(cl.log[i], 0)
				}
			}
			cl.log[i] = append(cl.log[i], cpuusage[i])
		}
		if (len(cl.log) > 0) {
			for i := 0; i < len(cl.log); i++ {
				if (len(cl.log[i])) > cl.maxlen {
					cl.log[i] = cl.log[i][len(cl.log[i])-cl.maxlen:]
				}
			}
		}
		cl.mu.Unlock()
		time.Sleep(cl.waittime)
	}
}