package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
)

var dimensions map[int]Position = map[int]Position{
	0: {
		X: 2000,
		Y: 2000,
	},
	1: {
		X: 4000,
		Y: 4000,
	},
	2: {
		X: 6000,
		Y: 6000,
	},
	3: {
		X: 8000,
		Y: 8000,
	},
	4: {
		X: 10000,
		Y: 10000,
	},
}

const squareSize = 400

var numOfSquares int

var imagesFilled [][]bool

type plainInfo struct {
	plain           [][]bool
	numberOfSquares int
}

func doSomething() (*GetImagesResponse, error) {
	pos := Position{
		X: 2000,
		Y: 2000,
	}

	plain := createPlain(pos)

	grid := plain.plain

	for i := 0; i < plain.numberOfSquares; i++ {
		grid[i] = make([]bool, pos.X/squareSize)
	}
	imagesFilled = grid

	var imagesResponse []Image
	// while grid not filled
	for !imagesFilled[numOfSquares-1][numOfSquares-1] || !imagesFilled[2][numOfSquares-1] {
		image, err := getImage(imagesFilled) //add to tge img slice
		if err != nil {
			return nil, err
		}
		imagesResponse = append(imagesResponse, *image)
	}

	return &GetImagesResponse{
		Images: imagesResponse,
	}, nil
}

func createPlain(position Position) *plainInfo {
	switch position.X {
	case 2000:
		const numberOfSquares = 2000 / 400
		numOfSquares = numberOfSquares
		imagesPlain := make([][]bool, numberOfSquares)
		return &plainInfo{
			plain:           imagesPlain,
			numberOfSquares: numberOfSquares}

	}
	return nil
}

func fillPlain(plain [][]bool, numOfSquares, x, y int) ([]int, []int) {
	var filledIndexesX []int
	var filledIndexesY []int
	for i := 0; i < numOfSquares; i++ {
		for j := 0; j < numOfSquares; j++ {
			if !plain[j][i] {
				for l := i; l < y+i; l++ {
					if l > numOfSquares-1 {
						continue
					}
					for k := j; k < x+j; k++ {
						if k > numOfSquares-1 {
							continue
						}
						plain[k][l] = true

						filledIndexesX = append(filledIndexesX, k)
					}
					filledIndexesY = append(filledIndexesY, l)
				}
				fmt.Println(plain)
				//fmt.Println(numOfSquares)
				//fmt.Println(filledIndexesX)
				//fmt.Println(filledIndexesY)
				//fmt.Println(imagesFilled)
				return filledIndexesX, filledIndexesY
			}
		}
	}
	return filledIndexesX, filledIndexesY
}

func getImage(plain [][]bool) (*Image, error) {
	files, err := ioutil.ReadDir(imagesDir + "/")
	if err != nil {
		log.Fatal(err)
	}

	numOfFiles := len(files)
	//rand.Seed(5)

	index := rand.Intn(numOfFiles)
	file := files[index]

	cfg, err := getImageDimensions(imagesDir + "/" + file.Name())
	if err != nil {
		return nil, fmt.Errorf("could not get image dimensions: %s, %s", file.Name(), err)
	}

	x, y := howManySquares(Position{X: cfg.Width, Y: cfg.Height})
	xIndexes, yIndexes := fillPlain(plain, numOfSquares, x, y)

	desiredPosition := calculatePositions(xIndexes, yIndexes)

	img, err := ioutil.ReadFile(imagesDir + "/" + file.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read file: %s, %s", file.Name(), err)
	}
	var base64Encoding string
	mimeType := http.DetectContentType(img)

	switch mimeType {
	case "image/jpeg":
		base64Encoding += base64jpeg
	case "image/png":
		base64Encoding += base64png
	}
	base64Encoding += toBase64(img)
	image := &Image{
		ID:       0,
		Position: desiredPosition,
		Dimensions: Position{
			cfg.Width,
			cfg.Height,
		},
		Base64: "", //base64Encoding,
	}
	return image, nil
}

func calculatePositions(x, y []int) Position {
	sort.Ints(x)
	sort.Ints(y)
	if len(x) == 0 || len(y) == 0 {
		fmt.Println("ZERO HEHE")
	}
	minX := x[0]
	mixY := y[0]
	return Position{
		X: minX * 400,
		Y: mixY * 400,
	}
}

func howManySquares(dimensions Position) (int, int) {
	var xSquares int
	var ySquares int

	xRatio := float64(dimensions.X) / float64(squareSize)
	yRatio := float64(dimensions.Y) / float64(squareSize)

	intNumberX := int(xRatio)
	decimalPartX := xRatio - float64(intNumberX)
	if decimalPartX < 0.5 {
		xSquares = intNumberX
	} else {
		xSquares = intNumberX + 1
	}

	intNumberY := int(yRatio)
	decimalPartY := yRatio - float64(intNumberY)
	if decimalPartY < 0.5 {
		ySquares = intNumberY
	} else {
		ySquares = intNumberY + 1
	}

	return xSquares, ySquares
}
