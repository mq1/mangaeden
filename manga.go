package mangaeden

import (
	"encoding/json"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
)

type rawChapterInfo []interface{}

func (c rawChapterInfo) getNumber() float64 {
	return c[0].(float64)
}

func (c rawChapterInfo) getTitle() string {
	return c[2].(string)
}

func (c rawChapterInfo) getID() string {
	return c[3].(string)
}

// rawMangaInfo contains the informations regarding a manga
type rawMangaInfo struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Image       string           `json:"image"`
	Chapters    []rawChapterInfo `json:"chapters"`
}

type chapterInfo struct {
	Number float64
	Title  string
	ID     string
}

func newChapterInfo(raw rawChapterInfo) chapterInfo {
	return chapterInfo{
		Number: raw.getNumber(),
		Title:  raw.getTitle(),
		ID:     raw.getID(),
	}
}

type MangaInfo struct {
	ID          string
	Title       string
	Description string
	Image       string
	Chapters    []chapterInfo
}

func newMangaInfo(raw rawMangaInfo, mangaID string) MangaInfo {
	var chapters []chapterInfo
	for _, c := range raw.Chapters {
		chapters = append(chapters, newChapterInfo(c))
	}

	return MangaInfo{
		ID:          mangaID,
		Title:       raw.Title,
		Description: html.UnescapeString(raw.Description), // unescape the special characters (like &egrave)
		Image:       raw.Image,
		Chapters:    chapters,
	}
}

func GetMangaID(mangaURL string) (string, error) {
	resp, err := http.Get("https://www.mangaeden.com" + mangaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyString := string(body)
	mangaIDStartPosition := strings.Index(bodyString, "window.manga_id2")
	trimmedBodyString := strings.Replace(bodyString, `window.manga_id2 = "`, "", 1)[mangaIDStartPosition:]

	var mangaID strings.Builder
	for i := 0; trimmedBodyString[i] != '"'; i++ {
		mangaID.WriteByte(trimmedBodyString[i])
	}

	return mangaID.String(), nil
}

// GetMangaInfo downloads the manga info corresponding to the id and returns a MangaInfo struct
func GetMangaInfo(mangaID string) (MangaInfo, error) {

	// manga info download
	resp, err := http.Get("https://www.mangaeden.com/api/manga/" + mangaID)
	if err != nil {
		return MangaInfo{}, err
	}
	defer resp.Body.Close()

	// download content reading
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MangaInfo{}, err
	}

	// json content processing
	var raw rawMangaInfo
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return MangaInfo{}, err
	}

	return newMangaInfo(raw, mangaID), nil
}
