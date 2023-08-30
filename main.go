package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

const urlToSearch string = "https://www.google.com/search?tbm=isch&q="

var (
	regexUrl      *regexp.Regexp
	regexFileName *regexp.Regexp
)

func init() {
	re, err := regexp.Compile(`(https:\/\/[^"]+)`)
	if err != nil {
		log.Panic(err)
	}
	regexUrl = re

	re, err = regexp.Compile(`q=tbn:([^&]+)`)
	if err != nil {
		log.Panic(err)
	}

	regexFileName = re
}

func main() {
	args := os.Args

	if len(args) == 1 {
		log.Panic("should pass as args the word to search the image")
	}

	chListToDownload := make(chan string)

	searchWord := args[1]

	CreateDirectory(searchWord)

	go Search(searchWord, chListToDownload)

	wg := sync.WaitGroup{}

	for v := range chListToDownload {
		wg.Add(1)

		go func(v string) {
			DownloadImage(v, searchWord)
			wg.Done()
		}(v)
	}

	wg.Wait()
}

func Search(text string, ch chan<- string) {
	values, err := url.ParseQuery(text)
	if err != nil {
		log.Panic(err)
	}

	resp, err := http.Get(fmt.Sprintf("%s%s", urlToSearch, values.Encode()))
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	urls := regexUrl.FindAllString(string(b), -1)

	for _, v := range urls {
		ch <- v
	}

	close(ch)
}

func DownloadImage(url, directory string) error {
	if !IsValidGoogleImage(url) {
		return errors.New("file is invalid")
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("error to get image from url: %s, error: %v\n", url, err)
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error to read body, error: %v\n", err)
		return err
	}

	fileName := fmt.Sprintf("%s.jpg", regexFileName.FindString(url))

	f, err := os.Create(fmt.Sprintf("%s/%s", directory, fileName))
	if err != nil {
		log.Printf("error to create file, error: %v\n", err)
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		log.Printf("error to write to file, error: %v\n", err)
		return err
	}

	return nil
}

func IsValidGoogleImage(url string) bool {
	return strings.Contains(url, "images?q=tbn")
}

func CreateDirectory(name string) {
	err := os.Mkdir(name, 0777)
	if err != nil {
		log.Panicf("error to create directory: %s - %v", name, err)
	}
}
