// Copyright ©2013 The bíogo.external Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmeans

import (
	"code.google.com/p/biogo.external"

	"gopkg.in/check.v1"
	"os/exec"
	"strings"
	"testing"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	_, err := exec.LookPath("kmeans")
	if err != nil {
		c.Skip("kmeans not present")
	}
}

func (s *S) TestSanityChecks(c *check.C) {
	for _, t := range []struct {
		cb  external.CommandBuilder
		err error
	}{
		{MakeUniverse{}, ErrMissingRequired},
		{Xmeans{}, ErrMissingRequired},
		{Xmeans{InFile: "test"}, ErrNoUniverse},
	} {
		cmd, err := t.cb.BuildCommand()
		c.Check(cmd, check.Equals, (*exec.Cmd)(nil))
		c.Check(err, check.Equals, t.err)
	}
}

func (s *S) TestBuild(c *check.C) {
	for _, cb := range []external.CommandBuilder{
		MakeUniverse{},
		Xmeans{},
	} {
		cmd, err := external.Build(cb)
		c.Check(cmd, check.Not(check.Equals), (*exec.Cmd)(nil))
		c.Check(err, check.Equals, nil)
	}
}

func (s *S) TestMembership(c *check.C) {
	mi, err := Membership(strings.NewReader(printclusters))
	c.Check(err, check.Equals, nil)
	c.Check(mi, check.DeepEquals, membership)
}
