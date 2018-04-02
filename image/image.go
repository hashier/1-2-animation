package image

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
)

const (
	apart  = 50
	margin = 50
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateCalibrationTestImageAndSave(file string, w, h int) error {
	fmt.Println("Generating calibration file", file)
	calImg := createCalibrationTestImage(w, h)
	err := savePNG(calImg, file)
	return err
}

func createOneColorImage(r, g, b float64, w, h int) image.Image {
	dc := gg.NewContext(w, h)
	dc.SetRGB(r, g, b)
	dc.Clear()
	return dc.Image()
}

func createCalibrationTestImage(width, height int) image.Image {
	dc := gg.NewContext(width, height)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	lw := 1.0
	for i := margin; i <= width; i += 2 * apart {
		x := float64(i)
		dc.DrawLine(x+0.5, 0, x+0.5, float64(height))
		dc.SetLineWidth(lw)
		dc.Stroke()
		lw += 2
	}
	lw = 2.0
	for i := 2 * margin; i <= width; i += 2 * apart {
		x := float64(i)
		dc.DrawLine(x, 0, x, float64(height))
		dc.SetLineWidth(lw)
		dc.Stroke()
		lw += 2
	}

	return dc.Image()
}

func createImage(images []image.Image, ppf int) image.Image {
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
			if ppfCounter == ppf {
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

func ExampleColorImages(w, h, ppf int) {
	var images []image.Image
	var img image.Image

	// gbryg
	images = append(images, createOneColorImage(0, 1, 0, w, h))
	images = append(images, createOneColorImage(0, 0, 1, w, h))
	images = append(images, createOneColorImage(1, 0, 0, w, h))
	images = append(images, createOneColorImage(1, 1, 0, w, h))
	images = append(images, createOneColorImage(0.5, 0.5, 0.5, w, h))
	fmt.Println("Generating 5 frame color example image")
	img = createImage(images, ppf)
	check(savePNG(img, "example-color-5-out.png"))

	// gbrygtp
	images = append(images, createOneColorImage(0, 1, 1, w, h))
	images = append(images, createOneColorImage(1, 0, 1, w, h))
	fmt.Println("Generating 7 frame color example image")
	img = createImage(images, ppf)
	check(savePNG(img, "example-color-7-out.png"))
}

func GenerateImageFrom(files []string, file string, ppf int) {
	fmt.Println("Loading input files")
	images := make([]image.Image, 0)
	for _, imgFile := range files {
		img, err := loadImage(imgFile)
		check(err)
		images = append(images, img)
	}
	outImg := createImage(images, ppf)
	check(savePNG(outImg, file))
}
