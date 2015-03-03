// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine

//
// Create goroot.zip for GAE.
//
// Example:
//	go run main.go
//
package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	file, err := os.Create("goroot.zip")
	if err != nil {
		log.Fatal("os.Create: ", err)
	}
	defer file.Close()

	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	// create /goroot/
	f, err := zipFile.Create("goroot/")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.Write([]byte("")); err != nil {
		log.Fatal(err)
	}

	filepath.Walk(runtime.GOROOT(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal("filepath.Walk: ", err)
		}
		if info.IsDir() {
			return nil
		}
		relpath, err := filepath.Rel(runtime.GOROOT(), path)
		if err != nil {
			log.Fatal("filepath.Rel: ", err)
		}

		filename := filepath.ToSlash(relpath)
		if isIngoreFile(filename) {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("ioutil.ReadFile: ", err)
		}

		f, err := zipFile.Create("goroot/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = f.Write(data); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", filename)
		return nil
	})
	fmt.Printf("Done\n")
}

func isIngoreFile(path string) bool {
	if strings.HasPrefix(path, "bin") {
		return true
	}
	if strings.HasPrefix(path, "pkg") {
		return true
	}
	if strings.HasPrefix(path, "translations/.git") {
		return true
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".exe", ".dll":
		return true
	}
	return false
}
