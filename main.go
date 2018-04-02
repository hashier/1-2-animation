package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hashier/1-2-animation/image"
	"github.com/hashier/1-2-animation/profile"
)

var (
	config = configuration{
		calibrateFile:  "calibration.png",
		cpuProfileFile: "cpu.prof",
		memProfileFile: "mem.prof",
	}
)

type configuration struct {
	outputFile        string
	calibrateFile     string
	exampleColorImage bool
	cpuProfile        bool
	cpuProfileFile    string
	memProfile        bool
	memProfileFile    string
	frames            int
	ppf               int
	width             int
	height            int
}

func init() {
	flag.StringVar(&config.outputFile, "o", "out.png", "output of generated PNG image")
	flag.StringVar(&config.calibrateFile, "calibrate", "", "Write calibration image to <file>. This image can be used to determine how many pixels the \"slit\" is wide")
	flag.BoolVar(&config.exampleColorImage, "example", false, "Write some demo color striped images to the drive. Those can be helpful to determine to \"number of frames\"")
	flag.BoolVar(&config.cpuProfile, "cp", false, "write CPU profile to cpu.prof")
	flag.BoolVar(&config.memProfile, "mp", false, "write memory profile to mem.prof")
	flag.IntVar(&config.frames, "f", 5, "How many input frames (set automatically if images provide, needed for example images)")
	flag.IntVar(&config.ppf, "ppf", 2, "How many pixel to take from the input image for every frame (how wide is the \"slit\" of the mask")
	flag.IntVar(&config.width, "w", 1680, "Width of image, only applies to calibration and example images")
	flag.IntVar(&config.height, "h", 1200, "Height of image, only applies to calibration and example images")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if config.cpuProfile {
		profile.CPUProfile(config.cpuProfileFile)
	}

	if config.calibrateFile != "" {
		check(image.GenerateCalibrationTestImageAndSave(config.calibrateFile, config.width, config.height))
	}

	if config.exampleColorImage {
		image.ExampleColorImages(config.width, config.height, config.ppf)
	}

	if len(flag.Args()) >= 2 {
		image.GenerateImageFrom(flag.Args(), config.outputFile, config.ppf)
	} else if !config.exampleColorImage && config.calibrateFile == "" {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		log.Fatal("Error: You need to either provide at least 2 input image OR enable calibration mode OR example mode")
	}

	if config.cpuProfile {
		profile.CPUProfileStop()
	}

	if config.memProfile {
		profile.MemProfile(config.memProfileFile)
	}
}
