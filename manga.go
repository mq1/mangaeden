package mangaeden

import (
	"encoding/json"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
)

type rawChapterInfo []interface{}

func (c rawChapterInfo) getNumber() float64 {
	return c[0].(float64)
}

func (c rawChapterInfo) getTitle() string {
	if c[2] == nil {
		return ""
	}
	return c[2].(string)
}

func (c rawChapterInfo) getID() string {
	if c[3] == nil {
		return ""
	}
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

// GetMangaID gets the page from a link like "/en/en-manga/soul-eater/" and extracts the manga's ID
func GetMangaID(mangaURL string) (string, error) {
	pattern := regexp.MustCompile(`window\.manga_id2 = "([^"]*)"`)

	resp, err := http.Get("https://www.mangaeden.com" + mangaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return pattern.FindStringSubmatch(string(body))[1], nil
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
