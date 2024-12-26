// cbz_generator.go
package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/fogleman/gg"
)

const (
	A3WidthPx       = 350
	A3HeightPx      = 496
	DefaultFontSize = 72
	DefaultBorder   = 50
)

func createSampleImage(filename string, pageNum int, fontSize, border int) error {
	img := image.NewRGBA(image.Rect(0, 0, A3WidthPx, A3HeightPx))
	bgColor := color.White
	borderColor := color.Black

	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw border
	for x := border; x < A3WidthPx-border; x++ {
		for y := border; y < A3HeightPx-border; y++ {
			if x == border || x == A3WidthPx-border-1 || y == border || y == A3HeightPx-border-1 {
				img.Set(x, y, borderColor)
			}
		}
	}

	dc := gg.NewContextForRGBA(img)
	dc.SetColor(color.Black)
	dc.LoadFontFace("arial", float64(fontSize))

	// Add main text
	mainText := fmt.Sprintf("Sample Page %d", pageNum)
	dc.DrawStringAnchored(mainText, float64(A3WidthPx)/2, float64(A3HeightPx)/2-50, 0.5, 0.5)

	// Add subtitle
	resolutionText := fmt.Sprintf("Resolution: %dx%d", A3WidthPx, A3HeightPx)
	dc.DrawStringAnchored(resolutionText, float64(A3WidthPx)/2, float64(A3HeightPx)/2+50, 0.5, 0.5)

	outputFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return jpeg.Encode(outputFile, img, nil)
}

func createCBZ(outputFile string, numPages, fontSize, border int) error {
	zipFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for i := 1; i <= numPages; i++ {
		filename := fmt.Sprintf("page_%03d.jpg", i)
		if err := createSampleImage(filename, i, fontSize, border); err != nil {
			return err
		}

		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer os.Remove(filename) // Clean up after
		defer file.Close()

		zipEntry, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(zipEntry, file); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	numPages := flag.Int("pages", 10, "Number of pages to generate")
	outputFile := flag.String("output", "output.cbz", "Name of the output CBZ file")
	fontSize := flag.Int("fontsize", DefaultFontSize, "Font size of the text")
	border := flag.Int("border", DefaultBorder, "Width of the border in pixels")

	flag.Parse()

	if err := createCBZ(*outputFile, *numPages, *fontSize, *border); err != nil {
		log.Fatalf("Failed to create CBZ file: %v", err)
	}
	log.Printf("Successfully created %s", *outputFile)
}
