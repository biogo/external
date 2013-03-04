// Copyright ©2013 The bíogo.external Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package kmeans implements interaction with the kmeans clustering tool.
// The kmeans tool is available from http://www.cs.cmu.edu/~dpelleg/kmeans.html.
package kmeans

import (
	"code.google.com/p/biogo.external"

	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ErrMissingRequired = errors.New("kmeans: missing required argument")
	ErrBadInput        = errors.New("kmeans: bad input")
	ErrNoUniverse      = errors.New("kmeans: no universe file")
)

type MakeUniverse struct {
	// Usage: kmeans makeuni in <infile>
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}kmeans{{end}}"` // kmeans

	// Make Universe:
	Kmeans struct{} `buildarg:"makeuni"` // makeuni

	// Files:
	Infile string `buildarg:"{{if .}}in{{split}}{{.}}{{end}}"` // in <file>
}

func (u MakeUniverse) BuildCommand() (*exec.Cmd, error) {
	if u.Infile == "" {
		return nil, ErrMissingRequired
	}
	cl := external.Must(external.Build(u))
	return exec.Command(cl[0], cl[1:]...), nil
}

type Xmeans struct {
	// Usage: kmeans kmeans [options] -in <infile>
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}kmeans{{end}}"` // kmeans

	// Kmeans:
	Kmeans struct{} `buildarg:"kmeans"` // kmeans

	// Files:
	InFile           string `buildarg:"{{if .}}-in{{split}}{{.}}{{end}}"`                     // -in <file>
	Constraints      string `buildarg:"{{if .}}-cons{{split}}{{.}}{{end}}"`                   // -cons <file>
	InitCenters      string `buildarg:"{{if .}}-init_ctrs{{split}}{{.}}{{end}}"`              // -init_ctrs <file>
	SaveCenters      string `buildarg:"{{if .}}-save_ctrs{{split}}{{.}}{{end}}"`              // -save_ctrs <file>
	PrintClusters    string `buildarg:"{{if .}}-printclusters{{split}}{{.}}{{end}}"`          // -printclusters <file>
	PrintNonClusters string `buildarg:"{{if .}}-print_no_cons_clusters{{split}}{{.}}{{end}}"` // -print_no_cons_clusters <file>

	// Options:
	InitialK         int     `buildarg:"{{if .}}-k{{split}}{{.}}{{end}}"`                     // -k <int>
	MaxCenters       int     `buildarg:"{{if .}}-max_ctrs{{split}}{{.}}{{end}}"`              // -max_ctrs <int>
	Method           string  `buildarg:"{{if .}}-method{{split}}{{.}}{{end}}"`                // -method <string>
	Splits           int     `buildarg:"{{if .}}-num_splits{{split}}{{.}}{{end}}"`            // -num_splits <int>
	DelSteps         int     `buildarg:"{{if .}}-del_steps_ratio{{split}}{{.}}{{end}}"`       // -del_steps_ratio <int>
	MaxIterations    int     `buildarg:"{{if .}}-max_iter{{split}}{{.}}{{end}}"`              // -max_iter <int>
	CutoffFactor     float64 `buildarg:"{{if .}}-cutoff_factor{{split}}{{.}}{{end}}"`         // -cutoff_factor <float>
	MaxLeafSize      int     `buildarg:"{{if .}}-max_leaf_size{{split}}{{.}}{{end}}"`         // -max_leaf_size <int>
	MinBoxWidth      float64 `buildarg:"{{if .}}-min_box_width{{split}}{{.}}{{end}}"`         // -min_box_width <float>
	CreateUniverse   bool    `buildarg:"{{if .}}-create_universe{{split}}true{{end}}"`        // -create_universe <bool>
	NeverKillCenters bool    `buildarg:"{{if .}}-never_kill_ctrs{{split}}true{{end}}"`        // -never_kill_ctrs <bool>
	SplitStatistic   string  `buildarg:"{{if .}}-S{{split}}{{.}}{{end}}"`                     // -S <string>
	Seed             int     `buildarg:"{{if .}}-seed{{split}}{{.}}{{end}}"`                  // -seed <int>
	RandStart        bool    `buildarg:"{{if .}}-randstart{{split}}true{{end}}"`              // -randstart <bool>
	ForceSplitFrac   float64 `buildarg:"{{if .}}-forced_split_fraction{{split}}{{.}}{{end}}"` // -forced_split_fraction <float>
	SplitConfLevel   float64 `buildarg:"{{if .}}-split_conf_level{{split}}{{.}}{{end}}"`      // -split_conf_level <float>

	// Display options:
	ShowEndCenters bool `buildarg:"{{if .}}-D_SHOW_END_CENTERS{{end}}"` // -D_SHOW_END_CENTERS
	DrawPoints     bool `buildarg:"{{if .}}-D_DRAWPOINTS{{end}}"`       // -D_DRAWPOINTS
	Interactive    bool `buildarg:"{{if .}}-D_INTERACTIVE{{end}}"`      // -D_INTERACTIVE
	ShowBValue     bool `buildarg:"{{if .}}-D_SHOW_BVALUE{{end}}"`      // -D_SHOW_BVALUE
}

func (x Xmeans) BuildCommand() (*exec.Cmd, error) {
	if x.InFile == "" {
		return nil, ErrMissingRequired
	}

	// We should be able to let kmeans fail this, but it waits for user input so
	// we just don't go there. Sane behaviour would require people patch their
	// kmeans source:
	//
	//  change void my_error_default(const char *string) in kmeans/utils/ambs.c
	//  to something less drastic.
	//
	if !x.CreateUniverse {
		_, err := os.Open(x.InFile + ".universe")
		if err != nil {
			return nil, ErrNoUniverse
		}
	}

	cl := external.Must(external.Build(x))
	return exec.Command(cl[0], cl[1:]...), nil
}

func Membership(r io.Reader) ([]int, error) {
	var (
		currId  = -1
		lastID  = -1
		seen    = make(map[int]int)
		classes = make(map[int]int)
		pnts    int
	)

	buf := bufio.NewReader(r)
	for {
		l, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		l = strings.TrimSpace(l)

		if len(l) < 2 || l[0] != '+' {
			return nil, ErrBadInput
		}

		if l[1] == '#' {
			l = strings.TrimSpace(l[2:])
			cname, err := strconv.Atoi(l)
			if err != nil {
				return nil, err
			}
			var ok bool
			if currId, ok = seen[cname]; !ok {
				lastID++
				currId = lastID
				seen[cname] = currId
			}
			continue
		}

		if currId < 0 {
			return nil, ErrBadInput
		}
		l = strings.TrimSpace(l[1:])
		pnt, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}
		pnts++
		classes[pnt] = currId
	}

	if lastID <= 0 {
		return nil, nil
	}

	mi := make([]int, pnts)
	for i := range mi {
		v, ok := classes[i]
		if ok {
			mi[i] = v
		} else {
			mi[i] = -1
		}
	}

	return mi, nil
}
