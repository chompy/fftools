package main

import (
	"os"
	"path/filepath"
	"sync"
)

type LogFileWriter struct {
	lock sync.Mutex
	name string
	fp   *os.File
	size int64
}

func (w *LogFileWriter) path() string {
	return filepath.Join(getLogPath(), w.name+".log")
}

func (w *LogFileWriter) open() error {
	if w.fp != nil {
		w.fp.Close()
	}
	var err error
	w.fp, err = os.OpenFile(w.path(), os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	stat, err := w.fp.Stat()
	if err != nil {
		return err
	}
	w.size = stat.Size()
	return nil
}

func (w *LogFileWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.fp == nil {
		if err := w.open(); err != nil {
			return 0, err
		}
	}
	//config := configAppLoad()
	w.size += int64(len(output))
	return w.fp.Write(output)
}

func (w *LogFileWriter) trim() error {
	/*w.lock.Lock()
	defer w.lock.Unlock()
	if w.fp == nil {
		if err := w.open(); err != nil {
			return err
		}
	}
	// stat log file
	s, err := w.fp.Stat()
	if err != nil {
		return err
	}
	// calculate how many bytes over size limit log file is
	config := configAppLoad()
	overage := s.Size() - config.LogMaxSize
	// not reached max size
	if overage < 0 {
		return nil
	}
	// overage is too large, just delete file
	if overage > config.LogMaxSize*2 {
		w.fp.Close()
		return os.Remove(w.path())
	}
	// scan line by line for overage content
	w.fp.Seek(0, os.SEEK_SET)
	scanner := bufio.NewScanner(w.fp)
	bytesScanned := int64(0)
	for scanner.Scan() {
		bytesScanned += int64(len(scanner.Bytes()))
		if bytesScanned >= overage {
			break
		}
	}
	// scan line by line for content that should be remain
	trimmedContents := make([]byte, 0)
	for scanner.Scan() {
		trimmedContents = append(trimmedContents, scanner.Bytes()...)
		trimmedContents = append(trimmedContents, '\n')
	}
	f.Close()
	// delete old file
	if err := os.Remove(LogPath); err != nil {
		return errors.WithStack(err)
	}
	// create new
	f, err = os.OpenFile(LogPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	if _, err := f.Write(trimmedContents); err != nil {
		return errors.WithStack(err)
	}
	if err := f.Sync(); err != nil {
		return errors.WithStack(err)
	}
	return nil*/
	return nil
}
