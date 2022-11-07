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
	"log"
	"net/http"
	"os"
	"strconv"
)

// GetImagesResponse contains Images with metadata
type GetImagesResponse struct {
	Images []Image `json:"images"`
}

// Image contains base64 encoded image with metadata
type Image struct {
	Name       string      `json:"name"`
	Position   Coordinates `json:"position"`
	Dimensions Coordinates `json:"dimensions"`
	Base64     string      `json:"base64"`
}

// Coordinates describes image position or dimensions
type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const (
	imagesDir       = "./images"
	base64jpeg      = "data:image/jpeg;base64,"
	base64png       = "data:image/png;base64,"
	correctUsername = "sample_user"
	correctPassword = "sample_password"
)

// GetImages handles requests at `/images`
func GetImages(w http.ResponseWriter, r *http.Request) {
	xRaw := r.URL.Query().Get("x")
	yRaw := r.URL.Query().Get("y")
	zRaw := r.URL.Query().Get("z")

	x, err := strconv.Atoi(xRaw)
	if err != nil {
		log.Fatalf("could not convert %s into int: %s", xRaw, err)
	}
	y, err := strconv.Atoi(yRaw)
	if err != nil {
		log.Fatalf("could not convert %s into int: %s", yRaw, err)
	}
	z, err := strconv.Atoi(zRaw)
	if err != nil {
		log.Fatalf("could not convert %s into int: %s", zRaw, err)
	}

	getImagesResponse, err := calculateImageMatrix(x, y, z)
	if err != nil {
		log.Fatal("calculateImageMatrix")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(getImagesResponse); err != nil {
		log.Fatal(err)
	}
}

// toBase64 encodes bytes as base64 string
func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// getImageDimensions retrieves image dimensions
func getImageDimensions(filePath string) (*image.Config, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %s, error: %s", filePath, err)
	}
	defer imgFile.Close()

	imgConfig, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return nil, fmt.Errorf("could not decode an image: %s, error: %s", filePath, err)
	}

	return &imgConfig, nil
}

// basicAuth performs basic authorization
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providedUsername, providedPassword, ok := r.BasicAuth()
		if ok {
			providedUsernameHash := sha256.Sum256([]byte(providedUsername))
			providedPasswordHash := sha256.Sum256([]byte(providedPassword))
			correctUsernameHash := sha256.Sum256([]byte(correctUsername))
			correctPasswordHash := sha256.Sum256([]byte(correctPassword))

			usernameMatch := subtle.ConstantTimeCompare(providedUsernameHash[:], correctUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(providedPasswordHash[:], correctPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}

// code if needed in the future
// resizeImage resizes given image to specified dimension(s)
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
