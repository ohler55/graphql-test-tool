// Copyright (c) 2019, Peter Ohler, All rights reserved.

package gtt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// extracting from json/native

func asString(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	switch tv := value.(type) {
	case string:
		return tv, nil
	case []interface{}:
		sa := make([]string, 0, len(tv))
		for _, v := range tv {
			if s, ok := v.(string); ok {
				sa = append(sa, s)
			} else {
				return "", fmt.Errorf("%T is not a valid string element type", v)
			}
		}
		return strings.Join(sa, "\n"), nil
	}
	return "", fmt.Errorf("%T is not a valid string element type", value)
}

// convert to a string or map[string]interface{}
func asMapOrString(value interface{}) (interface{}, error) {
	switch tv := value.(type) {
	case string, []interface{}:
		return asString(value)
	case map[string]interface{}:
		return tv, nil
	}
	return nil, fmt.Errorf("%T is not a valid type for a string or map[string]interface{}", value)
}

func asMapStrStr(value interface{}) (map[string]string, error) {
	if m, _ := value.(map[string]interface{}); m != nil {
		mss := map[string]string{}
		for k, v := range m {
			if s, ok := v.(string); ok {
				mss[k] = s
			} else {
				return nil, fmt.Errorf("expected a string, not a %T", v)
			}
		}
		return mss, nil
	}
	return nil, fmt.Errorf("%T is not a valid type for a map[string]string", value)
}

// json building helpers

// if multiple lines then convert to an array
func easyString(s string) interface{} {
	if strings.Contains(s, "\n") {
		return strings.Split(s, "\n")
	}
	return s
}

func addNotNil(m map[string]interface{}, key string, value interface{}) {
	if value != nil {
		m[key] = value
	}
}

func addAny(m map[string]interface{}, key string, value interface{}) {
	if value != nil {
		if str, _ := value.(string); 0 < len(str) {
			m[key] = easyString(str)
		} else {
			m[key] = value
		}
	}
}

// compare results

// Returns path, result, expected
func match(result interface{}, expect interface{}) ([]string, interface{}, interface{}) {
	switch x := expect.(type) {
	case map[string]interface{}:
		if rm, ok := result.(map[string]interface{}); ok {
			checked := map[string]bool{}
			exact := false
			for k, v := range x {
				if k == "*" {
					exact = true
				} else {
					checked[k] = true
				}
				if loc, av, xv := match(rm[k], v); loc != nil {
					return append([]string{k}, loc...), av, xv
				}
			}
			if exact {
				for k, v := range rm {
					if !checked[k] && v != nil {
						return []string{k}, v, nil
					}
				}
			}
		} else {
			return []string{}, result, expect
		}
	case []interface{}:
		if ra, ok := result.([]interface{}); ok {
			for i, v := range x {
				if len(ra) <= i {
					return []string{}, nil, v
				}
				if loc, av, xv := match(ra[i], v); loc != nil {
					return append([]string{strconv.Itoa(i)}, loc...), av, xv
				}
			}
			if len(ra) > len(x) {
				return []string{}, ra, x
			}
		} else {
			return []string{}, result, expect
		}
	case string:
		match := false
		if 2 < len(x) && x[0] == '/' && x[len(x)-1] == '/' {
			if rs, ok := result.(string); ok {
				match, _ = regexp.MatchString(x[1:len(x)-2], rs)
			} else {
				match, _ = regexp.MatchString(x[1:len(x)-2], fmt.Sprintf("%v", result))
			}
		} else if rs, ok := result.(string); ok {
			match = (rs == x)
		}
		if !match {
			return []string{}, result, expect
		}
	case float64:
		switch r := result.(type) {
		case float64:
			if r != x {
				return []string{}, result, expect
			}
		case int64:
			if float64(r) != x {
				return []string{}, result, expect
			}
		case int:
			if float64(r) != x {
				return []string{}, result, expect
			}
		default:
			return []string{}, result, expect
		}
	case int64:
		switch r := result.(type) {
		case float64:
			if r != float64(x) {
				return []string{}, result, expect
			}
		case int64:
			if r != x {
				return []string{}, result, expect
			}
		case int:
			if int64(r) != x {
				return []string{}, result, expect
			}
		default:
			return []string{}, result, expect
		}
	case int:
		switch r := result.(type) {
		case float64:
			if r != float64(x) {
				return []string{}, result, expect
			}
		case int64:
			if r != int64(x) {
				return []string{}, result, expect
			}
		case int:
			if r != x {
				return []string{}, result, expect
			}
		default:
			return []string{}, result, expect
		}
	default:
		if result != expect {
			return []string{}, result, expect
		}
	}
	return nil, nil, nil
}
