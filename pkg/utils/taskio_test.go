package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"taskTracker/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestSetStorage(t *testing.T){
	tempTask := t.TempDir()
	tempIndex := t.TempDir()
	tempLogs := t.TempDir()

	t.Run("success",func(t *testing.T) {
		err := SetStorage(tempTask, tempIndex, tempLogs)
		require.NoError(t, err)
	})
}

func TestGetTargetPath(t *testing.T) {
	t.Run("correct file Path", func(t *testing.T) {
		tempStorage := t.TempDir()
		fPath := GetTargetPath(tempStorage)
		log.Println(fPath)
		m := fmt.Sprintf("%d.json", time.Now().Local().Month())
		assert.Equal(t, fPath, filepath.Join(tempStorage, strconv.Itoa(time.Now().Local().Year()), m))
	})

	t.Run("no storage dir was created in advance", func(t *testing.T) {
		tempStorage := "tempStorage/tasks"
		fPath := GetTargetPath(tempStorage)
		log.Println(fPath)
		m := fmt.Sprintf("%d.json", time.Now().Local().Month())
		assert.Equal(t, fPath, filepath.Join(tempStorage, strconv.Itoa(time.Now().Local().Year()), m))
		assert.NoError(t, os.RemoveAll(strings.Split(fPath, "/")[0]))
	})
}

func TestPrepareTaskStorage(t *testing.T) {
	t.Run("correct task storage creation", func(t *testing.T) {
		input := "tempStorage/tasks"
		err := prepareTaskStorage(input)
		assert.NoError(t, err)
		p := filepath.Join(input, strconv.Itoa(time.Now().Local().Year()))
		data, err := os.Stat(p)
		require.NoError(t, err)
		require.True(t, data.IsDir())
		assert.NoError(t, os.RemoveAll(strings.Split(input, "/")[0]))
	})

	t.Run("no file owerwriting", func(t *testing.T) {
		input := "tempStorage/tasks"
		err := prepareTaskStorage(input) // create storage for the first time
		assert.NoError(t, err)
		p := filepath.Join(input, strconv.Itoa(time.Now().Local().Year()))
		fName := "stilExists.json"
		fPath := filepath.Join(p, fName)
		f, err := os.Create(fPath) // crete file in the storage
		require.NoError(t, err)
		defer f.Close()

		err = prepareTaskStorage(input) // create storage for the second time
		assert.NoError(t, err)
		data, err := os.Stat(fPath)
		assert.NoError(t, err) // file still exist
		log.Println(data.Name())
		assert.Equal(t, fName, data.Name()) // check its name
		assert.False(t, data.IsDir())
		assert.NoError(t, os.RemoveAll(strings.Split(input, "/")[0]))
	})
}

func TestPrepareLogStorage(t *testing.T) {
	t.Run("correct log storage creation", func(t *testing.T) {
		input := "tempStorage/Logs"
		err := prepareTaskStorage(input)
		assert.NoError(t, err)
		p := filepath.Join(input, strconv.Itoa(time.Now().Local().Year()))
		data, err := os.Stat(p)
		require.NoError(t, err)
		require.True(t, data.IsDir())
		assert.NoError(t, os.RemoveAll(strings.Split(input, "/")[0]))
	})

	t.Run("no file owerwriting", func(t *testing.T) {
		input := "tempStorage/logs"
		err := prepareTaskStorage(input) // create storage for the first time
		assert.NoError(t, err)
		p := filepath.Join(input, strconv.Itoa(time.Now().Local().Year()))
		fName := "stilExists.log"
		fPath := filepath.Join(p, fName)
		f, err := os.Create(fPath) // crete file in the storage
		require.NoError(t, err)
		defer f.Close()

		err = prepareTaskStorage(input) // create storage for the second time
		assert.NoError(t, err)
		data, err := os.Stat(fPath)
		assert.NoError(t, err) // file still exist
		log.Println(data.Name())
		log.Println(fPath)
		assert.Equal(t, fName, data.Name()) // check its name
		assert.False(t, data.IsDir())
		assert.NoError(t, os.RemoveAll(strings.Split(input, "/")[0]))
	})
}

