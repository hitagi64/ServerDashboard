package main

import (
	"fmt"
	"net/http"
)

type disksData struct {
	Disks []string
	Arrays []string
	ArraysData []mdInfo
	DisksData []diskInfo
}

func (tp PageTemplates) diskPageHandler(w http.ResponseWriter, r *http.Request){
	data := disksData{
	}
	data.Disks, data.Arrays = findStorageDevicesInSystem()
	for _, v := range data.Arrays {
		ardata, err := mdDeviceGetInfo(v)
		if (err == nil) {
			data.ArraysData = append(data.ArraysData, ardata)
		}else{
			fmt.Println("Error getting md info for " + v + ": " + err.Error())
		}
	}
	for _, v := range data.Disks {
		dsdata, err := diskGetInfo(v)
		if (err == nil) {
			data.DisksData = append(data.DisksData, dsdata)
		}else{
			fmt.Println("Error getting disk info for " + v + ": " + err.Error())
		}
	}
	data.ArraysData = append(data.ArraysData, makeFakeMD())
	tp.runBasePage(w, "Disks", tp.disks, data)
}