package main

import (
	"fmt"

	"github.com/mackerelio/go-osstat/memory"
)

func getMemoryUsage() (total float64, used float64, cached float64, free float64) {
	memory, err := memory.Get()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	return float64(memory.Total/1024.0/1024.0), float64(memory.Used/1024.0/1024.0), float64(memory.Cached/1024.0/1024.0), float64(memory.Free/1024.0/1024.0)
}