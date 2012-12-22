// Copyright ©2012 The bíogo.external Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package external

import (
	"bytes"
	check "launchpad.net/gocheck"
	"os/exec"
	"sort"
	"strings"
	"testing"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

type Ls struct {
	Cmd  string   `buildarg:"{{if .}}{{.}}{{else}}ls{{end}}"` // ls
	Glob []string `buildarg:"{{if .}}{{args .}}{{end}}"`      // "<in>"...
}

func (l Ls) BuildCommand() (*exec.Cmd, error) {
	cl := Must(Build(l))
	return exec.Command(cl[0], cl[1:]...), nil
}

func (s *S) TestBuildLs(c *check.C) {
	_, err := Ls{}.BuildCommand()
	c.Check(err, check.Equals, nil)
}

func (s *S) TestLs(c *check.C) {
	_, err := exec.LookPath("ls")
	if err != nil {
		c.Skip("ls not present")
	}
	files := []string{"external.go", "external_test.go", "external_example_test.go"}
	sort.Strings(files)
	ls, err := Ls{
		Glob: files,
	}.BuildCommand()
	if err != nil {
		c.Check(err, check.Equals, nil) // Build Ls command.
	}
	ls.Stdout = &bytes.Buffer{}
	ls.Stderr = &bytes.Buffer{}
	err = ls.Run()
	if err != nil {
		c.Check(err, check.Equals, nil) // Run Ls command.
	}
	list := strings.Fields(ls.Stdout.(*bytes.Buffer).String())
	sort.Strings(list)
	c.Check(files, check.DeepEquals, list)
	c.Check(ls.Stderr.(*bytes.Buffer).String(), check.Equals, "")
}

type Du struct {
	Cmd     string   `buildarg:"{{if .}}{{.}}{{else}}du{{end}}"`                       // du
	Exclude []string `buildarg:"{{if .}}{{mprintf \"--exclude=%s\" . | args}}{{end}}"` // --exclude="file"...
}

func (d Du) BuildCommand() (*exec.Cmd, error) {
	cl := Must(Build(d))
	return exec.Command(cl[0], cl[1:]...), nil
}

func (s *S) TestBuildDu(c *check.C) {
	_, err := Du{}.BuildCommand()
	c.Check(err, check.Equals, nil)
}

func (s *S) TestDu(c *check.C) {
	_, err := exec.LookPath("du")
	if err != nil {
		c.Skip("du not present")
	}
	files := []string{
		"external.go",
		"external_test.go",
		"external_example_test.go",
		".git",
	}
	du, err := Du{
		Exclude: files,
	}.BuildCommand()
	if err != nil {
		c.Check(err, check.Equals, nil) // Build Du command.
	}
	c.Check(du.Args, check.DeepEquals, []string{
		"du",
		"--exclude=external.go",
		"--exclude=external_test.go",
		"--exclude=external_example_test.go",
		"--exclude=.git",
	})
	du.Stdout = &bytes.Buffer{}
	du.Stderr = &bytes.Buffer{}
	err = du.Run()
	if err != nil {
		c.Check(err, check.Equals, nil) // Run Du command.
	}
	list := strings.Fields(du.Stdout.(*bytes.Buffer).String())
	c.Check(list, check.DeepEquals, []string{
		"16", "./last",
		"20", "./muscle",
		"128", "./mafft",
		"164", ".",
	})
	c.Check(du.Stderr.(*bytes.Buffer).String(), check.Equals, "")
}
