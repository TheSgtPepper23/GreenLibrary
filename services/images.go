package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Although this could cause duplicated images if the file is already downloaded and the main routine fails.
// It will be override later since the file name will be the same. Orphan images is a problem though
// Descarga la imagen del servidor y la almacena en un directorio local, devuelve la url de la imagen y un error si lo hubo
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
	//writes the image to the hard drive
	err = os.WriteFile(imgUrl, imgBytes, 0644)
	if err != nil {
		PrintRedError(err.Error())
		return
	}
	finalURL := fmt.Sprint(os.Getenv("IMG_URL"), fileName)
	updateFunction(finalURL)
}
