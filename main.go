package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"runtime/pprof"

	"github.com/fogleman/gg"
)

const (
	apart  = 50
	margin = 50
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
	flag.IntVar(&config.width, "x", 1680, "Width of image, only applies to calibration and example images")
	flag.IntVar(&config.height, "h", 1200, "Height of image, only applies to calibration and example images")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateCalibrationTestImageAndSave() {
	calImg := createCalibrationTestImage()
	err := savePNG(calImg, config.calibrateFile)
	check(err)
}

func createOneColorImage(r, g, b float64) image.Image {
	dc := gg.NewContext(config.width, config.height)
	dc.SetRGB(r, g, b)
	dc.Clear()
	return dc.Image()
}

func createCalibrationTestImage() image.Image {
	fmt.Println("Generating calibration file", config.calibrateFile)
	dc := gg.NewContext(config.width, config.height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	w := 1.0
	for i := margin; i <= config.width; i += 2 * apart {
		x := float64(i)
		dc.DrawLine(x+0.5, 0, x+0.5, float64(config.height))
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 2
	}
	w = 2.0
	for i := 2 * margin; i <= config.width; i += 2 * apart {
		x := float64(i)
		dc.DrawLine(x, 0, x, float64(config.height))
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 2
	}

	return dc.Image()
}

func createImage(images []image.Image) image.Image {
	ppfCounter := 0
	index := 0
	width := images[0].Bounds().Size().X
	height := images[0].Bounds().Size().Y

	fmt.Printf("Generating %dx%d pixel image\n", width, height)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := images[index].At(x, y)
			img.Set(x, y, pixel)
			ppfCounter++
			if ppfCounter == config.ppf {
				ppfCounter = 0
				index = (index + 1) % len(images)
			}
		}
		ppfCounter = 0
		index = 0
	}

	return img
}

func savePNG(img image.Image, path string) error {
	fmt.Println("Writing image to", path)
	fd, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := png.Encode(fd, img); err != nil {
		fd.Close()
		return err
	}

	if err := fd.Close(); err != nil {
		return err
	}

	return nil
}

func loadImage(path string) (image.Image, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	im, _, err := image.Decode(fd)
	return im, err
}

func memProfile() {
	fmt.Println("Writing memory profile to", config.memProfileFile)
	f, err := os.Create(config.memProfileFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pprof.WriteHeapProfile(f)
}

func exampleColorImages() {
	var images []image.Image
	var img image.Image

	// gbryg
	images = append(images, createOneColorImage(0, 1, 0))
	images = append(images, createOneColorImage(0, 0, 1))
	images = append(images, createOneColorImage(1, 0, 0))
	images = append(images, createOneColorImage(1, 1, 0))
	images = append(images, createOneColorImage(0.5, 0.5, 0.5))
	fmt.Println("Generating 5 frame color example image")
	img = createImage(images)
	check(savePNG(img, "example-color-5-out.png"))

	// gbrygtp
	images = append(images, createOneColorImage(0, 1, 1))
	images = append(images, createOneColorImage(1, 0, 1))
	fmt.Println("Generating 7 frame color example image")
	img = createImage(images)
	check(savePNG(img, "example-color-7-out.png"))
}

func generateImageFrom(files []string) {
	fmt.Println("Loading input files")
	images := make([]image.Image, 0)
	for _, imgFile := range files {
		img, err := loadImage(imgFile)
		check(err)
		images = append(images, img)
	}
	outImg := createImage(images)
	check(savePNG(outImg, config.outputFile))
}

func main() {
	flag.Parse()

	if config.cpuProfile {
		fmt.Println("Starting CPU profiling to file", config.cpuProfileFile)
		f, err := os.Create(config.cpuProfileFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if config.calibrateFile != "" {
		generateCalibrationTestImageAndSave()
	}

	if config.exampleColorImage {
		exampleColorImages()
	}

	if len(flag.Args()) >= 2 {
		generateImageFrom(flag.Args())
	} else if !config.exampleColorImage && config.calibrateFile == "" {
		log.Fatal("Error: You need to either provide at least 2 input image OR enable calibration mode OR example mode")
	}

	if config.memProfile {
		memProfile()
	}
}
