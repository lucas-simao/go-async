package main

import (
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	searchWord := "pink home"

	ch := make(chan string)

	go Search(searchWord, ch)

	var countResult int

	for v := range ch {
		if len(v) > 0 {
			countResult++
		}
	}

	if countResult < 10 {
		t.Errorf("expect 10 but got %v", countResult)
	}
}

func TestDownloadImage(t *testing.T) {
	tests := map[string]struct {
		imageUrl string
		hasError bool
	}{
		"success": {
			imageUrl: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSVsq7Up_V1vTGhvamYj7EEkYmvD2NGidL-j3sOjKCeQBnIRPFBSluVU3N7aw&amp;s",
			hasError: false,
		},
		"error, not is valid url": {
			imageUrl: "https://www.google.com/imgres?imgurl=https%3A%2F%2Fodia.ig.com.br%2F_midias%2Fjpg%2F2021%2F10%2F29%2F458x280%2F1_whatsapp_image_2021_10_28_at_21_10_09__1_-23431481.jpeg&tbnid=kItz5UrnRfjtcM&vet=12ahUKEwir89GCrYWBAxXLNbkGHdVlDqkQMyg1egUIARDwAQ..i&imgrefurl=https%3A%2F%2Fodia.ig.com.br%2Fbarra-mansa%2F2021%2F10%2F6264966-casa-rosa-atende-pacientes-em-tratamento-contra-o-cancer-em-barra-mansa.html&docid=jnLkX3xmbo95xM&w=458&h=279&q=casa%20rosa&ved=2ahUKEwir89GCrYWBAxXLNbkGHdVlDqkQMyg1egUIARDwAQ",
			hasError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tempDir := os.TempDir()

			err := DownloadImage(test.imageUrl, tempDir)
			if test.hasError {
				if err == nil {
					t.Errorf("expect error but got %v", err)
				}
			}
		})
	}
}

func TestIsValidGoogleImage(t *testing.T) {
	tests := map[string]struct {
		url     string
		isValid bool
	}{
		"expect true": {
			url:     "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSVsq7Up_V1vTGhvamYj7EEkYmvD2NGidL-j3sOjKCeQBnIRPFBSluVU3N7aw&amp;s",
			isValid: true,
		},
		"expect false": {
			url:     "https://www.google.com/imgres?imgurl=https%3A%2F%2Fodia.ig.com.br%2F_midias%2Fjpg%2F2021%2F10%2F29%2F458x280%2F1_whatsapp_image_2021_10_28_at_21_10_09__1_-23431481.jpeg&tbnid=kItz5UrnRfjtcM&vet=12ahUKEwir89GCrYWBAxXLNbkGHdVlDqkQMyg1egUIARDwAQ..i&imgrefurl=https%3A%2F%2Fodia.ig.com.br%2Fbarra-mansa%2F2021%2F10%2F6264966-casa-rosa-atende-pacientes-em-tratamento-contra-o-cancer-em-barra-mansa.html&docid=jnLkX3xmbo95xM&w=458&h=279&q=casa%20rosa&ved=2ahUKEwir89GCrYWBAxXLNbkGHdVlDqkQMyg1egUIARDwAQ",
			isValid: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resp := IsValidGoogleImage(test.url)

			if resp != test.isValid {
				t.Errorf("expect %v but got %v", test.isValid, resp)
			}
		})
	}
}
