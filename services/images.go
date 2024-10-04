package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func ResizeImage(data []byte) ([]byte, error) {
	imageData, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	resizedImage := imaging.Resize(imageData, 180, 0, imaging.Lanczos)
	buff := bytes.NewBuffer(nil)
	err = jpeg.Encode(buff, resizedImage, nil)
	return buff.Bytes(), nil
}

// Downloads the image from the service, resize it, and store it in this server, then returns updates the data with said URL
func ProcessImage(url, bookKey string, done chan (bool), updateFunction func(string)) {
	result := <-done

	if !result {
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		PrintRedError(err.Error())
		return
	}
	defer resp.Body.Close()

	//Gets the image bytes from the URL
	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		PrintRedError(err.Error())
		return
	}

	fileName := fmt.Sprint(bookKey, ".jpg")
	imgUrl := filepath.Join(os.Getenv("IMG_DIR"), fileName)

	//Creates the directory in case it doesnt exists
	err = os.MkdirAll(filepath.Dir(imgUrl), 0777)
	if err != nil {
		PrintRedError(err.Error())
		return
	}

	// resizedImg, err := ResizeImage(imgBytes)
	// if err != nil {
	// 	PrintRedError(err.Error())
	// 	return
	// }
	//writes the image to the hard drive
	err = os.WriteFile(imgUrl, imgBytes, 0644)
	if err != nil {
		PrintRedError(err.Error())
		return
	}
	finalURL := fmt.Sprint(os.Getenv("IMG_URL"), fileName)
	updateFunction(finalURL)
	return
}
