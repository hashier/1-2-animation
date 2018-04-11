// Package profile provides functions to start/stop profiling of go apps
package profile

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

// MemProfile runs the GC and dump memory profile
func MemProfile(file string) {
	fmt.Println("Writing memory profile to", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}

// CPUProfile starts CPU profiling and save the result to file
func CPUProfile(file string) {
	fmt.Println("Starting CPU profiling to file", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
}

// CPUProfileStop stops CPU profiling
func CPUProfileStop() {
	pprof.StopCPUProfile()
}
