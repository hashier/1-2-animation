package main

import "github.com/fogleman/gg"

const (
	maxWidth  = 1680
	maxHeight = 1000
)

func main() {
	dc := gg.NewContext(maxWidth, maxHeight)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	w := 1.0
	for i := 50; i <= maxWidth; i += 100 {
		x := float64(i)
		dc.DrawLine(x+0.5, 0, x+0.5, maxHeight)
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 2
	}
	w = 2.0
	for i := 100; i <= maxWidth; i += 100 {
		x := float64(i)
		dc.DrawLine(x, 0, x, maxHeight)
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 2
	}

	dc.SavePNG("/Users/chl/Desktop/out.png")
}
