package node

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/nightlyone/lockfile"
)

// TODO put in common

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// TODO TEST LOCK SETUP
func SafeFileWrite(path string, data []byte) (err error) {
	path, err = overwriteProofPath(path)
	if err != nil {
		return
	}
	lf, err := lockfile.New(path)
	if err != nil {
		// LOG ERR
		return
	}
	err = lf.TryLock()
	if err != nil {
		// LOG ERR
		return
	}
	defer func() {
		err = lf.Unlock()
	}()
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		// TODO Log
		return
	}
	return
}

// TODO TEST
func overwriteProofPath(path string) (updatedPath string, err error) {
	if !FileExists(path) {
		updatedPath = path
		return
	}
	pathSplit := strings.Split(path, "/") // does not support windows yet
	filename := pathSplit[len(pathSplit)-1]
	fnSplit := strings.Split(filename, ".")
	// cat.jpg.2
	lastToken := fnSplit[len(fnSplit)-1]
	i_idx := 1
	idx, err := strconv.Atoi(lastToken)
	if err == nil { // lastToken was int
		i_idx = idx + 1
	}
	sIdx := strconv.Itoa(i_idx)
	fnSplit = append(fnSplit, sIdx)
	newFilename := strings.Join(fnSplit, ".")
	return overwriteProofPath(newFilename)
}

// TODO TEST LOCK
func SafeFileRead(path string) (data []byte, err error) {
	lf, err := lockfile.New(path)
	if err != nil {
		// LOG ERR
		return
	}
	err = lf.TryLock()
	if err != nil {
		// LOG ERR
		return
	}
	defer func() {
		err = lf.Unlock()
	}()
	data, err = ioutil.ReadFile(path)
	if err != nil {
		// LOG
		return
	}
	return
}
