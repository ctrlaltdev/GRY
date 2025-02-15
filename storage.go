package main

import (
	"errors"
	"os"
	"path/filepath"
)

func GetURL(slug string) (url string, err error) {
	content, fileErr := os.ReadFile(filepath.Join(STORAGE_PATH, slug))
	if fileErr != nil {
		return "", fileErr
	}

	url = string(content)

	return url, nil
}

func CreateURL(slug string, url string) (err error) {
	content := []byte(url)

	_, err = os.Stat(filepath.Join(STORAGE_PATH, slug))
	if err == nil {
		return errors.New("slug already exists")
	}

	err = os.WriteFile(filepath.Join(STORAGE_PATH, slug), content, 0640)
	if err != nil {
		LogErr(err)
		return err
	}

	return nil
}

func UpdateURL(slug string, url string) (err error) {
	content := []byte(url)

	_, err = os.Stat(filepath.Join(STORAGE_PATH, slug))
	if err != nil {
		return errors.New("slug does not exist")
	}

	err = os.WriteFile(filepath.Join(STORAGE_PATH, slug), content, 0640)
	if err != nil {
		LogErr(err)
		return err
	}

	return nil
}

func DeleteURL(slug string) (err error) {
	_, err = os.Stat(filepath.Join(STORAGE_PATH, slug))
	if err != nil {
		return errors.New("slug does not exist")
	}

	err = os.Remove(filepath.Join(STORAGE_PATH, slug))
	if err != nil {
		LogErr(err)
		return err
	}

	return nil
}
