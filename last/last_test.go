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

package last

import (
	"code.google.com/p/biogo.external"
	check "launchpad.net/gocheck"
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
