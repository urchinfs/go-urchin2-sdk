package utils

import (
	"fmt"
	"github.com/ipfs/boxo/files"
	"os"
	"path/filepath"
)

func appendFile(fpath string, recursive bool, filter *files.Filter) (files.Node, error) {
	stat, err := os.Lstat(fpath)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		if !recursive {
			return nil, fmt.Errorf("not support")
		}
	} else if (stat.Mode() & os.ModeNamedPipe) != 0 {
		file, err := os.Open(fpath)
		if err != nil {
			return nil, err
		}

		return files.NewReaderFile(file), nil
	}
	return files.NewSerialFileWithFilter(fpath, filter, stat)
}

func WarpPath(path string) (files.Directory, error) {
	filter, err := files.NewFilter("", []string{}, false)
	if err != nil {
		return nil, err
	}

	fileNode, err := appendFile(path, true, filter)
	if err != nil {
		return nil, err
	}

	fileArgs := make([]files.DirEntry, 0)
	fileArgs = append(fileArgs, files.FileEntry(filepath.Base(filepath.Clean(path)), fileNode))
	wrapDataDir := files.NewSliceDirectory(fileArgs)

	return wrapDataDir, nil
}

func WarpPathWithFilter(path string, filter *files.Filter) (files.Directory, error) {
	fileNode, err := appendFile(path, true, filter)
	if err != nil {
		return nil, err
	}

	fileArgs := make([]files.DirEntry, 0)
	fileArgs = append(fileArgs, files.FileEntry(filepath.Base(filepath.Clean(path)), fileNode))
	wrapDataDir := files.NewSliceDirectory(fileArgs)

	return wrapDataDir, nil
}
