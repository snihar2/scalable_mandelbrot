package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"math/rand"
	"os"
	"time"
)

const (
	maxEsc = 30
	rMin   = -2.
	rMax   = .5
	iMin   = -1.
	iMax   = 1.
	width  = 4000
)

var palette []color.RGBA
var escapeColor color.RGBA

func mandelbrot(a complex128) int {
	z := a
	i := 0
	// must end absolute value of z is superior to 2 or we attain n max
	for ; i < maxEsc - 1 ; i++ {
		if cmplx.Abs(z) > 2 {
			return i
		}
		z = cmplx.Pow(z, 2) + a
	}
	return maxEsc-i
}

func main() {
	log.Printf("Starting ... \n")

	done := make(chan struct{})
	ticker := time.NewTicker(time.Millisecond * 1000)

	go func() {
		i := 0
		for {
			select {
			case <-ticker.C:
				if i % 15 == 0 {
					fmt.Print(".")
				}
				i++
			case <-done:
				ticker.Stop()
				fmt.Printf("\n\nMandelbrot set rendered into `%s` in %d seconds\n", "mandelbrot_.png", i)
				log.Printf("End 2 ... \n")

			}
		}
	}()

	render(done)
	log.Printf("End 1 ... \n")


}

func render(done chan struct{}) {
	palette = make([]color.RGBA, maxEsc)
	for i := 0; i < maxEsc-1; i++ {
		palette[i] = color.RGBA{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			255}
	}
	escapeColor = color.RGBA{0, 0, 0, 0}

	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)

	b := image.NewNRGBA(bounds)
	// draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			coord := complex(float64(x)/scale+rMin, float64(y)/scale+iMin)
			fEsc := mandelbrot(coord)
			if fEsc == maxEsc - 1 {
				b.Set(x, y, escapeColor)
			}
			b.Set(x, y, palette[fEsc])

		}
	}
	done <- struct{}{}
	f, _ := os.Create("vertical_scalability/mandelbrot_.png")

	err := png.Encode(f, b)

	if err != nil {
		log.Println("png.Encode:", err)
	}

	if err = png.Encode(f, b); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}
