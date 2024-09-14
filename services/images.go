package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Debería de agregar sincronización con un canal, pero de momento está bien
// Descarga la imagen del servidor y la almacena en un directorio local, devuelve la url de la imagen y un error si lo hubo
func ProcessImage(url, bookKey string, done chan (any), updateFunction func(string)) {
	// Procesa la imagen y la almacena en el servidor
	resp, err := http.Get(url)
	if err != nil {
		PrintRedError(err.Error())
		return
	}
	defer resp.Body.Close()

	//obtener la imagen de la respuesta del servidor
	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		PrintRedError(err.Error())
		return
	}

	fileName := fmt.Sprint(bookKey, ".jpg")
	imgUrl := filepath.Join(os.Getenv("IMG_DIR"), fileName)

	//crea el directorio en caso de que no exista
	err = os.MkdirAll(filepath.Dir(imgUrl), 0777)
	if err != nil {
		PrintRedError(err.Error())
		return
	}
	//escribe la imagen en el disco duro
	err = os.WriteFile(imgUrl, imgBytes, 0644)
	if err != nil {
		PrintRedError(err.Error())
		return
	}
	finalURL := fmt.Sprint(os.Getenv("IMG_URL"), fileName)

	<-done
	updateFunction(finalURL)
}
