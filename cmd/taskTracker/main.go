package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"taskTracker/pkg/task"
	"taskTracker/pkg/types"
	"taskTracker/pkg/utils"

	"log"
)

const (
	STORAGE_LAST_ID = "storage/lastID.json"
	TASK_STORAGE    = "storage/tasks"
	INDEX_STORAGE   = "storage/index"
	LOG_STORAGE     = "storage/logs"
)

func main() {
	if err := utils.SetStorage(TASK_STORAGE, INDEX_STORAGE, LOG_STORAGE); err != nil {
		log.Fatal(err)
	}

	if err := execute(); err != nil{
		log.Fatal(err)
	}
}

func execute() error {
	fPath := utils.GetTargetPath(TASK_STORAGE)
	lastID, err := utils.ReadLastID(STORAGE_LAST_ID)
	if err != nil {
		return err
	}

	createFlag := flag.Bool("c", false, "create task")
	updateFlag := flag.Bool("u", false, "update task. used with -id flag")
	deleteFlag := flag.Bool("d", false, "delete task. used with -id flag")
	getByIDFlag := flag.Bool("g", false, "get one task by id. use with -id flag")
	getTodayFlag := flag.Bool("today", false, "get todays tasks")
	descFlag := flag.String("desc", "", "description fot your task")
	doneFlag := flag.Bool("done", false, "task status (done or not)")
	idFlag := flag.Int64("id", 0, "indicate in case you want to update a task")

	listFlag := flag.Bool("ld", false, "get tasks according to the -day, -m, -y")
	dayFlag := flag.Int("day", -1, "indicate day")
	monthFlag := flag.Int("m", -1, "indicate month")
	yearFlag := flag.Int("y", -1, "indicate year")

	helpFlag := flag.Bool("h", false, "help")

	flag.Parse()
	if *helpFlag {
		utils.Help()
	}

	if *createFlag {
		if *descFlag == "" {
			return errors.New("provide task description")
		}
		if err := task.CreateTask(fPath, *descFlag, *doneFlag, lastID); err != nil {
			return err
		}
	}
	if *updateFlag && *idFlag > 0 {
		targetFile, err := task.SearchByID(*idFlag, INDEX_STORAGE, TASK_STORAGE)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("task does not exist")
				return nil
			}
			return err
		}
		if err := task.Update(*idFlag, *doneFlag, *descFlag, targetFile); err != nil {
			log.Fatal(err)
		}
	}
	if *deleteFlag && *idFlag > 0 {
		targetFile, err := task.SearchByID(*idFlag, INDEX_STORAGE, TASK_STORAGE)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("task does not exist")
				return os.ErrNotExist
			}
			return err
		}
		if err := task.Delete(*idFlag, targetFile); err != nil {
			log.Fatal(err)
		}
	}
	if *getTodayFlag {
		targetFile, err := task.SearchByID(lastID, INDEX_STORAGE, TASK_STORAGE)
		if err != nil {
			return err
		}
		arr, err := task.GetToday(lastID, targetFile)
		if err != nil {
			return err
		}

		for _, t := range arr {
			utils.ShowTask(*t)
		}
		fmt.Println("Total tasks:", len(arr))
	}

	if *getByIDFlag && *idFlag > 0 {
		targetFile, err := task.SearchByID(*idFlag, INDEX_STORAGE, TASK_STORAGE)
		if err != nil {
			return err
		}
		t, err := task.GetByID(*idFlag, targetFile)
		if err != nil {
			return err
		}
		utils.ShowTask(*t)
	}

	if *listFlag {
		f := types.NewFilter()
		if *dayFlag > 0 && *dayFlag < 32 {
			f.Day = *dayFlag
		}
		if *monthFlag > 0 && *monthFlag < 13 {
			f.Month = *monthFlag
		}
		if *yearFlag > 2025 {
			f.Year = *yearFlag
		}
		arr, err := task.GetByDate(TASK_STORAGE, f)
		if err != nil {
			return err
		}
		for _, elem := range arr {
			utils.ShowTask(*elem)
		}
	}

	log.Println("Finished successfully")
	return nil
}
