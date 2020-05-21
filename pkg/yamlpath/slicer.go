/*
 * Copyright 2020 Go YAML Path Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package yamlpath

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func slice(index string, length int) ([]int, error) {
	if union := strings.Split(index, ","); len(union) > 1 {
		combination := []int{}
		for i, idx := range union {
			sl, err := slice(idx, length)
			if err != nil {
				return nil, fmt.Errorf("error in union member %d: %s", i, err)
			}
			combination = combine(combination, sl)
		}
		return combination, nil
	}
	from := 0
	step := 1
	var to int
	if index == "*" {
		to = length
	} else {
		sliceParms := strings.Split(index, ":")
		if len(sliceParms) > 3 {
			return nil, errors.New("malformed array index")
		}
		p := []int{}
		for i, s := range sliceParms {
			if i == 0 && s == "" {
				p = append(p, 0)
				continue
			}
			if i == 1 && s == "" {
				p = append(p, length)
				continue
			}
			if i == 2 && s == "" {
				p = append(p, 1)
				continue
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, errors.New("non-integer array index")
			}
			p = append(p, n)
		}
		from = p[0]
		if from < 0 {
			from = length + from
			to = from - 1
			step = -1
		} else {
			to = from + 1
		}
		if len(p) >= 2 {
			if p[1] >= 0 {
				to = p[1]
			} else {
				to = length + p[1]
			}
			if from < to {
				step = 1
			}
		}
		if len(p) == 3 {
			step = p[2]
		}
		if step < 0 && from <= to {
			from, to = to-1, from-1
		}
	}
	slice := []int{}
	if step > 0 {
		for i := from; i < to; i += step {
			slice = append(slice, i)
		}
	} else if step < 0 {
		for i := from; i > to; i += step {
			slice = append(slice, i)
		}
	}
	return slice, nil
}

func combine(a, b []int) []int {
	if a == nil {
		return b
	}
	if b == nil {
		return b
	}
	present := map[int]struct{}{}
	result := []int{}
	for _, i := range a {
		if _, pres := present[i]; !pres {
			result = append(result, i)
			present[i] = struct{}{}
		}
	}
	for _, i := range b {
		if _, pres := present[i]; !pres {
			result = append(result, i)
			present[i] = struct{}{}
		}
	}
	return result
}
