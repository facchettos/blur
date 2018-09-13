package main

import (
	"os"
	"image/jpeg"
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
	"strconv"
)

var wg sync.WaitGroup

func copyExtremities(img *image.NRGBA, original *image.Image) {
	boundY := (*original).Bounds().Dy()
	boundX := (*original).Bounds().Dx()
	for i := 0; i < img.Bounds().Dx(); i++ {
		for j:=0;j<3 ;j++ {
			img.Set(i,j,(*original).At(i,j))
			img.Set(i,boundY-j,(*original).At(i,boundY-j))
		}
	}

	for i := 0; i < img.Bounds().Dy(); i++ {
		for j:=0;j<3 ;j++ {
			img.Set(j,i,(*original).At(i,j))
			img.Set(boundX-j,i,(*original).At(boundX-i,j))
		}
	}

}

func worker(img *image.NRGBA, total int, rank int, original *image.Image) {
	rangee := (*original).Bounds().Dx() / total
	start := rangee * rank
	var effectiveStart int
	if rangee == 0 {
		effectiveStart=start+3
	}else{
		effectiveStart=start
	}
	var nr, ng, nb = 0.0, 0.0, 0.0
	for i := effectiveStart; i < start+rangee; i++ {
		for j := 3; j < (*original).Bounds().Dy()-2; j++ {
			nr, nb, ng = 0.0, 0.0, 0.0
			for k := -2; k < 3; k++ {
				for l := -2; l < 3; l++ {
					r, g, b, a := (*original).At(i+k, j+l).RGBA()
					nr += (float64(r) / float64(a) * 255 / 25.0)
					ng += (float64(g) / float64(a) * 255 / 25.0)
					nb += (float64(b) / float64(a) * 255 / 25.0)
				}
			}
			(*img).Set(i, j, color.RGBA{uint8(nr), uint8(ng), uint8(nb), 255})
		}
	}
	wg.Done()
	return
}

func main() {

	outputFile := os.Args[2]
	numberofcore, err := strconv.Atoi(os.Args[3])
	if err != nil {
		wg.Add(1)
	} else {
		wg.Add(numberofcore)
	}
	filePath := os.Args[1]

	fmt.Println(filePath)

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println("error opening the file")
	}

	myImage, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("error decoding image")
		return
	}

	//var nr, ng, nb = 0, 0, 0
	blurred := image.NewNRGBA(myImage.Bounds())
	start := time.Now()
	for i := 0; i < numberofcore; i++ {
		go worker(blurred, numberofcore, i, &myImage)
	}



	var opt = jpeg.Options{}
	opt.Quality = 90
	out, err := os.Create(outputFile)
	wg.Wait()
	copyExtremities(blurred,&myImage)
	fmt.Println(time.Since(start))
	jpeg.Encode(out, blurred, &opt)
	fmt.Println("finished")

}
