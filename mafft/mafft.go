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

// Package mafft implements interaction with the MAFFT multiple alignment tool.
// MAFFT is available from http://mafft.cbrc.jp/alignment/software/
package mafft

import (
	"code.google.com/p/biogo.external"
	"errors"
	"os/exec"
)

var ErrMissingRequired = errors.New("mafft: missing required argument")

type Mafft struct {
	// Usage: mafft <inputfile> > <outputfile>
	//
	// For details relating to options and parameters, see the MAFFT manual.
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}mafft{{end}}"` // mafft

	// Algorithm:
	Auto          bool    `buildarg:"{{if .}}--auto{{end}}"`              // --auto
	HexamerPair   bool    `buildarg:"{{if .}}--6merpair{{end}}"`          // --6merpair
	GlobalPair    bool    `buildarg:"{{if .}}--globalpair{{end}}"`        // --globalpair
	LocalPair     bool    `buildarg:"{{if .}}--localpair{{end}}"`         // --localpair
	GenafPair     bool    `buildarg:"{{if .}}--genafpair{{end}}"`         // --genafpair
	FastaPair     bool    `buildarg:"{{if .}}--fastapair{{end}}"`         // --fastapair
	Weighting     float64 `buildarg:"{{if .}}--weighti||{{.}}{{end}}"`    // --weighti <f.>
	ReTree        int     `buildarg:"{{if .}}--retree||{{.}}{{end}}"`     // --retree <n>
	MaxIterate    int     `buildarg:"{{if .}}--maxiterate||{{.}}{{end}}"` // --maxiterate <n>
	Fft           bool    `buildarg:"{{if .}}--fft{{end}}"`               // --fft
	NoFft         bool    `buildarg:"{{if .}}--nofft{{end}}"`             // --nofft
	NoScore       bool    `buildarg:"{{if .}}--noscore{{end}}"`           // --noscore
	MemSave       bool    `buildarg:"{{if .}}--memsave{{end}}"`           // --memsave
	Partree       bool    `buildarg:"{{if .}}--parttree{{end}}"`          // --parttree
	DPPartTree    bool    `buildarg:"{{if .}}--dpparttree{{end}}"`        // --dpparttree
	FastaPartTree bool    `buildarg:"{{if .}}--fastaparttree{{end}}"`     // --fastaparttree
	PartSize      int     `buildarg:"{{if .}}--partsize||{{.}}{{end}}"`   // --partsize <n>
	GroupSize     int     `buildarg:"{{if .}}--groupsize||{{.}}{{end}}"`  // --groupsize <n>

	// Parameter:
	GapOpenCost          float64 `buildarg:"{{if .}}--op||{{.}}{{end}}"`       // --op <f.>
	ExtensionCost        float64 `buildarg:"{{if .}}--ep||{{.}}{{end}}"`       // --ep <f.>
	LocalOpenCost        float64 `buildarg:"{{if .}}--lop||{{.}}{{end}}"`      // --lop <f.>
	LocalPairOffset      float64 `buildarg:"{{if .}}--lep||{{.}}{{end}}"`      // --lep <f.>
	LocalExtensionCost   float64 `buildarg:"{{if .}}--lexp||{{.}}{{end}}"`     // --lexp <f.>
	GapOpenSkipCost      float64 `buildarg:"{{if .}}--LOP||{{.}}{{end}}"`      // --LOP <f.>
	GapExtensionSkipCost float64 `buildarg:"{{if .}}--LEXP||{{.}}{{end}}"`     // --LEXP <f.>
	Blosum               byte    `buildarg:"{{if .}}--bl||{{.}}{{end}}"`       // --bl <n>
	JttPAM               uint    `buildarg:"{{if .}}--jtt||{{.}}{{end}}"`      // --jtt <n>
	TransMembranePAM     uint    `buildarg:"{{if .}}--tm||{{.}}{{end}}"`       // --tm <n>
	AminoMatrix          string  `buildarg:"{{if .}}--aamatrix||{{.}}{{end}}"` // --aamatrix <file>
	FModel               bool    `buildarg:"{{if .}}--fmodel{{end}}"`          // --fmodel

	// Output:
	ClustalOut bool `buildarg:"{{if .}}--clustalout{{end}}"` // --clustalout
	InputOrder bool `buildarg:"{{if .}}--inputorder{{end}}"` // --inputorder
	Reorder    bool `buildarg:"{{if .}}--reorder{{end}}"`    // --reorder
	TreeOut    bool `buildarg:"{{if .}}--treeout{{end}}"`    // --treeout
	Quiet      bool `buildarg:"{{if .}}--quiet{{end}}"`      // --quiet

	// Input:
	Nucleic bool     `buildarg:"{{if .}}--nuc{{end}}"`         // --nuc
	Amino   bool     `buildarg:"{{if .}}--amino{{end}}"`       // --amino
	Seed    []string `buildarg:"{{if .}}--seed||{{.}}{{end}}"` // --seed <file>...

	// Performance:
	Threads int `buildarg:"{{if .}}--thread||{{.}}{{end}}"` // --thread <n>

	// Files:
	InFile string `buildarg:"{{if .}}{{.}}{{end}}"` // <inputfile> - use "-" for stdin.
}

func (m Mafft) BuildCommand() (*exec.Cmd, error) {
	if m.InFile == "" {
		return nil, ErrMissingRequired
	}
	cl, err := external.Build(m)
	if err != nil {
		return nil, err
	}
	return exec.Command(cl[0], cl[1:]...), nil
}
