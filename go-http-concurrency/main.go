package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Photos []struct {
	AlbumId      int    `json:"albumId"`
	Id           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailUrl string `json:"thumbnailUrl"`
}

func downloadPhoto(url string, path string, semaphore chan int) {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	writeErr := os.WriteFile(path, body, 0644)
	if writeErr != nil {
		log.Fatal(writeErr)
	}

	<-semaphore
}

func main() {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, "https://jsonplaceholder.typicode.com/photos", nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
		os.Exit(1)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
		os.Exit(1)
	}

	var photos = Photos{}

	jsonErr := json.Unmarshal(body, &photos)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		os.Exit(1)
	}

	mkdirErr := os.MkdirAll("photos", os.ModePerm)
	if mkdirErr != nil {
		log.Fatal(mkdirErr)
		os.Exit(1)
	}
	semaphore := make(chan int, 20)
	for _, photo := range photos {
		semaphore <- 1
		log.Println("Processing " + photo.Title)
		go downloadPhoto(photo.URL, "photos/"+photo.Title+".png", semaphore)

	}

}
