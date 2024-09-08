package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	width, height := 300, 300
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Отрисовка простого вертикального паттерна
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x%50 < 25 {
				img.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255}) // Blue
			}
		}
	}

	file, err := os.Create("amazing_logo.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
