package main

import (
	"fmt"
	"net/http"
)

type dashboardData struct {
	SystemName string
	Kernel string
	KernelVersion string
	DistroName string
	SystemArchitecture string
	SystemUptime string

	FreeMemory float64
	UsedMemory float64
	CachedMemory float64
	TotalMemory float64

	CPUUsage float64
	CPUName string
	CPUCoreCount int
	CPUThreadCount int
	CPUHighestUsage string
	CPUHasVirtualization bool
	CPUMaxHistoryLength int
	CPULogHistory [][]float64

	ArraysData []mdInfo
}

type dashboardRefreshData struct {
	SystemUptime string

	FreeMemory float64
	UsedMemory float64
	CachedMemory float64
	TotalMemory float64

	CPUUsage float64
	CPUHighestUsage string
	CPUMaxHistoryLength int
	CPULogHistory [][]float64

	ArraysData []mdInfo
}

func (tp PageTemplates) dashboardPageHandler(w http.ResponseWriter, r *http.Request){
	_, machine, nodename, release, sysname, _ := getUname()
	distroname := getDistroName()
	total, used, cached, free := getMemoryUsage()
	data := dashboardData{
		SystemName: nodename,
		Kernel: sysname,
		KernelVersion: release,
		DistroName: distroname,
		SystemArchitecture: machine,
		SystemUptime: getSystemUptime().String(),

		FreeMemory: free,
		UsedMemory: used,
		CachedMemory: cached,
		TotalMemory: total,

		CPUUsage: getCPUUsage(),
		CPUName: getCPUName(),
		CPUCoreCount: getCPUCoreCount(),
		CPUThreadCount: getCPUThreadCount(),
		CPUHighestUsage: getHighestCPUUsageCoreAsString(),
		CPUHasVirtualization: getCPUVirtualizationSupported(),
		CPUMaxHistoryLength: cl1.getCPULogMaxLength(),
		CPULogHistory: cl1.getCPULog(),
	}
	_, arrays := findStorageDevicesInSystem()
	for _, v := range arrays {
		ardata, err := mdDeviceGetInfo(v)
		if (err == nil) {
			data.ArraysData = append(data.ArraysData, ardata)
		}else{
			fmt.Println("Error getting md info for " + v + ": " + err.Error())
		}
	}
	data.ArraysData = append(data.ArraysData, makeFakeMD())
	tp.runBasePage(w, "Dashboard", tp.dashboard, data)
}

func (tp PageTemplates) dashboardRefreshPageHandler(w http.ResponseWriter, r *http.Request){
	total, used, cached, free := getMemoryUsage()
	data := dashboardRefreshData{
		SystemUptime: getSystemUptime().String(),

		FreeMemory: free,
		UsedMemory: used,
		CachedMemory: cached,
		TotalMemory: total,

		CPUUsage: getCPUUsage(),
		CPUHighestUsage: getHighestCPUUsageCoreAsString(),
		CPUMaxHistoryLength: cl1.getCPULogMaxLength(),
		CPULogHistory: cl1.getCPULog(),
	}
	req, ok := r.URL.Query()["req"]
    
    if ok && len(req[0]) > 0 {
        if (req[0] == "md") {
			_, arrays := findStorageDevicesInSystem()
			for _, v := range arrays {
				ardata, err := mdDeviceGetInfo(v)
				if (err == nil) {
					data.ArraysData = append(data.ArraysData, ardata)
				}else{
					fmt.Println("Error getting md info for " + v + ": " + err.Error())
				}
			}
			data.ArraysData = append(data.ArraysData, makeFakeMD())
		}
    }
	err := tp.dashboardRefreshData.Execute(w, data); 
	if err != nil {
		fmt.Println("2:" + err.Error())
		http.Error(w, "500 Internal server error", 500)
		return;
	}
}