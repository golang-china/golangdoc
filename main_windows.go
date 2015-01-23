// Copyright 2015 ChaiShushan <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Godoc service
//
//	# install as windows service
//	golangdoc -service-install -http=:6060
//
//	# start/stop service
//	golangdoc -service-start
//	golangdoc -service-stop
//
//	# remove service
//	golangdoc -service-remove
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"code.google.com/p/winsvc/svc"
)

const (
	ServiceName = "golangdoc"
	ServiceDesc = "Go Documentation Server"
)

var (
	flagServiceInstall   = flag.Bool("service-install", false, "Install service")
	flagServiceUninstall = flag.Bool("service-remove", false, "Remove service")
	flagServiceStart     = flag.Bool("service-start", false, "Start service")
	flagServiceStop      = flag.Bool("service-stop", false, "Stop service")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	playEnabled = *showPlayground

	if *flagServiceInstall {
		var args []string
		args = append(args, fmt.Sprintf("-goroot=%s", *goroot))
		for i := 1; i < len(os.Args); i++ {
			if strings.HasPrefix(os.Args[i], "-service-install") {
				continue
			}
			if strings.HasPrefix(os.Args[i], "-goroot") {
				continue
			}
			args = append(args, os.Args[i])
		}
		if *httpAddr == "" {
			args = append(args, "-http=:6060")
		}
		if err := installService(ServiceName, ServiceDesc, args...); err != nil {
			log.Fatalf("installService(%s, %s): %v", ServiceName, ServiceDesc, err)
		}
		fmt.Printf("Done\n")
		return
	}
	if *flagServiceUninstall {
		if err := removeService(ServiceName); err != nil {
			log.Fatalf("removeService: %v\n", err)
		}
		fmt.Printf("Done\n")
		return
	}
	if *flagServiceStart {
		if err := startService(ServiceName); err != nil {
			log.Fatalf("startService: %v\n", err)
		}
		fmt.Printf("Done\n")
		return
	}
	if *flagServiceStop {
		if err := controlService(ServiceName, svc.Stop, svc.Stopped); err != nil {
			log.Fatalf("stopService: %v\n", err)
		}
		fmt.Printf("Done\n")
		return
	}

	// Check usage: either server and no args, command line and args, or index creation mode
	if (*httpAddr != "" || *urlFlag != "") != (flag.NArg() == 0) && !*writeIndex {
		usage()
	}

	// run as service
	if isIntSess, err := svc.IsAnInteractiveSession(); err == nil && !isIntSess {
		runService(ServiceName)
		return
	}

	runGodoc()
}

type GodocService struct{}

func (m *GodocService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	go runGodoc()

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				// warning: unexpected control request ${c}
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string) {
	if err := svc.Run(name, &GodocService{}); err != nil {
		return
	}
}
