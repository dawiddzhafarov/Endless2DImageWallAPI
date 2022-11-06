package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
)

var dimensions = map[int]Coordinates{
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

type MatrixInfo struct {
	Matrix [][]bool
}

func calculateImageMatrix(x, y, z int) (*GetImagesResponse, error) {
	//dim := dimensions[z]
	//matrixInfo := createMatrix(dim)
	//grid := matrixInfo.Matrix

	grid := createMatrix(dimensions[z]).Matrix

	for i := 0; i < numOfSquares; i++ {
		grid[i] = make([]bool, dimensions[z].X/squareSize)
	}

	imagesFilled = grid
	var imagesResponse []Image

	for !checkLastRow() {
		image, err := getImage(imagesFilled)
		if err != nil {
			return nil, fmt.Errorf("getImage error: %s", err)
		}
		imagesResponse = append(imagesResponse, *image)
	}

	return &GetImagesResponse{Images: imagesResponse}, nil
}

func getImage(currentGrid [][]bool) (*Image, error) {
	files, err := ioutil.ReadDir(imagesDir + "/")
	if err != nil {
		return nil, fmt.Errorf("could not read dir: %s", err)
	}

	numOfFiles := len(files)
	//rand.Seed(5)

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

	desiredPosition := calculatePositions(xIndexes, yIndexes) // TODO

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
		Dimensions: Coordinates{
			cfg.Width,
			cfg.Height,
		},
		Base64: "", //base64Encoding,
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

func checkLastRow() bool {
	for i := 0; i < numOfSquares; i++ {
		if !imagesFilled[numOfSquares-1][i] {
			return false
		}
	}

	return true
}

func createMatrix(dim Coordinates) *MatrixInfo {
	num := dim.X / squareSize
	matrix := make([][]bool, num)
	numOfSquares = num
	return &MatrixInfo{
		Matrix: matrix,
	}
	//switch dim.X {
	//case 2000:
	//	const numberOfSquares = 2000 / 400
	//	numOfSquares = numberOfSquares
	//	imagesPlain := make([][]bool, numberOfSquares)
	//	return &plainInfo{
	//		plain:           imagesPlain,
	//		numberOfSquares: numberOfSquares}
	//
	//}
}

func calculatePositions(x, y []int) Coordinates {
	sort.Ints(x)
	sort.Ints(y)
	if len(x) == 0 || len(y) == 0 {
		fmt.Println("ZERO HEHE")
	}
	minX := x[0]
	mixY := y[0]
	return Coordinates{
		X: minX * 400,
		Y: mixY * 400,
	}
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
