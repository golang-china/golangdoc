// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blog

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/tools/godoc/vfs"
)

// SkipDir is used as a return value from WalkFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = errors.New("skip this directory")

// WalkFunc is the type of the function called for each file or directory
// visited by Walk. The path argument contains the argument to Walk as a
// prefix; that is, if Walk is called with "dir", which is a directory
// containing the file "a", the walk function will be called with argument
// "dir/a". The info argument is the os.FileInfo for the named path.
//
// If there was a problem walking to the file or directory named by path, the
// incoming error will describe the problem and the function can decide how
// to handle that error (and Walk will not descend into that directory). If
// an error is returned, processing stops. The sole exception is that if path
// is a directory and the function returns the special value SkipDir, the
// contents of the directory are skipped and processing continues as usual on
// the next file.
type WalkFunc func(fs vfs.FileSystem, path string, info os.FileInfo, err error) error

// walk recursively descends path, calling w.
func walk(fs vfs.FileSystem, path string, info os.FileInfo, walkFn WalkFunc) error {
	err := walkFn(fs, path, info, nil)
	if err != nil {
		if info.IsDir() && err == SkipDir {
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return nil
	}

	names, err := readDirNames(fs, path)
	if err != nil {
		return walkFn(fs, path, info, err)
	}

	for _, name := range names {
		filename := filepath.Join(path, name)
		fileInfo, err := fs.Lstat(filename)
		if err != nil {
			if err := walkFn(fs, filename, fileInfo, err); err != nil && err != SkipDir {
				return err
			}
		} else {
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != SkipDir {
					return err
				}
			}
		}
	}
	return nil
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root. All errors that arise visiting files
// and directories are filtered by walkFn. The files are walked in lexical
// order, which makes the output deterministic but means that for very
// large directories Walk can be inefficient.
// Walk does not follow symbolic links.
func Walk(fs vfs.FileSystem, root string, walkFn WalkFunc) error {
	info, err := fs.Lstat(root)
	if err != nil {
		return walkFn(fs, root, nil, err)
	}
	return walk(fs, root, info, walkFn)
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(fs vfs.FileSystem, dirname string) ([]string, error) {
	fileInfoList, err := fs.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, fi := range fileInfoList {
		names = append(names, fi.Name())
	}
	sort.Strings(names)
	return names, nil
}
