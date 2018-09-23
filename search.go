package mangaeden

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type SearchResult struct {
	Title string
	Class string
	Link  string
}

func newSearchResult(match []string) SearchResult {
	return SearchResult{
		Title: match[3],
		Class: match[2],
		Link:  match[1],
	}
}

// SearchManga searches a manga using the mangaeden search engine
func SearchManga(title, language string) ([]SearchResult, error) {
	pattern := regexp.MustCompile(`<a href="([^"]*)" class="(.*Manga)">([^<]*)<\/a>`)

	// use the mangaeden search engine
	resp, err := http.Get(fmt.Sprintf("https://www.mangaeden.com/en/%s-directory/?title=%s", language, title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResults []SearchResult

	for _, match := range pattern.FindAllStringSubmatch(string(body), -1) {
		searchResults = append(searchResults, newSearchResult(match))
	}

	return searchResults, nil
}
