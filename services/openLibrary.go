package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/TheSgtPepper23/GreenLibrary/models"
	"golang.org/x/text/unicode/norm"
)

type response struct {
	NumFound int   `json:"numFound"`
	Docs     []doc `json:"docs"`
}

type doc struct {
	AuthorKey       []string `json:"author_key"`
	AuthorName      []string `json:"author_name"`
	CoverEditinoKey string   `json:"cover_edition_key"`
	FirstPulishYear int      `json:"first_publish_year"`
	NumberOfPages   int      `json:"number_of_pages_median"`
	Title           string   `json:"title"`
	AvgRating       float32  `json:"ratings_average"`
}

func SearchBook(bookTitle string, target *[]models.Book, wg *sync.WaitGroup, mu *sync.Mutex, errChan chan (error)) {
	defer wg.Done()
	start := time.Now()
	bookTitle = normalizeString(bookTitle)
	bookTitle = strings.ReplaceAll(bookTitle, " ", "+")

	req, err := http.NewRequest("GET", os.Getenv("OPEN_LIBRARY_URL")+bookTitle, nil)
	if err != nil {
		errChan <- err
		return
	}

	req.Header.Add("User-Agent", "bluefive.xyz:greenLibrary:andresdglez@gmail.com")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		errChan <- err
		return
	}

	defer resp.Body.Close()

	var response response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		errChan <- err
		return
	}
	baseImage := os.Getenv("IMAGE_URL")

	for i := 0; i < len(response.Docs); i++ {
		if response.Docs[i].CoverEditinoKey == "" {
			continue
		}
		currentDoc := response.Docs[i]

		authorName := "Unknown"
		if len(currentDoc.AuthorName) > 0 {
			authorName = strings.Join(currentDoc.AuthorName, ", ")
		}

		authorKey := ""
		if len(currentDoc.AuthorKey) > 0 {
			authorKey = strings.Join(currentDoc.AuthorKey, ", ")
		}

		var tempBook = models.Book{
			Title:       currentDoc.Title,
			Author:      authorName,
			Key:         currentDoc.CoverEditinoKey,
			AuthorKey:   authorKey,
			ReleaseYear: currentDoc.FirstPulishYear,
			AVGRating:   currentDoc.AvgRating,
			PageCount:   currentDoc.NumberOfPages,
			CoverURL:    buildImageURL(currentDoc.CoverEditinoKey, baseImage),
		}
		mu.Lock()
		*target = append(*target, tempBook)
		mu.Unlock()
	}

	fmt.Println(time.Since(start).Milliseconds())
	return
}

func buildImageURL(key, baseURL string) string {
	return fmt.Sprint(baseURL, key, "-", "M", ".jpg")
}

func normalizeString(original string) string {
	normInput := norm.NFD.String(original)

	var sb strings.Builder
	for _, r := range normInput {
		if !unicode.Is(unicode.Mn, r) {
			sb.WriteRune(r)
		}
	}
	return norm.NFC.String(sb.String())
}
