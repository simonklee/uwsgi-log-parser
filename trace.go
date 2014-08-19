// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package trace

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func AnalyzePycall(dst io.Writer, src io.Reader) error {
	return nil
}

func AnalyzePyline(dst io.Writer, src io.Reader) error {
	res, err := pylineParser(src)
	if err != nil {
		return err
	}
	return analyze(res)
}

type result struct {
	lcnt       int
	aggregated map[string][]uint64
}

func newResult() *result {
	return &result{
		aggregated: make(map[string][]uint64),
	}
}

func analyze(res *result) error {
	var topk string
	var toptot, topmax, topmin, topmed uint64

	for k, v := range res.aggregated {
		var tot, min, max, med uint64

		for _, n := range v {
			tot += n
			if n > max {
				max = n
			}
			if n > min {
				min = n
			}
		}
		med = tot / uint64(len(v))

		if tot > toptot {
			topk = k
			toptot = tot
			topmin = min
			topmax = max
			topmed = med
		}
	}

	fmt.Printf("%s tot: %d, min: %d, max: %d, med: %d\n", topk, toptot, topmin, topmax, topmed)

	return nil
}

// [19/Aug/2014:19:09:22 +0000] - [uWSGI Python profiler 13] file /srv/test/src/ve/local/lib/python2.7/site-packages/werkzeug/wsgi.py line 694: close argc:1
func pylineParser(r io.Reader) (*result, error) {
	scanner := bufio.NewScanner(r)
	res := newResult()

	for scanner.Scan() {
		res.lcnt++
		line := scanner.Text()

		// find dt
		i := strings.Index(line, "uWSGI Python profiler ")

		if i == -1 {
			continue
		}

		//println(line)
		line = line[i+len("uWSGI Python profiler "):]
		i = strings.Index(line, "]")

		if i == -1 {
			continue
		}

		//usec := line[:i]
		usec, err := strconv.ParseUint(line[:i], 10, 0)

		if err != nil {
			continue
		}

		line = line[i+len("]")+len(" file "):]

		i = strings.Index(line, " ")

		if i == -1 {
			continue
		}

		filename := line[:i]
		line = line[len(filename)+len(" line "):]

		i = strings.Index(line, ": ")

		if i == -1 {
			continue
		}

		linenr := line[:i]
		line = line[i+len(": "):]

		// find func name
		i = strings.Index(line, " ")

		if i == -1 {
			continue
		}

		fname := line[:i]

		k := filename + ":" + linenr + ":" + fname + "()"
		v, ok := res.aggregated[k]

		if !ok {
			v = make([]uint64, 0, 8)
		}

		res.aggregated[k] = append(v, usec)
	}

	return res, nil
}
