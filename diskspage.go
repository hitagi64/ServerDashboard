package main

import (
	"fmt"
	"net/http"
)

type disksData struct {
	Disks []string
	Arrays []string
	ArrayData []mdInfo
}

func (tp PageTemplates) diskPageHandler(w http.ResponseWriter, r *http.Request){
	data := disksData{
	}
	data.Disks, data.Arrays = findStorageDevicesInSystem()
	for _, v := range data.Arrays {
		ardata, err := mdDeviceGetInfo(v)
		if (err == nil) {
			data.ArrayData = append(data.ArrayData, ardata)
		}else{
			fmt.Println("Error getting md info for " + v + ": " + err.Error())
		}
	}
	data.ArrayData = append(data.ArrayData, makeFakeMD())
	tp.runBasePage(w, "Disks", tp.disks, data)
}