// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows
// +build !appengine

package main

import (
	"flag"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	playEnabled = *showPlayground

	// Check usage: either server and no args, command line and args, or index creation mode
	if (*httpAddr != "" || *urlFlag != "") != (flag.NArg() == 0) && !*writeIndex {
		usage()
	}

	runGodoc()
}
