package main

import (
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/jpeg"
    "math"
    "os"
    "sync"
    "image/draw"
)


func main() {
    var inputFile, outputFile string
    var jpegQuality int
    flag.StringVar(&inputFile, "input", "input.jpg", "Input image file")
    flag.StringVar(&outputFile, "output", "output.jpg", "Output image file")
    flag.IntVar(&jpegQuality, "quality", 75, "JPEG quality for the output image")
    flag.Parse()

    img, err := loadImage(inputFile)
    if err != nil {
        fmt.Printf("Error loading image '%s': %v\n", inputFile, err)
        return
    }

    grayscaleImg := grayscale(img)
    edgeImg := applyConvolution(grayscaleImg, getEdgeDetectionKernel())

    outputImage := image.NewRGBA(grayscaleImg.Bounds())
    draw.Draw(outputImage, outputImage.Bounds(), edgeImg, image.Point{}, draw.Src)

    if err := saveImage(outputFile, outputImage, jpegQuality); err != nil {
        fmt.Printf("Error saving image '%s': %v\n", outputFile, err)
        return
    }

    fmt.Println("Edge detection applied and saved to", outputFile)
}

func loadImage(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    return img, err
}

func saveImage(filename string, img image.Image, quality int) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
}

func grayscale(img image.Image) *image.Gray {
    bounds := img.Bounds()
    gray := image.NewGray(bounds)

    var wg sync.WaitGroup
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        wg.Add(1)
        go func(y int) {
            defer wg.Done()
            for x := bounds.Min.X; x < bounds.Max.X; x++ {
                grayColor := color.GrayModel.Convert(img.At(x, y))
                gray.Set(x, y, grayColor.(color.Gray))
            }
        }(y)
    }
    wg.Wait()

    return gray
}

func applyConvolution(img *image.Gray, kernel [][]float64) *image.Gray {
    bounds := img.Bounds()
    newImg := image.NewGray(bounds)
    width, height := bounds.Dx(), bounds.Dy()

    var wg sync.WaitGroup
    for y := 1; y < height-1; y++ {
        wg.Add(1)
        go func(y int) {
            defer wg.Done()
            for x := 1; x < width-1; x++ {
                var sum float64
                for ky := -1; ky <= 1; ky++ {
                    for kx := -1; kx <= 1; kx++ {
                        pixel := img.GrayAt(x+kx, y+ky)
                        sum += float64(pixel.Y) * kernel[ky+1][kx+1]
                    }
                }
                newImg.SetGray(x, y, color.Gray{Y: uint8(math.Abs(sum))})
            }
        }(y)
    }
    wg.Wait()

    return newImg
}

func getEdgeDetectionKernel() [][]float64 {
    return [][]float64{
        {-1, -1, -1},
        {-1, 8, -1},
        {-1, -1, -1},
    }
}
