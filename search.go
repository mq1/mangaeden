package mangaeden

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// XMLSearchResult is the search result parsed from the html search results page
type XMLSearchResult struct {
	Title string `xml:",chardata"`
	Class string `xml:"class,attr"`
	Link  string `xml:"href,attr"`
}

type xmlSearchResults struct {
	XMLName       xml.Name          `xml:"tbody"`
	SearchResults []XMLSearchResult `xml:"tr>td>a"`
}

func (x *xmlSearchResults) getSearchResults() []XMLSearchResult {
	var searchResults []XMLSearchResult
	for _, result := range x.SearchResults {
		// the link could be of a single chapter (it doesn't have a class)
		if result.Class != "" {
			searchResults = append(searchResults, result)
		}
	}

	return searchResults
}

// SearchManga searches a manga using the mangaeden search engine
func SearchManga(title, language string) ([]XMLSearchResult, error) {
	// use the mangaeden search engine
	resp, err := http.Get(fmt.Sprintf("https://www.mangaeden.com/%s/%s-directory/?title=%s", language, language, title))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// extract the results table from the html body
	bodyString := string(body)
	tbodyStartPosition := strings.Index(bodyString, "<tbody>")
	tbodyEndPosition := strings.Index(bodyString, "</tbody>")
	toParse := bodyString[tbodyStartPosition:tbodyEndPosition] + "</tbody>"

	// parse the data
	var rawSearchResults xmlSearchResults
	if err := xml.Unmarshal([]byte(toParse), &rawSearchResults); err != nil {
		return nil, err
	}

	// get only the relevant links
	searchResults := rawSearchResults.getSearchResults()

	return searchResults, nil
}
