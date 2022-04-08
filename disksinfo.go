package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// for info on the values these can have, check out:
// https://www.kernel.org/doc/html/v4.15/admin-guide/md.html

type diskInfo struct {
	Name string // /sys/block/mdX/md/dev-'name'
	Size int64  // /sys/block/mdX/md/dev-name/size
	SizeShortened string
	Model string // /sys/block/mdX/md/dev-name/block/device/model
	Revision string // /sys/block/mdX/md/dev-name/block/device/rev
	Serial string // /sys/block/mdX/md/dev-name/block/device/vpd_pg80
	State string // /sys/block/mdX/md/dev-name/state
	StateIsGood bool // true if state of drive is usable and working
}

type mdInfo struct {
	Disks []diskInfo

	Name string
	Sync_action string // /sys/block/mdX/md/sync_action
	UUID string // /sys/block/mdX/md/uuid
	ArrayIsGood bool // true if array is healthy with no failed drive
	ArrayIsDegraded bool // true if array is working but one or more drives have failed
	Degraded int // /sys/block/mdX/md/degraded
	Raid_disks int // /sys/block/mdX/md/raid_disks
	Array_state string // /sys/block/mdX/md/array_state
	Consistency_policy string // /sys/block/mdX/md/consistency_policy
	Level string // /sys/block/mdX/md/level
}

func mdDeviceGetInfo(mdname string) (mdinfo mdInfo, err error) {
	err = nil

	mdpath := "/sys/block/" + mdname + "/md/"
	if _, err = os.Stat(mdpath); os.IsNotExist(err) {
		err = errors.New("MD device does not exist")
		return
	}

	mdinfo.Name = mdname

	var val []byte

	// disks
	files, err := ioutil.ReadDir(mdpath)
	if (err != nil) {
		return
	}
	
	for _, f := range files {
		if (strings.HasPrefix(f.Name(), "dev-")) {
			var disk diskInfo
			disk.Name = f.Name()[4:]

			// state
			val, err = ioutil.ReadFile(mdpath + f.Name() + "/state")
			if (err != nil) {
				return
			}
			disk.State = stringStripNewline(string(val))

			disk.StateIsGood = false
			if (disk.State == "in_sync" || disk.State == "writemostly" || disk.State == "spare" || disk.State == "replacement") {
				disk.StateIsGood = true
			}

			// size
			val, err = ioutil.ReadFile(mdpath + f.Name() + "/size")
			if (err != nil) {
				return
			}
			var valI int64
			valI, err = strconv.ParseInt(stringStripNewline(string(val)), 10, 64)
			if (err != nil) {
				return
			}
			disk.Size = valI

			// SizeShortened
			disk.SizeShortened = normalizeKBValue(disk.Size)
			
			// model
			val, err = ioutil.ReadFile(mdpath + f.Name() + "/block/device/model")
			disk.Model = stringStripNewline(string(val))
			if (err != nil) {
				disk.Model = "Unknown model"
			}

			// revision
			val, err = ioutil.ReadFile(mdpath + f.Name() + "/block/device/rev")
			disk.Revision = stringStripNewline(string(val))
			if (err != nil) {
				disk.Revision = "Unknown rev"
			}

			// serial
			val, err = ioutil.ReadFile(mdpath + f.Name() + "/block/device/vpd_pg80")
			for _, c := range stringStripNewline(string(val)) {
				if (c >= 33 && c <= 126) {
					disk.Serial += string(c)
				}
			}
			if (err != nil || disk.Serial == "") {
				disk.Serial = "Unknown serial"
			}

			mdinfo.Disks = append(mdinfo.Disks, disk)
		}
	}

	// sync_action
	val, err = ioutil.ReadFile(mdpath + "sync_action")
	mdinfo.Sync_action = stringStripNewline(string(val))
	if (err != nil) {
		mdinfo.Sync_action = ""
	}

	// uuid
	val, err = ioutil.ReadFile(mdpath + "uuid")
	mdinfo.UUID = stringStripNewline(string(val))
	if (err != nil) {
		mdinfo.UUID = ""
	}

	// degraded
	val, err = ioutil.ReadFile(mdpath + "degraded")
	if (err != nil) {
		return
	}
	valI, err := strconv.Atoi(stringStripNewline(string(val)))
	if (err != nil) {
		return
	}
	mdinfo.Degraded = valI

	// raid_disks
	val, err = ioutil.ReadFile(mdpath + "raid_disks")
	if (err != nil) {
		return
	}
	valI, err = strconv.Atoi(stringStripNewline(string(val)))
	mdinfo.Raid_disks = valI
	if (err != nil) {
		mdinfo.Raid_disks = 0
	}

	// array_state
	val, err = ioutil.ReadFile(mdpath + "array_state")
	if (err != nil) {
		return
	}
	mdinfo.Array_state = stringStripNewline(string(val))

	// ArrayIsGood
	mdinfo.ArrayIsGood = false
	if (mdinfo.Array_state == "clean" || mdinfo.Array_state == "active" || mdinfo.Array_state == "active_idle"){
		mdinfo.ArrayIsGood = true
	}

	// ArrayIsDegraded
	mdinfo.ArrayIsDegraded = false
	if (mdinfo.Array_state == "clear" || mdinfo.Array_state == "inactive" || mdinfo.Array_state == "readonly"){
		mdinfo.ArrayIsDegraded = true
	}
	if (mdinfo.Degraded != 0) {
		mdinfo.ArrayIsGood = false
		mdinfo.ArrayIsDegraded = false
		// both false means failure
	}

	// consistency_policy
	val, err = ioutil.ReadFile(mdpath + "consistency_policy")
	mdinfo.Consistency_policy = stringStripNewline(string(val))
	if (err != nil) {
		mdinfo.Consistency_policy = ""
	}

	// level
	val, err = ioutil.ReadFile(mdpath + "level")
	if (err != nil) {
		return
	}
	mdinfo.Level = stringStripNewline(string(val))

	return
}