func TestDecodeTasks(t *testing.T) {
	tempName := "tempIndex.json"
	tempDir := os.TempDir()
	fpath := path.Join(tempDir, tempName)
	tempFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0755)
	require.NoError(t, err)
	defer tempFile.Close()
	tempM := make(map[int64]*types.Task)
	tempM[10] = &types.Task{ID: 10, Description: "Success", Done: true}
	err = json.NewEncoder(tempFile).Encode(tempM)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		m := make(map[int64]*types.Task)
		err := DecodeTasks(fpath, m)
		assert.NoError(t, err)
		require.Equal(t, m[10].Description, tempM[10].Description)
		require.Equal(t, m[10], tempM[10])

	})

	t.Run("file does not exists -> error", func(t *testing.T) {
		m := make(map[int64]*types.Task)
		err := DecodeTasks("nonExistingPath/index/i.json", m)
		assert.Error(t, err)
		require.EqualError(t, err, os.ErrNotExist.Error())
	})
}

func TestEncodeTasks(t *testing.T) {
	tempName := "tempIndex.json"
	tempDir := os.TempDir()
	fpath := path.Join(tempDir, tempName)
	tempFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0755)
	require.NoError(t, err)
	defer tempFile.Close()

	t.Run("success", func(t *testing.T) {
		tempM := make(map[int64]*types.Task)
		tempM[10] = &types.Task{ID: 10, Description: "Success", Done: true}
		err = EncodeTasks(fpath, tempM)
		require.NoError(t, err)

		newM := make(map[int64]*types.Task)
		f, err := os.Open(fpath)
		require.NoError(t, err)
		defer f.Close()
		err = json.NewDecoder(f).Decode(&newM)
		require.NoError(t, err)

		require.Equal(t, tempM[10].Description, newM[10].Description)
		require.Equal(t, tempM[1], newM[1])
	})

	t.Run("file does not exists -> error", func(t *testing.T) {
		tempM := make(map[int64]*types.Task)
		tempM[10] = &types.Task{ID: 10, Description: "foooo fail!!!!"}
		err := EncodeTasks("somerandomFile/2025.json", tempM)
		require.Error(t, err)
	})

}

func TestReadLastID(t *testing.T) {
	tempDir := t.TempDir()
	fpath := path.Join(tempDir, "lastID.json")
	t.Run("first Launch -> id == 0", func(t *testing.T) {
		res, err := ReadLastID(fpath)
		require.NoError(t, err)
		require.Equal(t, res, int64(0))
	})

	t.Run("File has already been used", func(t *testing.T) {
		err := WriteLastID(10, fpath)
		require.NoError(t, err)

		res, err := ReadLastID(fpath)
		require.NoError(t, err)
		require.Equal(t, res, int64(10))
	})

	t.Run("incorrect file", func(t *testing.T) {
		res, err := ReadLastID("nonExisitingfilr/2034.json")
		require.Error(t, err)
		require.Equal(t, res, int64(-1))
	})
}

func TestWriteLastID(t *testing.T) {
	tempDir := t.TempDir()
	fpath := path.Join(tempDir, "lastID.json")

	t.Run("write 10 to the file", func(t *testing.T) {
		err := WriteLastID(10, fpath)
		require.NoError(t, err)

		res, err := ReadLastID(fpath)
		require.NoError(t, err)
		require.Equal(t, res, int64(10))
	})

	t.Run("rewrite data", func(t *testing.T) {
		err := WriteLastID(10, fpath)
		require.NoError(t, err)

		err = WriteLastID(44, fpath)
		require.NoError(t, err)

		res, err := ReadLastID(fpath)
		require.NoError(t, err)
		require.Equal(t, res, int64(44))
	})

	t.Run("incorrect file", func(t *testing.T) {
		err := WriteLastID(100, "nonExisitingfilr/2034.json")
		require.Error(t, err)
	})
}