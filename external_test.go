// Copyright ©2012 Dan Kortschak <dan.kortschak@adelaide.edu.au>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
	Glob []string `buildarg:"{{.}}"`                          // "<in>"...
}

func (l Ls) BuildCommand() (*exec.Cmd, error) {
	cl, err := Build(l)
	if err != nil {
		return nil, err
	}
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
	Cmd     string   `buildarg:"{{if .}}{{.}}{{else}}du{{end}}"` // du
	Exclude []string `buildarg:"{{if .}}--exclude={{.}}{{end}}"` // --exclude="file"...
}

func (d Du) BuildCommand() (*exec.Cmd, error) {
	cl, err := Build(d)
	if err != nil {
		return nil, err
	}
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
