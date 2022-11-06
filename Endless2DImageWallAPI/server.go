package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

// Position describes image position or dimensions
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

	resp, err := doSomething()
	if err != nil {
		log.Fatal("error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Fatal(err)
	}
	//files, err := ioutil.ReadDir(imagesDir)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//resp, err := createResponse(files)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//if err = json.NewEncoder(w).Encode(resp); err != nil {
	//	log.Fatal(err)
	//}
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
			return nil, fmt.Errorf("could not read file: %s, %s", file.Name(), err)
		}

		cfg, err := getImageDimensions(imagesDir + "/" + file.Name())
		if err != nil {
			return nil, fmt.Errorf("could not get image dimensions: %s, %s", file.Name(), err)
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
		return nil, err
	}

	return &imgConfig, nil
}

func basicAuth(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			correctUsername := sha256.Sum256([]byte("sample_user"))
			correctPassword := sha256.Sum256([]byte("sample_password"))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], correctUsername[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], correctPassword[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func resizeImage(oldImage []byte, width, height uint) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(oldImage))
	if err != nil {
		return nil, err
	}

	newImage := resize.Resize(width, height, img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImage, nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//resizedImage, err := resizeImage(img, 100, 100)
//if err != nil {
//return nil, err
//}
//err = os.WriteFile("dat1.jpg", resizedImage, 0644)
//if err != nil {
//return nil, err
//}
