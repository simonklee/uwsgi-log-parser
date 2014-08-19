// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/simonz05/trace"
	"github.com/simonz05/util/log"
)

var (
	help       = flag.Bool("h", false, "show help text")
	version    = flag.Bool("version", false, "show version number and exit")
	profiler   = flag.String("profiler", "pycall", "profiler type. pycall, pyline [default pycall]")
	cpuprofile = flag.String("debug.cpuprofile", "", "write cpu profile to file")
)

var Version = "0.1.0"

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [OPTIONS]

analayze uwsgi profiler logs.
`, os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.Println("start trace…")

	if *version {
		fmt.Fprintln(os.Stdout, Version)
		return
	}

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	switch *profiler {
	case "pycall":
		trace.AnalyzePycall(os.Stdout, os.Stdin)
	case "pyline":
		trace.AnalyzePyline(os.Stdout, os.Stdin)
	default:
		flag.Usage()
		os.Exit(1)
	}
}
