package main

import (
	"fmt"
	"strings"
	"unicode"
)

func stringStripNewline(s string) (o string) {
	o = strings.Replace(s, "\n", "", -1)
	o = strings.Replace(o, "\r", "", -1)
	o = strings.Replace(o, "\\n", "", -1)
	return o
}

func stringContainsNumbers(s string) bool {
	for _, c := range s {
		if unicode.IsNumber(c) {
			return true
		}
	}
	return false
}

func normalizeKBValue(val int64) (out string) {
	if (val < 1024) {
		out = fmt.Sprintf("%.2f KB", float64(val))
	} else if (val < 1024 * 1024) {
		out = fmt.Sprintf("%.2f MB", float64(val) / 1024.0)
	} else if (val < 1024 * 1024 * 1024) {
		out = fmt.Sprintf("%.2f GB", float64(val) / 1024.0 / 1024.0)
	} else {
		out = fmt.Sprintf("%.2f TB", float64(val) / 1024.0 / 1024.0 / 1024.0)
	}
	return
}

func makeFakeMD() (disk mdInfo) {
	disk.Name = "md0"
	disk.Sync_action = "none"
	disk.UUID = "UUID209573498"
	disk.Degraded = 0
	disk.Raid_disks = 2
	disk.Array_state = "in_sync"
	disk.Consistency_policy = "bitmap"
	disk.Level = "Raid1"
	disk.Disks = append(disk.Disks, diskInfo{
		Name: "sdc",
		Model: "Samsung SSD 850 PRO",
		Revision: "S1",
		Serial: "DF74JG",
		Size: 104857600,
		SizeShortened: "100 GB",
		State: "active",
		StateIsGood: true,
	})
	disk.Disks = append(disk.Disks, diskInfo{
		Name: "sdd",
		Model: "Samsung SSD 850 PRO",
		Revision: "S1",
		Serial: "DF74J4",
		Size: 104857600,
		SizeShortened: "100 GB",
		State: "failed",
		StateIsGood: false,
	})
	return
}