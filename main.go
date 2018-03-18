package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	path := os.Args[1]
	scale, _ := strconv.Atoi(os.Args[2])

	img, err := getImage(path)
	if err != nil {
		log.Fatal(err)
	}

	newImg, err := pixelate(img, scale)
	if err != nil {
		log.Fatal(err)
	}

	newPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".pxl" + filepath.Ext(path)
	saveImage(newImg, newPath)
}

func getImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func saveImage(img image.Image, path string) {
	file, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	jpeg.Encode(file, img, nil)
}

func pixelate(img image.Image, scale int) (image.Image, error) {
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y
	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x += scale {
		for y := 0; y < height; y += scale {
			averageColor := computeAverageColor(img, x, y, scale)

			for i := x; i < x+scale; i++ {
				for j := y; j < y+scale; j++ {
					newImg.Set(i, j, averageColor)
				}
			}
		}
	}

	return newImg, nil
}

func computeAverageColor(img image.Image, x, y, scale int) color.Color {
	numColors := uint32(scale * scale)
	var r, g, b, a uint32

	for i := x; i < x+scale; i++ {
		for j := y; j < y+scale; j++ {
			rI, gI, bI, aI := img.At(i, j).RGBA()
			r += rI
			g += gI
			b += bI
			a += aI
		}
	}

	avgR := uint16(r / numColors)
	avgG := uint16(g / numColors)
	avgB := uint16(b / numColors)
	avgA := uint16(a / numColors)

	return color.NRGBA64{avgR, avgG, avgB, avgA}
}
