package util

import (
	"encoding/json"
	"net/http"
)

func GetRandomDogPhotoURL() (string, error) {
	resp, err := http.Get("https://random.dog/woof.json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result["url"], nil
}
