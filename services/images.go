package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Descarga la imagen del servidor y la almacena en un directorio local, devuelve la url de la imagen y un error si lo hubo
func ProcessImage(url, bookKey string) (string, error) {
	// Procesa la imagen y la almacena en el servidor
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//obtener la imagen de la respuesta del servidor
	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	imgUrl := filepath.Join(homeDir, "images", fmt.Sprint(bookKey, ".jpg"))
	directory := filepath.Dir(imgUrl)
	//crea el directorio en caso de que no exista
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", err
	}

	return imgUrl, os.WriteFile(imgUrl, imgBytes, 0644)

}
