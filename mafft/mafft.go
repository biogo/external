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
	"os/exec"
	"text/template"
)

type Mafft struct {
	// Usage: mafft <inputfile> > <outputfile>
	//
	// For details relating to options and parameters, see the MAFFT manual.
	//
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}mafft{{end}}"` // mafft

	// Algorithm:
	Auto          bool    `buildarg:"{{if .}}--auto{{end}}"`               // --auto
	HexamerPair   bool    `buildarg:"{{if .}}--6merpair{{end}}"`           // --6merpair
	GlobalPair    bool    `buildarg:"{{if .}}--globalpair{{end}}"`         // --globalpair
	LocalPair     bool    `buildarg:"{{if .}}--localpair{{end}}"`          // --localpair
	GenafPair     bool    `buildarg:"{{if .}}--genafpair{{end}}"`          // --genafpair
	FastaPair     bool    `buildarg:"{{if .}}--fastapair{{end}}"`          // --fastapair
	Weighting     float64 `buildarg:"{{with .}}--weighti {{.}}{{end}}"`    // --weighti <f.>
	ReTree        int     `buildarg:"{{with .}}--retree {{.}}{{end}}"`     // --retree <n>
	MaxIterate    int     `buildarg:"{{with .}}--maxiterate {{.}}{{end}}"` // --maxiterate <n>
	Fft           bool    `buildarg:"{{if .}}--fft{{end}}"`                // --fft
	NoFft         bool    `buildarg:"{{if .}}--nofft{{end}}"`              // --nofft
	NoScore       bool    `buildarg:"{{if .}}--noscore{{end}}"`            // --noscore
	MemSave       bool    `buildarg:"{{if .}}--memsave{{end}}"`            // --memsave
	Partree       bool    `buildarg:"{{if .}}--parttree{{end}}"`           // --parttree
	DPPartTree    bool    `buildarg:"{{if .}}--dpparttree{{end}}"`         // --dpparttree
	FastaPartTree bool    `buildarg:"{{if .}}--fastaparttree{{end}}"`      // --fastaparttree
	PartSize      int     `buildarg:"{{with .}}--partsize {{.}}{{end}}"`   // --partsize <n>
	GroupSize     int     `buildarg:"{{with .}}--groupsize {{.}}{{end}}"`  // --groupsize <n>

	// Parameter:
	GapOpenCost          float64 `buildarg:"{{with .}}--op {{.}}{{end}}"`           // --op <f.>
	ExtensionCost        float64 `buildarg:"{{with .}}--ep {{.}}{{end}}"`           // --ep <f.>
	LocalOpenCost        float64 `buildarg:"{{with .}}--lop {{.}}{{end}}"`          // --lop <f.>
	LocalPairOffset      float64 `buildarg:"{{with .}}--lep {{.}}{{end}}"`          // --lep <f.>
	LocalExtensionCost   float64 `buildarg:"{{with .}}--lexp {{.}}{{end}}"`         // --lexp <f.>
	GapOpenSkipCost      float64 `buildarg:"{{with .}}--LOP {{.}}{{end}}"`          // --LOP <f.>
	GapExtensionSkipCost float64 `buildarg:"{{with .}}--LEXP {{.}}{{end}}"`         // --LEXP <f.>
	Blosum               byte    `buildarg:"{{with .}}--bl {{.}}{{end}}"`           // --bl <n>
	JttPAM               uint    `buildarg:"{{with .}}--jtt {{.}}{{end}}"`          // --jtt <n>
	TransMembranePAM     uint    `buildarg:"{{with .}}--tm {{.}}{{end}}"`           // --tm <n>
	AminoMatrix          string  `buildarg:"{{with .}}--aamatrix \"{{.}}\"{{end}}"` // --aamatrix <file>
	FModel               bool    `buildarg:"{{if .}}--fmodel{{end}}"`               // --fmodel

	// Output:
	ClustalOut bool `buildarg:"{{if .}}--clustalout{{end}}"` // --clustalout
	InputOrder bool `buildarg:"{{if .}}--inputorder{{end}}"` // --inputorder
	Reorder    bool `buildarg:"{{if .}}--reorder{{end}}"`    // --reorder
	TreeOut    bool `buildarg:"{{if .}}--treeout{{end}}"`    // --treeout
	Quiet      bool `buildarg:"{{if .}}--quiet{{end}}"`      // --quiet

	// Input:
	Nucleic bool     `buildarg:"{{if .}}--nuc{{end}}"`                     // --nuc
	Amino   bool     `buildarg:"{{if .}}--amino{{end}}"`                   // --amino
	Seed    []string `buildarg:"{{mprintf \"--seed %q\" . | join \" \"}}"` // --seed <file>...

	// Files:
	InFile  string `buildarg:"{{with .}}{{softquote .}}{{end}}"` // <inputfile> - use "-" for stdin.
	OutFile string `buildarg:"{{with .}} >\"{{.}}\"{{end}}"`     // ><outfile>
}

func sq(s string) string {
	if s == "-" {
		return s
	}
	return `"` + s + `"`
}

func (m Mafft) BuildCommand() (*exec.Cmd, error) {
	cl, err := external.Build(m, template.FuncMap{"softquote": sq})
	if err != nil {
		return nil, err
	}
	return exec.Command(cl[0], cl[1:]...), nil
}
