// Copyright (c) 2019, Peter Ohler, All rights reserved.

package gtt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	if err = uc.addSteps(m["steps"]); err != nil {
		return
	}
	return
}

// The arg can be either a string, array, or a map. A map is assumed to be a
// single step while a string is a relative path to a file to include. The
// included file should be an array of steps or steps and additional includes.
func (uc *UseCase) addSteps(v interface{}) error {
	switch tv := v.(type) {
	case []interface{}:
		for _, v = range tv {
			if err := uc.addSteps(v); err != nil {
				return err
			}
		}
	case string:
		filepath := filepath.Join(filepath.Dir(uc.Filepath), tv)
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			return err
		}
		var steps []interface{}

		if err = json.Unmarshal(data, &steps); err != nil {
			return err
		}
		return uc.addSteps(steps)
	case map[string]interface{}:
		step := Step{}
		if err := step.Set(v); err != nil {
			return err
		}
		uc.Steps = append(uc.Steps, &step)
	default:
		return fmt.Errorf("%T is not a valid steps type", v)
	}
	return nil
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
		if 80 <= len(path) {
			path = underline + path + normal
		} else {
			path = underline + path + strings.Repeat(" ", (80-len(path))) + normal
		}
	}
	if 0 < len(uc.Comment) {
		r.Log(aComment, "\n%s\n%s\n", path, uc.Comment)
	} else {
		r.Log(aComment, "\n%s\n", path)
	}
	for _, step := range uc.Steps {
		if err == nil {
			err = step.Execute(uc)
		} else if step.Always {
			_ = step.Execute(uc)
		}
	}
	return
}
