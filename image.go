package mangaeden

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func DownloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get("https://cdn.mangaeden.com/mangasimg/" + imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type page []interface{}

type chapter struct {
	Images []page `json:"images"`
}

func (p page) getPageNumber() float64 {
	return p[0].(float64)
}

func (p page) getPageImageURL() string {
	return p[1].(string)
}

// extension detection (jpg/png)
func getImageExtension(imageURL string) string {
	if imageURL[len(imageURL)-3] == 'j' || imageURL[len(imageURL)-3] == 'J' {
		return ".jpg"
	}
	return ".png"
}

func saveImage(fileWithoutExtension, imageURL string) error {
	// create the file
	out, err := os.Create(fileWithoutExtension + getImageExtension(imageURL))
	if err != nil {
		return err
	}
	defer out.Close()

	// download the image
	bytes, err := DownloadImage(imageURL)
	if err != nil {
		return err
	}

	// save the image to the file
	if _, err := out.Write(bytes); err != nil {
		return err
	}

	return nil
}

func DownloadChapter(chapterID, directory string) error {
	// download the chapter's image list
	resp, err := http.Get("https://www.mangaeden.com/api/chapter/" + chapterID)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read the download content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parse the json
	var raw chapter
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return err
	}

	for _, page := range raw.Images {
		if err := os.MkdirAll(directory, os.ModePerm); err != nil {
			return err
		}
		if err := saveImage(fmt.Sprintf("%s/%v", directory, page.getPageNumber()), page.getPageImageURL()); err != nil {
			return err
		}
	}

	return nil
}
