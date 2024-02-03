package main

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "math"
    "os"
)

func main() {
    // Load the image
    inputFile := "input.jpg" // Change this to your input image file
    img, err := loadImage(inputFile)
    if err != nil {
        fmt.Println("Error loading image:", err)
        return
    }

    // Apply grayscale conversion
    grayscaleImg := grayscale(img)

    // Apply edge detection using convolution
    edgeImg := applyConvolution(grayscaleImg, getEdgeDetectionKernel())

    // Convert the resulting image array back to an image
    outputImage := image.NewRGBA(grayscaleImg.Bounds())
    draw.Draw(outputImage, outputImage.Bounds(), edgeImg, image.Point{}, draw.Src)

    // Save the resulting image
    outputFile := "output.jpg" // Change this to your output image file
    saveImage(outputFile, outputImage)
    fmt.Println("Edge detection applied and saved to", outputFile)
}

// Load an image from file
func loadImage(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }

    return img, nil
}

// Save an image to file
func saveImage(filename string, img image.Image) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    err = image.JPEG.Encode(file, img, nil)
    if err != nil {
        return err
    }

    return nil
}

// Convert an image to grayscale
func grayscale(img image.Image) *image.Gray {
    bounds := img.Bounds()
    gray := image.NewGray(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            grayColor := color.GrayModel.Convert(img.At(x, y))
            gray.Set(x, y, grayColor.(color.Gray))
        }
    }

    return gray
}

// Apply convolution to an image using a given kernel
func applyConvolution(img image.Image, kernel [][]float64) *image.Gray {
    bounds := img.Bounds()
    gray := image.NewGray(bounds)
    width, height := bounds.Dx(), bounds.Dy()

    for y := 1; y < height-1; y++ {
        for x := 1; x < width-1; x++ {
            var sum float64
            for ky := -1; ky <= 1; ky++ {
                for kx := -1; kx <= 1; kx++ {
                    pixel := color.GrayModel.Convert(img.At(x+kx, y+ky)).(color.Gray)
                    sum += float64(pixel.Y) * kernel[ky+1][kx+1]
                }
            }
            grayColor := color.Gray{Y: uint8(math.Abs(sum))}
            gray.Set(x, y, grayColor)
        }
    }

    return gray
}

// Define a kernel for edge detection
func getEdgeDetectionKernel() [][]float64 {
    return [][]float64{
        {-1, -1, -1},
        {-1, 8, -1},
        {-1, -1, -1},
    }
}
