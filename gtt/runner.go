// Copyright (c) 2019, Peter Ohler, All rights reserved.

package gtt

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	aComment  = ""
	aRequest  = "\x1b[36m"   // dark cyan
	aResponse = "\x1b[32;1m" // green
	underline = "\x1b[4m"
	normal    = "\x1b[m" // back to normal
	// Debug for debug logging.
	Debug = "\x1b[35m" // dark cyan
	red   = "\x1b[31m" // red
)

// Runner runs UseCases. It provide focal point for an assembly of use case
// tests. Test output is feature based. A flag is available for each category
// of output.
type Runner struct {

	// Server is the host and port of a URL. Example: http://localhost:8080
	Server string

	// Base of the URL path. By convention this is usually "/graphql".
	Base string

	// ShowComments if true will cause comments to be printed during a run.
	ShowComments bool

	// ShowRequests if true will cause request URL and content to be printed
	// during a run.
	ShowRequests bool

	// ShowResponses if true will cause response JSON to be printed during a run.
	ShowResponses bool

	// NoColor if true turns off colorized output.
	NoColor bool

	// Indent is the number of spaces to indent JSON resonses. If 0 no
	// modification to responses are made otherwise the JSON is unmarshalled
	// and re-marshalled with an indentation.
	Indent int

	// UseCases to run.
	UseCases []*UseCase

	// Writer is an alternate io.Writer that will be used in place of writing
	// to Stdout when logging if not nil.
	Writer io.Writer
}

// Run the usecases.
func (r *Runner) Run() (err error) {
	for _, uc := range r.UseCases {
		if err = uc.Run(r); err != nil {
			break
		}
	}
	if r.ShowComments || r.ShowRequests || r.ShowResponses {
		fmt.Println()
	}
	return
}

func (r *Runner) String() string {
	return string(r.JSON())
}

func (r *Runner) JSON(indents ...int) []byte {
	indent := 0
	if 0 < len(indents) {
		indent = indents[0]
	}
	var j []byte
	if 0 < indent {
		j, _ = json.MarshalIndent(r.Native(), "", strings.Repeat(" ", indent))
	} else {
		j, _ = json.Marshal(r.Native())
	}
	return j
}

// Native version of the Runner. Used for JSON() which is mostly for
// debugging.
func (r *Runner) Native() interface{} {
	cases := make([]interface{}, 0, len(r.UseCases))
	for _, uc := range r.UseCases {
		cases = append(cases, uc.Native())
	}
	native := map[string]interface{}{
		"useCases":      cases,
		"server":        r.Server,
		"base":          r.Base,
		"showComments":  r.ShowComments,
		"showRequests":  r.ShowRequests,
		"showResponses": r.ShowResponses,
		"noColor":       r.NoColor,
		"indent":        r.Indent,
	}
	return native
}

// Log output for one of the categories.
func (r *Runner) Log(color string, format string, args ...interface{}) {
	switch color {
	case aComment:
		if !r.ShowComments {
			return
		}
	case aRequest:
		if !r.ShowRequests {
			return
		}
	case aResponse:
		if !r.ShowResponses {
			return
		}
	}
	format += "\n"
	if !r.NoColor && color != aComment {
		format = color + format + normal
	}
	if r.Writer != nil {
		_, _ = r.Writer.Write([]byte(fmt.Sprintf(format, args...)))
	} else {
		fmt.Printf(format, args...)
	}
}
