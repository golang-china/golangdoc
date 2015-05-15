// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Godoc for Golang, support translate.

Download Chinese Translate:

	git clone https://github.com/golang-china/golangdoc.translations.git $(GOROOT)/translations

Start Chinese Godoc Server:

	go get github.com/golang-china/golangdoc
	golangdoc -http=:6060 -lang=zh_CN

See:

	http://github.com/golang-china/golangdoc.translations
	http://godoc.org/github.com/golang-china/golangdoc/docgen
	http://godoc.org/github.com/golang-china/golangdoc/local


Run as windows service:

	# install as windows service
	golangdoc -service-install -http=:6060

	# start/stop service
	golangdoc -service-start
	golangdoc -service-stop

	# remove service
	golangdoc -service-remove


BUGS

Report bugs to chaishushan@gmail.com.

Thanks!
*/
package main
