package profile

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

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

func CPUProfileStop() {
	pprof.StopCPUProfile()
}
