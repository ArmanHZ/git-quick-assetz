package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type URLHistory map[string]string

func LoadURLHistory(path string) (URLHistory, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		panic(err)
	}

	hist := make(URLHistory)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return hist, nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return hist, nil
	}

	err = json.Unmarshal(data, &hist)
	return hist, err
}

func WriteURLHistory(urlHistory URLHistory, path string) error {
	data, err := json.MarshalIndent(urlHistory, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LookupURLHash(url string) {

}

func AddNewURL(url string, histFilePath string) {
	urlHistory, err := LoadURLHistory(histFilePath)
	if err != nil {
		panic(err)
	}

	hashBytes := sha256.Sum256([]byte(url))
	hash := hex.EncodeToString(hashBytes[:])

	if _, exists := urlHistory[hash]; !exists {
		urlHistory[hash] = url

		if err := WriteURLHistory(urlHistory, histFilePath); err != nil {
			panic(err)
		}
	}

}
