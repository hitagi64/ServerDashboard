package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

func getCPUUsage() (usage float64) {
	pct, err := cpu.Percent(time.Duration(1) * time.Second / 16, false)
	if (err != nil) {
		return -1
	}
	if (len(pct) < 1) {
		return -1
	}

	return pct[0]
}

func getCPUUsageAll() (usage []float64) {
	pct, err := cpu.Percent(time.Duration(1) * time.Second / 16, true)
	if (err != nil) {
		return []float64{-1}
	}

	return pct
}

func getCPUUsageForCore(core int) (usage float64) {
	pct, err := cpu.Percent(time.Duration(1) * time.Second / 16, true)
	if (err != nil) {
		return -1
	}
	if (len(pct) <= core) {
		return -1
	}

	return pct[core]
}

func getHighestCPUUsageCore() (usage float64, coreid int) {
	var pct, err = cpu.Percent(time.Duration(1) * time.Second / 16, true)
	if (err != nil) {
		return -1, -1
	}

	var highestusage = 0.0
	var highestcore = 0
	for i, v := range pct {
		if (v > highestusage) {
			highestusage = v
			highestcore = i
		}
	}

	return highestusage, highestcore
}

func getHighestCPUUsageCoreAsString() (result string) {
	var usage, coreid = getHighestCPUUsageCore()
	if (coreid == -1) {
		return "Unknown"
	}
	return fmt.Sprintf("#%d (%.0f%%)", coreid, usage)
}

func getCPUName() (name string) {
	cpuinfo, err := cpu.Info()
	if (err != nil || len(cpuinfo) < 1) {
		return "Unknown CPU model"
	}
	
	return cpuinfo[0].ModelName
}

func getCPUCoreCount() (name int) {
	cpucount, err := cpu.Counts(false)
	if (err != nil) {
		return 0
	}
	
	return cpucount
}

func getCPUThreadCount() (name int) {
	cpucount, err := cpu.Counts(true)
	if (err != nil) {
		return 0
	}
	
	return cpucount
}

func getCPUVirtualizationSupported() (supported bool) {
	cpuinfo, err := cpu.Info()
	if (err != nil || len(cpuinfo) < 1) {
		return false
	}
	for _, v := range cpuinfo[0].Flags {
		if (v == "vmx" || v == "svm") {
			return true
		}
	}
	return false
}