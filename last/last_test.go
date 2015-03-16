// Copyright ©2012 The bíogo.external Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package last

import (
	"github.com/biogo/external"
	"gopkg.in/check.v1"
	"os/exec"
	"testing"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	for _, f := range []string{"lastdb", "lastal", "lastex"} {
		_, err := exec.LookPath(f)
		if err != nil {
			c.Skip("last suite not present")
		}
	}
}

func (s *S) TestSanityChecks(c *check.C) {
	for _, cb := range []external.CommandBuilder{
		DB{},
		Align{},
		Expect{},
	} {
		cmd, err := cb.BuildCommand()
		c.Check(cmd, check.Equals, (*exec.Cmd)(nil))
		c.Check(err, check.Equals, ErrMissingRequired)
	}
}

func (s *S) TestBuild(c *check.C) {
	for _, cb := range []external.CommandBuilder{
		DB{OutFile: "out", InFiles: []string{"in"}},
		Align{DB: "db", InFiles: []string{"in"}},
		Expect{Ref: "ref", Query: "query"},
	} {
		_, err := cb.BuildCommand()
		// c.Check(cmd, check.Equals, (*exec.Cmd)(nil))
		c.Check(err, check.Equals, nil)
	}
}
