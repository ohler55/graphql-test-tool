// Copyright (c) 2019, Peter Ohler, All rights reserved.

package gtt

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ohler55/ojg/oj"
	"github.com/ohler55/ojg/sen"
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
	var p sen.Parser
	var v interface{}
	if v, err = p.Parse(data); err != nil {
		return
	}
	if m, _ = v.(map[string]interface{}); m == nil {
		return nil, fmt.Errorf("expected a map, not a %T", v)
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
		var p sen.Parser
		var pd interface{}
		if pd, err = p.Parse(data); err != nil {
			return err
		}
		if steps, _ = pd.([]interface{}); steps == nil {
			return fmt.Errorf("expected a array, not a %T", pd)
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
	return []byte(oj.JSON(uc, indent))
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

func (uc *UseCase) Simplify() interface{} {
	return uc.Native()
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

func (uc *UseCase) replaceVars(s string) string {
	for k, v := range uc.memory {
		pat := fmt.Sprintf("$(%s)", k)
		rep := fmt.Sprintf("%s", v)
		s = strings.ReplaceAll(s, pat, rep)
	}
	return s
}
