package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func SearchBook(bookTitle string) (*[]models.Book, error) {
	bookTitle = normalizeString(bookTitle)
	bookTitle = strings.Replace(bookTitle, " ", "+", -1)
	foundBooks := []models.Book{}
	resp, err := http.Get("https://openlibrary.org/search.json?q=" + bookTitle)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var response response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	fmt.Println(response.Docs)

	for i := 0; i < len(response.Docs); i++ {
		currentDoc := response.Docs[i]
		var tempBook = models.Book{
			Title:       currentDoc.Title,
			Author:      currentDoc.AuthorName[0],
			Key:         currentDoc.CoverEditinoKey,
			AuthorKey:   currentDoc.AuthorKey[0],
			ReleaseYear: currentDoc.FirstPulishYear,
			AVGRating:   currentDoc.AvgRating,
			PageCount:   currentDoc.NumberOfPages,
		}
		foundBooks = append(foundBooks, tempBook)
	}

	return &foundBooks, nil
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