// This function reuses the diskInfo struct for md drive devices
func diskGetInfo(name string) (disk diskInfo, err error){
	disk.Name = name

	basePath := "/sys/block/" + name + "/"

	// state
	val, err := ioutil.ReadFile(basePath + "device/" + "state")
	if (err != nil) {
		return disk, err
	}
	disk.State = stringStripNewline(string(val))

	disk.StateIsGood = false
	if (disk.State == "running") {
		disk.StateIsGood = true
	}

	// size
	val, err = ioutil.ReadFile(basePath + "size")
	if (err != nil) {
		return disk, err
	}
	var valI int64
	valI, err = strconv.ParseInt(stringStripNewline(string(val)), 10, 64)
	if (err != nil) {
		return disk, err
	}
	disk.Size = valI

	// SizeShortened
	disk.SizeShortened = normalizeKBValue(disk.Size/2.0)// /2.0 because size is in 512 byte sectors
	
	// model
	val, err = ioutil.ReadFile(basePath + "device/model")
	disk.Model = stringStripNewline(string(val))
	if (err != nil) {
		disk.Model = "Unknown model"
	}

	// revision
	val, err = ioutil.ReadFile(basePath + "device/rev")
	disk.Revision = stringStripNewline(string(val))
	if (err != nil) {
		disk.Revision = "Unknown rev"
	}

	// serial
	val, err = ioutil.ReadFile(basePath + "device/vpd_pg80")
	for _, c := range stringStripNewline(string(val)) {
		if (c >= 33 && c <= 126) {
			disk.Serial += string(c)
		}
	}
	if (err != nil || disk.Serial == "") {
		disk.Serial = "Unknown serial"
	}

	return disk, err
}

func findStorageDevicesInSystem() (disks []string, arrays []string) {
	files, err := ioutil.ReadDir("/sys/block/")
	if (err != nil) {
		return
	}
	
	for _, f := range files {
		if (strings.HasPrefix(f.Name(), "sd")) {
			disks = append(disks, f.Name())
		}else if (strings.HasPrefix(f.Name(), "nvme")) {
			disks = append(disks, f.Name())
		}else if (strings.HasPrefix(f.Name(), "md")) {
			arrays = append(arrays, f.Name())
		}
	}
	return
}