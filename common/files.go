package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	fp "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nightlyone/lockfile"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func SafeFileWrite(path string, data []byte) (err error) {
	path, err = overwriteProofPath(path)
	if err != nil {
		return
	}
	err = TouchFile(path)
	if err != nil {
		// TODO LOG
		return
	}
	absPath, err := fp.Abs(path)
	if err != nil {
		// TODO LOG
		return
	}
	path = absPath
	lf, err := lockfile.New(path + ".lock")
	if err != nil {
		// LOG ERR
		return
	}
	err = TryLock(lf)
	if err != nil {
		// LOG ERR
		return
	}
	defer func() {
		err = lf.Unlock()
		// TODO LOG ERR
	}()
	err = ioutil.WriteFile(path, data, 0644)
	// this frees the lock
	if err != nil {
		// TODO Log
		return
	}
	return
}

func overwriteProofPath(path string) (updatedPath string, err error) {
	if !FileExists(path) {
		updatedPath = path
		return
	}
	hasSlash := false // like for /home/user . linux specific!
	if []rune(path)[0] == []rune("/")[0] {
		hasSlash = true
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
		fnSplit = fnSplit[:len(fnSplit)-1]
	}
	sIdx := strconv.Itoa(i_idx)
	fnSplit = append(fnSplit, sIdx)
	newFilename := strings.Join(fnSplit, ".")
	newPath := ""
	if len(pathSplit) < 2 {
		newPath = newFilename
	} else {
		pathSplit = pathSplit[:len(pathSplit)-1]
		pathSplit = append(pathSplit, newFilename)
		newPath = FSJoinSlice(pathSplit)
		if hasSlash {
			newPath = "/" + newPath
		}
	}
	return overwriteProofPath(newPath)
}

func SafeFileRead(path string) (data []byte, err error) {
	absPath, err := fp.Abs(path)
	if err != nil {
		// LOG ERR
		return
	}
	lf, err := lockfile.New(absPath + ".lock")
	if err != nil {
		// LOG ERR
		return
	}
	if !FileExists(absPath) {
		err = fmt.Errorf("SafeFileRead: File does not exist.")
		return
	}
	err = TryLock(lf)
	if err != nil {
		// LOG ERR
		return
	}
	defer func() {
		err = lf.Unlock()
		// TODO LOG ERR
	}()
	data, err = ioutil.ReadFile(absPath)
	// this frees the lock
	if err != nil {
		// LOG
		return
	}
	return
}

func TryLock(lf lockfile.Lockfile) (err error) {
	maxTries := 1000
	sleepInterval := 10 * time.Millisecond
	// 10 seconds worth of tries
	numTries := 0
	for numTries < maxTries {
		err = lf.TryLock()
		/* From the docs:
		   TryLock tries to own the lock. It Returns nil, if successful
		   and error describing the reason, it didn't work out. Please
		   note, that existing lockfiles containing pids of dead
		   processes and lockfiles containing no pid at all are simply
		   deleted.
		*/
		if err == nil {
			return
		}
		time.Sleep(sleepInterval)
		numTries++
	}
	return
}

func FSJoin(folders ...string) (res string) {
	for i := 0; i < len(folders); i++ {
		res = fsJoin(res, folders[i])
	}
	return
}

func FSJoinSlice(folders []string) (res string) {
	for i := 0; i < len(folders); i++ {
		res = fsJoin(res, folders[i])
	}
	return
}

func fsJoin(dir, subpath string) string {
	if IsFullPath(subpath) {
		return subpath
	}
	return fp.Join(dir, subpath)
}

func IsFullPath(path string) bool {
	l := len(path)
	if l == 0 {
		return false
	}
	if path[0] == '/' || path[0] == '\\' {
		return true
	}
	if l >= 2 && path[1] == ':' {
		// Windows
		return true
	}
	return false
}

func CreateFile(fullPath string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(fullPath), os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(fullPath)
	return file, err
}

func TouchFile(absPath string) (err error) {
	if !FileExists(absPath) {
		_, err = CreateFile(absPath)
		if err != nil {
			// TODO LOG
			return
		}
	}
	return
}
