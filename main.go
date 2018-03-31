package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
)

const (
	maxWidth  = 1000
	maxHeight = 1000
	apart     = 50
	margin    = 50
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateCalibrationTestImage() image.Image {
	dc := gg.NewContext(maxWidth, maxHeight)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	w := 1.0
	for i := margin; i <= maxWidth; i += 2 * apart {
		x := float64(i)
		dc.DrawLine(x+0.5, 0, x+0.5, maxHeight)
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 2
	}
	w = 2.0
	for i := 2 * margin; i <= maxWidth; i += 2 * apart {
		x := float64(i)
		dc.DrawLine(x, 0, x, maxHeight)
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 2
	}

	return dc.Image()
}

func createImage(images []image.Image, ppf int) image.Image {
	ppfCounter := 0
	index := 0
	img := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))

	for y := 0; y < maxHeight; y++ {
		for x := 0; x < maxWidth; x++ {
			pixel := images[index].At(x, y)
			img.Set(x, y, pixel)
			// log.Printf("Writing at %d/%d colour: %v from index: %d", y, x, pixel, index)
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

func main() {
	calImg := generateCalibrationTestImage()
	check(savePNG(calImg, "/Users/chl/Desktop/1-2-calibration.png"))

	var images []image.Image
	eins, _ := loadImage("/Users/chl/Desktop/red.png")
	zwei, _ := loadImage("/Users/chl/Desktop/green.png")
	drei, _ := loadImage("/Users/chl/Desktop/blue.png")
	vier, _ := loadImage("/Users/chl/Desktop/grey.png")
	funf, _ := loadImage("/Users/chl/Desktop/pink.png")
	sech, _ := loadImage("/Users/chl/Desktop/turk.png")
	sieb, _ := loadImage("/Users/chl/Desktop/yellow.png")
	images = append(images, eins, zwei, drei, vier, funf, sech, sieb)

	img := createImage(images, 10)
	check(savePNG(img, "/Users/chl/Desktop/1-2-out.png"))
}
