package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
)

// squareSize is a square size in 2D grid
const squareSize = 400

var (
	numOfSquares int
	imagesFilled [][]bool
	dimensions   map[int]Coordinates
)

// MatrixInfo represents 2D matrix
type MatrixInfo struct {
	Matrix [][]bool
}

// calculateImageMatrix creates GetImageResponse based on provided parameters
func calculateImageMatrix(x, y, z int) (*GetImagesResponse, error) {
	dimensions = generateDimensions()
	grid := createMatrix(dimensions[z]).Matrix

	for i := 0; i < numOfSquares; i++ {
		grid[i] = make([]bool, dimensions[z].X/squareSize)
	}

	imagesFilled = grid
	var imagesResponse []Image
	rand.Seed(5) // set seed here to always get same indexes

	for !checkLastRow() {
		image, err := getImage(imagesFilled)
		if err != nil {
			return nil, fmt.Errorf("getImage error: %s", err)
		}
		imagesResponse = append(imagesResponse, *image)
	}

	return &GetImagesResponse{Images: imagesResponse}, nil
}

// getImage returns Image struct with calculated position
func getImage(currentGrid [][]bool) (*Image, error) {
	files, err := ioutil.ReadDir(imagesDir + "/")
	if err != nil {
		return nil, fmt.Errorf("could not read dir: %s", err)
	}

	numOfFiles := len(files)

	index := rand.Intn(numOfFiles)
	file := files[index]

	cfg, err := getImageDimensions(imagesDir + "/" + file.Name())
	if err != nil {
		return nil, fmt.Errorf("could not get image dimensions: %s, error: %s", file.Name(), err)
	}

	x, y := howManySquares(Coordinates{X: cfg.Width, Y: cfg.Height})
	xIndexes, yIndexes, err := fillGrid(currentGrid, x, y)
	if err != nil {
		return nil, fmt.Errorf("fillGrid error: %s", err)
	}

	desiredPosition, err := calculatePosition(xIndexes, yIndexes)
	if err != nil {
		return nil, fmt.Errorf("calculatePosition error: %s", err)
	}

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
		Name:     file.Name(),
		Position: *desiredPosition,
		Dimensions: Coordinates{
			X: cfg.Width,
			Y: cfg.Height,
		},
		Base64: base64Encoding,
	}

	return image, nil
}

// fillGrid fills the squares in matrix based on provided number of squares used by an image. Function
// return indexes in the matrix, which have been filled.
func fillGrid(currentGrid [][]bool, x, y int) ([]int, []int, error) {
	var filledIndexesX, filledIndexesY []int

	for i := 0; i < numOfSquares; i++ {
		for j := 0; j < numOfSquares; j++ {
			if !currentGrid[j][i] {
				for l := i; l < y+i; l++ {
					if l > numOfSquares-1 {
						continue
					}
					for k := j; k < x+j; k++ {
						if k > numOfSquares-1 {
							continue
						}
						currentGrid[k][l] = true
						filledIndexesX = append(filledIndexesX, k)
					}
					filledIndexesY = append(filledIndexesY, l)
				}
				return filledIndexesX, filledIndexesY, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("could not fill the matrix") // if error return indexes, not error
}

// checkLastRow() checks whether the last row of a matrix is filled with images
func checkLastRow() bool {
	for i := 0; i < numOfSquares; i++ {
		if !imagesFilled[numOfSquares-1][i] {
			return false
		}
	}

	return true
}

// createMatrix creates a matrix based on received position
func createMatrix(dim Coordinates) *MatrixInfo {
	num := dim.X / squareSize
	matrix := make([][]bool, num)
	numOfSquares = num

	return &MatrixInfo{
		Matrix: matrix,
	}
}

// calculatePosition calculates position for given image based on indexes that this image occupy
func calculatePosition(x, y []int) (*Coordinates, error) {
	sort.Ints(x)
	sort.Ints(y)
	if len(x) == 0 || len(y) == 0 {
		return nil, fmt.Errorf("there are not indexes to be used as coordinates")
	}

	// at first index will be the minimum
	return &Coordinates{
		X: x[0] * 400,
		Y: y[0] * 400,
	}, nil
}

// howManySquares counts how many squares given image takes up horizontally and vertically
func howManySquares(dimensions Coordinates) (int, int) {
	var xSquares, ySquares int

	xRatio := float64(dimensions.X) / float64(squareSize)
	yRatio := float64(dimensions.Y) / float64(squareSize)

	intPartX := int(xRatio)
	decimalPartX := xRatio - float64(intPartX)
	if decimalPartX < 0.5 {
		xSquares = intPartX
	} else {
		xSquares = intPartX + 1
	}

	intPartY := int(yRatio)
	decimalPartY := yRatio - float64(intPartY)
	if decimalPartY < 0.5 {
		ySquares = intPartY
	} else {
		ySquares = intPartY + 1
	}

	return xSquares, ySquares
}

// generateDimensions generates dimensions map which specifies the dimension range based on z parameter
func generateDimensions() map[int]Coordinates {
	dimensionsMap := make(map[int]Coordinates, 100)
	for i := 0; i < 100; i++ {
		dimensionsMap[i] = Coordinates{
			X: 2000 + i*squareSize,
			Y: 2000 + i*squareSize,
		}
	}

	return dimensionsMap
}
