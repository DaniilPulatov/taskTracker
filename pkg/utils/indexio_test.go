package utils

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareIndexStorage(t *testing.T) {
	t.Run("correct task storage creation", func(t *testing.T) {
		input := "tempStorage/index"
		err := prepareTaskStorage(input)
		assert.NoError(t, err)
		p := filepath.Join(input, strconv.Itoa(time.Now().Local().Year()))
		data, err := os.Stat(p)
		require.NoError(t, err)
		require.True(t, data.IsDir())
		assert.NoError(t, os.RemoveAll(strings.Split(input, "/")[0]))
	})

	t.Run("no file owerwriting", func(t *testing.T) {
		input := "tempStorage/index"
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

func TestDecodeIndex(t *testing.T) {
	tempName := "tempIndex.json"
	tempDir := os.TempDir()
	fpath := path.Join(tempDir, tempName)
	tempFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0755)
	require.NoError(t, err)
	defer tempFile.Close()
	tempM := make(map[int][]int64)
	tempM[10] = []int64{1, 2, 3, 4, 5}
	tempM[1] = []int64{10, 11}
	err = json.NewEncoder(tempFile).Encode(tempM)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		m := make(map[int][]int64)
		err := DecodeIndex(fpath, m)
		assert.NoError(t, err)
		require.Equal(t, m[10][0], int64(1))
		require.Equal(t, m[10], tempM[10])
		require.Equal(t, len(m[1]), 2)
	})

	t.Run("file does not exists -> error", func(t *testing.T) {
		m := make(map[int][]int64)
		err := DecodeIndex("nonExistingPath/index/i.json", m)
		assert.Error(t, err)
		require.EqualError(t, err, os.ErrNotExist.Error())
	})
}

func TestEncodeIndex(t *testing.T) {
	tempName := "tempIndex.json"
	tempDir := os.TempDir()
	fpath := path.Join(tempDir, tempName)
	tempFile, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0755)
	require.NoError(t, err)
	defer tempFile.Close()

	t.Run("success", func(t *testing.T) {
		tempM := make(map[int][]int64)
		tempM[10] = []int64{1, 2, 3, 4, 5}
		tempM[1] = []int64{10, 11}
		err = EncodeIndex(fpath, tempM)
		require.NoError(t, err)

		newM := make(map[int][]int64)
		f, err := os.Open(fpath)
		require.NoError(t, err)
		defer f.Close()
		err = json.NewDecoder(f).Decode(&newM)
		require.NoError(t, err)

		require.Equal(t, len(tempM[10]), len(newM[10]))
		require.Equal(t, tempM[1], newM[1])
	})

	t.Run("file does not exists -> error", func(t *testing.T) {
		tempM := make(map[int][]int64)
		tempM[10] = []int64{1, 2, 3, 4, 5}
		tempM[1] = []int64{10, 11}
		err = EncodeIndex("somerandomFile/2025.json", tempM)
		require.Error(t, err)
	})
}