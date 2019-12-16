// Copyright (c) 2019, Peter Ohler, All rights reserved.

package gtt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// UseCase encapsulates a use case composed of multiple steps. The use case
// can be read from a JSON file in which case the Filepath member will be set.
type UseCase struct {

	// Comment is the description of the use case.
	Comment string

	// Filepath is the path to the file that the use case was read from or
	// will be written to.
	Filepath string

	// Steps are the steps to be taken in the use case.
	Steps []*Step

	runner *Runner
	memory map[string]interface{}
}

// NewUseCase creates a new UseCase from a file.
func NewUseCase(filepath string) (uc *UseCase, err error) {
	var data []byte

	if data, err = ioutil.ReadFile(filepath); err != nil {
		return
	}
	var m map[string]interface{}

	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	uc = &UseCase{Filepath: filepath}
	if uc.Comment, err = asString(m["comment"]); err != nil {
		return
	}
	if steps, ok := m["steps"].([]interface{}); ok {
		for _, v := range steps {
			step := Step{}
			if err = step.Set(v); err != nil {
				return
			}
			uc.Steps = append(uc.Steps, &step)
		}
	} else {
		return nil, fmt.Errorf("%T is not a valid steps type", m["steps"])
	}
	return
}

func (uc *UseCase) String() string {
	return string(uc.JSON())
}

func (uc *UseCase) JSON(indents ...int) []byte {
	indent := 0
	if 0 < len(indents) {
		indent = indents[0]
	}
	var j []byte
	if 0 < indent {
		j, _ = json.MarshalIndent(uc.Native(), "", strings.Repeat(" ", indent))
	} else {
		j, _ = json.Marshal(uc.Native())
	}
	return j
}

func (uc *UseCase) Native() interface{} {
	steps := make([]interface{}, 0, len(uc.Steps))
	for _, step := range uc.Steps {
		steps = append(steps, step.Native())
	}
	native := map[string]interface{}{
		"steps": steps,
	}
	if 0 < len(uc.Comment) {
		native["comment"] = easyString(uc.Comment)
	}
	return native
}

func (uc *UseCase) Run(r *Runner) (err error) {
	uc.runner = r
	// Start with a fresh memory cache as each run is separate from any other.
	uc.memory = map[string]interface{}{}
	path := uc.Filepath
	if !r.NoColor {
		path = underline + path + strings.Repeat(" ", (80-len(path))) + normal
	}
	if 0 < len(uc.Comment) {
		r.Log(aComment, "\n%s\n%s\n", path, uc.Comment)
	} else {
		r.Log(aComment, "\n%s\n", path)
	}
	for _, step := range uc.Steps {
		if err = step.Execute(uc); err != nil {
			break
		}
	}
	return
}
