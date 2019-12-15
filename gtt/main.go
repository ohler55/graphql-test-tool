// Copyright (c) 2019, Peter Ohler, All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ohler55/graphql-test-tool/gtt"
)

var verbose = false
var debug = false
var server = ""
var base = "/graphql"
var showComments = false
var showResponses = false
var showRequests = false
var noColor = false
var indent = 0

func init() {
	flag.StringVar(&server, "s", server, "server URL, host and port (example: http://localhost:8080)")
	flag.StringVar(&base, "b", base, "URL base path")
	flag.BoolVar(&showComments, "comment", showComments, "show comments")
	flag.BoolVar(&showRequests, "request", showRequests, "show requests")
	flag.BoolVar(&showResponses, "response", showResponses, "show responses")
	flag.BoolVar(&noColor, "no-color", noColor, "no color")
	flag.BoolVar(&verbose, "v", verbose, "verbose")
	flag.BoolVar(&debug, "d", debug, "debug")
	flag.IntVar(&indent, "i", indent, "indent")
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `

usage: %s [<options>] <json-file>...

`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprint(os.Stderr, "\n")
	}
	flag.Parse()

	r := gtt.Runner{
		Server:        server,
		Base:          base,
		ShowComments:  showComments,
		ShowResponses: showResponses,
		ShowRequests:  showRequests,
		NoColor:       noColor,
		Indent:        indent,
	}
	if verbose {
		r.ShowComments = true
		r.ShowResponses = true
		r.ShowRequests = true
	}
	for _, filepath := range flag.Args() {
		uc, err := gtt.NewUseCase(filepath)
		if err != nil {
			fmt.Printf("*-*-* Error: %s\n", err)
			os.Exit(1)
		}
		r.UseCases = append(r.UseCases, uc)
	}
	if debug {
		r.Log(gtt.Debug, string(r.JSON(2)))
	}
	if err := r.Run(); err != nil {
		fmt.Printf("*-*-* Error: %s\n", err)
		os.Exit(1)
	}
}

// TBD
// README
// usecase file description (format.md)
// doc.go - same as format.md
