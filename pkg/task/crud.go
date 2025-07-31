package task

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"taskTracker/pkg/types"
	"taskTracker/pkg/utils"
	"time"
)

const (
	STORAGE_LAST_ID = "storage/lastID.json"
	TASK_STORAGE    = "storage/tasks"
	INDEX_STORAGE   = "storage/index"
	LOG_STORAGE     = "storage/logs"
)

func CreateTask(fPath, desc string, status bool, lastID int64) error {
	year, month, _ := time.Now().Local().Date()
	tMap := make(map[int64]*types.Task)
	iMap := make(map[int][]int64)

	task := types.Task{ID: lastID + 1, Description: desc, Done: status, CreatedAt: time.Now().Local(), UpdateAt: time.Now().Local()}

	if err := utils.DecodeTasks(fPath, tMap); err != nil {
		log.Fatal(err)
	}
	tMap[task.ID] = &task

	if err := utils.EncodeTasks(fPath, tMap); err != nil {
		return err
	}
	if err := utils.WriteLastID(task.ID, STORAGE_LAST_ID); err != nil {
		return err
	}

	iFile := filepath.Join(INDEX_STORAGE, fmt.Sprintf("%v.json", year))
	if err := utils.DecodeIndex(iFile, iMap); err != nil {
		return err
	}
	arr := iMap[int(month)]
	if len(arr) > 1 {
		arr[1] = task.ID // in order to keep only first id and last id of the tasks created at month
	} else {
		arr = append(arr, task.ID)
	}
	iMap[int(month)] = arr
	if err := utils.EncodeIndex(iFile, iMap); err != nil {
		return err
	}
	return nil
}

// Update updates the task
func Update(id int64, done bool, desc, targetFile string) error {
	tMap := make(map[int64]*types.Task)
	if err := utils.DecodeTasks(targetFile, tMap); err != nil {
		return err
	}
	t := tMap[id]
	if desc != "" {
		t.Description = desc
	}
	t.Done = done
	t.UpdateAt = time.Now().Local()
	tMap[id] = t

	if err := utils.EncodeTasks(targetFile, tMap); err != nil {
		return err
	}
	return nil
}

func Delete(id int64, targetFile string) error {
	tMap := make(map[int64]*types.Task)
	if err := utils.DecodeTasks(targetFile, tMap); err != nil {
		return err
	}
	delete(tMap, id)
	if err := utils.EncodeTasks(targetFile, tMap); err != nil {
		return err
	}
	return nil
}

// SearchByID search path to file where task was created using index (iStorage should be INDEX_STORAGE)
func SearchByID(id int64, iStorage, tStoarge string) (string, error) {
	wg := sync.WaitGroup{}
	resChan := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	files, err := os.ReadDir(iStorage)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		iMap := make(map[int][]int64)
		if file.IsDir() {
			continue
		}
		if err := utils.DecodeIndex(filepath.Join(iStorage, file.Name()), iMap); err != nil {
			continue
		}
		wg.Add(1)
		go func(id int64, m map[int][]int64, fName string) {
			defer wg.Done()
			for k, val := range m {
				select {
				case <-ctx.Done():
					return
				default:
					if len(val) < 2 {
						if id == val[0] {
							resChan <- fmt.Sprintf("%v/%v/%v.json", tStoarge, fName, k)
							return
						}
						continue
					}
					if val[0] <= id && id <= val[1] {
						resChan <- fmt.Sprintf("%v/%v/%v.json", tStoarge, fName, k)
						return
					}
				}
			}
		}(id, iMap, strings.Split(file.Name(), ".")[0]) // split in order to get only year without file type
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()
	res, ok := <-resChan
	if !ok {
		cancel()
		return "", os.ErrNotExist
	}
	return res, nil
}

func GetToday(lastID int64, fPath string) ([]*types.Task, error) {
	dst := make([]*types.Task, 0)
	m := make(map[int64]*types.Task)
	if err := utils.DecodeTasks(fPath, m); err != nil {
		return nil, err
	}
	y, mt, d := time.Now().Local().Date()
	_, ok := m[lastID]
	if !ok {
		return nil, os.ErrInvalid
	}
	elem := m[lastID].CreatedAt
	for ok && elem.Year() == y && elem.Month() == mt && elem.Day() == d {
		dst = append(dst, m[lastID])
		lastID--
		_, ok = m[lastID]
		if !ok {
			break
		}
		elem = m[lastID].CreatedAt
	}
	return dst, nil
}

func GetByID(id int64, fPath string) (*types.Task, error) {
	m := make(map[int64]*types.Task)
	if err := utils.DecodeTasks(fPath, m); err != nil {
		return nil, err
	}
	if _, ok := m[id]; !ok {
		return nil, os.ErrNotExist
	}
	return m[id], nil
}

func GetByDate(tStoragePath string, f *types.Filter) ([]*types.Task, error){
	arr := make([]*types.Task, 0)
	tMap := make(map[int64]*types.Task)
	fPath := filepath.Join(tStoragePath, strconv.Itoa(f.Year), fmt.Sprintf("%d.json", f.Month))
	if err := utils.DecodeTasks(fPath, tMap); err != nil{
		return nil, err
	}
	for _, t := range tMap{
		arr = append(arr, t)
	}

	return arr, nil
}
