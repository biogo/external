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

package mafft

import (
	"bytes"
	check "launchpad.net/gocheck"
	"os/exec"
	"strings"
	"testing"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func (s *S) SetUpSuite(c *check.C) {
	_, err := exec.LookPath("mafft")
	if err != nil {
		c.Skip("mafft not present")
	}
}

func (s *S) TestBuild(c *check.C) {
	_, err := Mafft{}.BuildCommand()
	c.Check(err, check.Equals, nil)
}

func (s *S) TestSeed(c *check.C) {
	cmd, err := Mafft{Seed: []string{"a", "b", "c"}}.BuildCommand()
	c.Check(err, check.Equals, nil)
	c.Check(cmd.Args, check.DeepEquals, []string{"mafft", "--seed \"a\" --seed \"b\" --seed \"c\""})
}

func (s *S) TestMafft(c *check.C) {
	for _, t := range []struct{ in, out, err string }{
		{
			in: `>71.2259 lcl|scaffold_41:8288143+
CCCCAAATTCTCATAAAAAGACCAGACTTAATGGTCTGACTGAGACTAGAGGAATCCCGG
TGGTCATGGTCCCCAAACCTTCTGTTGGCCCAGGACAGGAACCATTCCCGAAGACAACTC
ATCAGACACGGAAGGGACTGGACAATGGGTAGGAGAGAGATGCTGACGAAGAGTGAGCTA
CTTGTATCAGGTGGACACTTGAGACTGTGTTGGCATCTCCTGTCTGGAGGGGAGATAGGA
GGGTAGAGAGGGTTAGAAACTGGCAAAATCGTCATGAAAGGAGGGACTGGAAGGAGGGAG
CGGGCTGACTCAGTAGGGGGAGAGTAAGTGGGAGTATGGAGTAAGGTGTATATAAGCTTA
TATGTGACAGATTGACTTGATTTGTAAACTTTCACTTAAAGCACAATAAAAATTATTTTT
TAAAAAATTGTTT
>71.2259 lcl|scaffold_41:11597466-
ATTATTATTTTTTTAAATAATTTTTATTGTGTTTTAAGGGAAAGTTTGCAAATCAAGTCA
GTCTCTCACATATAACCTTATATACACCTTACTCCATACTCCCATTTACTCTCCCCCTAA
TGAGTCAGCCCGCTCCCTCCTTCCGGTCTCTCCTTTCTTGACGATTTTGTCAGTTTCTAA
CCCTCTCTACCCTTCTATCTCTCCTCCAGACAGGAGATGCCAACACTGTCTCAAGTGTCC
ACTTGATACAAGTAGCTCACTCTTCGTCAGCATCTCTCTCCAACCCATTGTCCAGTCCCT
GCCATGTCTGATGAGTTGTCTTTGGGAATGGTTCCTGTCCTGGGCCAACAGAAGGTTTGG
GGACCATGACCGCTGGGATTCCTCTAGTCTCAGTCAGACCATTAAGTCTGGTCTTTTTAT
GAGA
>71.2259 lcl|scaffold_45:2724255+
ATAAAAAGACCAGACTTAATGGTCTGACTGAGACTAGAAGAATCCCGGTGGCCATGGTCC
CCAAACCTTCTGTTGGCCCAGGACAGGAACCATTCCCGAAGACAATTCATCAGACATGGA
AGGGACTGGACAATGGGTTGGAGAGAGATGCTGATAAAGAGTGAGCTACTTGTATCAGGT
GGACGTTTGAGACTGTATTGGCATCTCCTGTCTGGAGGGGAGATAGGGTAGAGAGGGTTA
GAAACTGGCAAAACGGTCACGAAAGGAGAGACTGGAAGAAGGGAGCAGGCTGACTCATTA
GGGGGAGAGTAAATGGGAGTATGTAGTAAGGTGTATATAAGCTTACATGTGACAGACTGA
CTTGATTTGTAAACTTTCACTTAAAGCACAATAAAAATTATTTTTTAAAAATTTGCC
`,
			out: `>71.2259 lcl|scaffold_41:8288143+
ccccaaattctcataaaaagaccagacttaatggtctgactgagactagaggaatcccgg
tggtcatggtccccaaaccttctgttggcccaggacagga--------------------
--------------------------------accattcccgaaga--------------
--------------------caactcatcagacacggaagggactggacaatgggtagga
gagagatgctgacgaagagtgagctacttgtatcaggtggacacttgaga----------
---ctgtgttggcatctcc---------------------------tgtctggaggggag
ataggagggtagagagggttagaaactggcaaaatcgtcatgaaaggagggactggaagg
agggagcgggctgactcagtagggggagagtaagtgggagtatggagtaaggtgtatata
agcttatatgtgacagattgacttgatttgtaaactttcacttaaagcacaataaaaatt
attttttaaaaaattgttt
>71.2259 lcl|scaffold_41:11597466-
attattatttttttaaataatttttatt--gtgttttaagggaaagtttgcaaatcaagt
cagtctctcacatataaccttatatacaccttactccatactcccatttactctccccct
aatgagtcagcccgctccctccttccggtctctcctttcttgacgattttgtcagtttct
aaccctctctacccttctatctctcctccagaca--------------------------
-ggagatgccaaca-------------ctgtctcaagtgtccacttgatacaagtagctc
actcttcgtcagcatctctctccaacccattgtccagtccctgccatgtctgatgagttg
tc-------------------------------------------------tttgggaat
ggttcctgtcctgggccaacagaaggtttg-----gggaccat-----------------
---------------gaccgctgggattcctctagtctcagtcagaccattaagtctggt
ctttttatgaga-------
>71.2259 lcl|scaffold_45:2724255+
------------ataaaaagaccagacttaatggtctgactgagactagaagaatcccgg
tggccatggtccccaaaccttctgttggcccaggacagga--------------------
--------------------------------accattcccgaaga--------------
--------------------caattcatcagacatggaagggactggacaatgggttgga
gagagatgctgataaagagtgagctacttgtatcaggtggacgtttgaga----------
---ctgtattggcatctcc---------------------------tgtctggaggggag
at---agggtagagagggttagaaactggcaaaacggtcacgaaaggagagactggaaga
agggagcaggctgactcattagggggagagtaaatgggagtatgtagtaaggtgtatata
agcttacatgtgacagactgacttgatttgtaaactttcacttaaagcacaataaaaatt
attttttaaaaatttgcc-
`,
			err: "\n" +
				"nseq =  3\n" +
				"distance =  ktuples\n" +
				"iterate =  0\n" +
				"cycle =  2\n" +
				"nthread = 0\n" +
				"generating 200PAM scoring matrix for nucleotides ... done\n" +
				"done\n" +
				"done\n" +
				"scoremtx = -1\n" +
				"Gap Penalty = -1.53, +0.00, +0.00\n" +
				"\n" +
				"tuplesize = 6, dorp = d\n" +
				"\n" +
				"\n" +
				"Making a distance matrix ..\n" +
				"\r    1 / 3\n" +
				"done.\n" +
				"\n" +
				"Constructing a UPGMA tree ... \n" +
				"\r    0 / 3\n" +
				"done.\n" +
				"\n" +
				"Progressive alignment ... \n" +
				"\rSTEP     1 / 2 f\rSTEP     2 / 2 f\n" +
				"done.\n" +
				"\n" +
				"disttbfast (nuc) Version 7.012b alg=A, model=DNA200 (2),  1.530 ( 4.590), -0.000 (-0.000)\n" +
				"0 thread(s)\n" +
				"nthread = 0\n" +
				"blosum 62 / kimura 200\n" +
				"generating 200PAM scoring matrix for nucleotides ... done\n" +
				"done\n" +
				"done\n" +
				"scoremtx = -1\n" +
				"Gap Penalty = -1.53, +0.00, +0.00\n" +
				"Making a distance matrix .. \n" +
				"\r    0 / 2\n" +
				"done.\n" +
				"\n" +
				"Constructing a UPGMA tree ... \n" +
				"\r    0 / 3\n" +
				"done.\n" +
				"\n" +
				"Progressive alignment ... \n" +
				"\rSTEP     1 /2 f\rSTEP     2 /2 f\n" +
				"done.\n" +
				"tbfast (nuc) Version 7.012b alg=A, model=DNA200 (2),  1.530 ( 4.590), -0.000 (-0.000)\n" +
				"0 thread(s)\n" +
				"\n" +
				"\n" +
				"Strategy:\n" +
				" FFT-NS-2 (Fast but rough)\n" +
				" Progressive method (guide trees were built 2 times.)\n" +
				"\n" +
				"If unsure which option to use, try 'mafft --auto input > output'.\n" +
				"For more information, see 'mafft --help', 'mafft --man' and the mafft page.\n" +
				"\n",
		},
	} {
		cmd, err := Mafft{InFile: "-"}.BuildCommand()
		c.Logf("%#v\n", cmd.Args)
		c.Check(err, check.Equals, nil)
		cmd.Stdin = strings.NewReader(t.in)
		bOut := &bytes.Buffer{}
		bErr := &bytes.Buffer{}
		cmd.Stdout = bOut
		cmd.Stderr = bErr
		cmd.Run()
		c.Check(bOut.String(), check.Equals, t.out)
		c.Check(bErr.String(), check.Equals, t.err)
	}
}
