package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func DirExists(dir string, create bool) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if create {
			err := os.MkdirAll(dir, 0755)
			return err
		}
		return err
	}
	return err
}

func FileExists(filename string, create bool) error {
	s, err := os.Stat(filename)
	if os.IsNotExist(err) {
		if create {
			_, err := os.Create(filename)
			return err
		}
		return err
	}
	if s.IsDir() {
		return errors.New("the path is a directory, not a normal file")
	}
	return nil
}

func WriteToJson(content interface{}, toFile string) error {
	file, err := os.Create(toFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	jsonData, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func ReadJson(fromFile string, mute bool) ([]byte, error) {
	file, err := os.Open(fromFile)
	if !mute && HandleError(err, "Error open file") {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	jsonData, err := os.ReadFile(fromFile)
	if !mute && HandleError(err, "Error read file") {
		return nil, err
	}

	return jsonData, nil
}

func WalkDirGlob(dir string) ([]string, error) {
	matches, err := filepath.Glob(dir)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
