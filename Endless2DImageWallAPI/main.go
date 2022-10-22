package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	handler := http.HandlerFunc(GetImages)

	http.Handle("/image", handler)
	fmt.Println("Server started at port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return
	}
}

// GetImagesResponse contains Images with metadata
type GetImagesResponse struct {
	Images []Image `json:"images"`
}

// Image contains base64 encoded image with metadata
type Image struct {
	ID         int      `json:"id"`
	Position   Position `json:"position"`
	Dimensions Position `json:"dimensions"`
	Base64     string   `json:"base64"`
}

// Position describes image postition or dimensions
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const (
	imagesDir  = "./images"
	base64jpeg = "data:image/jpeg;base64,"
	base64png  = "data:image/png;base64,"
)

func GetImages(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(imagesDir)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := createResponse(files)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatal(err)
	}
}

// toBase64 encodes bytes as base64 string
func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// createResponse creates JSON response
func createResponse(files []fs.FileInfo) (*GetImagesResponse, error) {
	var imagesList []Image
	for i, file := range files {
		img, err := ioutil.ReadFile(imagesDir + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		cfg, err := getImageDimensions(imagesDir + "/" + file.Name())
		if err != nil {
			return nil, err
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

		imagesList = append(imagesList, Image{
			ID: i,
			Position: Position{
				X: i + 100,
				Y: i - 100,
			},
			Dimensions: Position{
				X: cfg.Width,
				Y: cfg.Height,
			},
			Base64: base64Encoding,
		})
	}
	return &GetImagesResponse{Images: imagesList}, nil
}

// getImageDimensions retrieves image dimensions
func getImageDimensions(filePath string) (*image.Config, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	imgConfig, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	return &imgConfig, nil
}
