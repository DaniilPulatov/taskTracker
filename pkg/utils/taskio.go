package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"taskTracker/pkg/types"
	"time"
)

var (
	label  = "lastID"
)

func SetStorage(taskPath, indexPath, logsPath string) error {
	if err := prepareTaskStorage(taskPath); err != nil {
		return err
	}
	if err := prepareLogStorage(logsPath); err != nil {
		return err
	}
	if err := prepareIndexStorage(indexPath); err != nil {
		return err
	}
	return nil
}

// GetTargetPath returns path to the file to write
func GetTargetPath(storagePath string) (string) {
	t := time.Now().Local()
	fName := fmt.Sprintf("%d.json", t.Month())
	fPath := filepath.Join(storagePath, strconv.Itoa(t.Year()), fName)

	return fPath
}

// prepareTaskStorage creates a full path to the directory where data will be saved.
func prepareTaskStorage(storagePath string) error {
	curDir := fmt.Sprintf("%d", time.Now().Local().Year())
	fPath := filepath.Join(storagePath, curDir)
	if err := os.MkdirAll(fPath, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	return nil
}

// prepareLogStorage creates full path for logs
func prepareLogStorage(storagePath string) error {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	return nil
}

// DecodeFile parse json data from file into array
func DecodeTasks(fPath string, dst map[int64]*types.Task) error {
	rFile, err := os.OpenFile(fPath, os.O_RDONLY, 0755)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			//return errors.New("file not found")
			return os.ErrNotExist
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

// EncodeFile saves data from src to json file.
func EncodeTasks(fPath string, src map[int64]*types.Task) error {
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

func WriteLastID(id int64, fPath string) error {
	m := make(map[string]int64)
	m[label] = id
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(&m); err != nil {
		return err
	}
	return nil
}

func ReadLastID(fPath string) (int64, error) {
	m := make(map[string]int64)
	file, err := os.OpenFile(fPath, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return -1, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&m); err != nil {
		if !errors.Is(err, io.EOF) {
			return -1, err
		}
	}
	_, ok := m[label] // in case file has never been opened before
	if !ok {
		return 0, nil
	}
	return m[label], nil

}

func ShowTask(t types.Task){
	fmt.Printf("%v. %v / finished: %v --- Created: %v --- Updated: %v\n", t.ID, t.Description, t.Done, t.CreatedAt.Format(time.RFC822), t.UpdateAt.Format(time.RFC822))
}

func Help(){
	println("********************************************************************")
	println()
	fmt.Println("-c: flag for creating task. Indicate Description and status of the task via -desc, -done flags")
	fmt.Println("-u: flag for updating task. Require ID (-id flag )of the target task. Indicate description and status with -desc and -done flags")
	fmt.Println("-d: flag for deleting task. Require ID of the task (-id flag)")
	fmt.Println("-g: flag for getting full info about the task with specified ID (-id flag)")
	fmt.Println("-ld: flag for getting tasks created on specified date -day, -m, -y")
	fmt.Println("-today: flag that will provide all the tasks created today")
	fmt.Println("-desc: flag for indicating description of the task")
	fmt.Println("-done: flag for indicating status of task (true if done else false)")
	fmt.Println("-id: flag for indicating id of the target task")
	fmt.Println("-day: flag for indicating day for filter")
	fmt.Println("-month: flag for indicating month for filter")
	fmt.Println("-year: flag for indicating year for filter")
	println()
	println("********************************************************************")
}