package utils

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// PrepareTaskStorage creates a full path to the directory where data will be saved.
func prepareIndexStorage(indexPath string) error {
	if err := os.Mkdir(indexPath, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	return nil
}

func DecodeIndex(fPath string, dst map[int][]int64) error {
	rFile, err := os.OpenFile(fPath, os.O_RDONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			//return errors.New("file not found")
			return os.ErrNotExist // CHANGED: WAS NIL!!!
		}
		return err
	}
	defer rFile.Close()
	if err := json.NewDecoder(rFile).Decode(&dst); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}

func EncodeIndex(fPath string, src map[int][]int64) error {
	wFile, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer wFile.Close()
	if err := json.NewEncoder(wFile).Encode(&src); err != nil {
		return err
	}
	return nil
}
