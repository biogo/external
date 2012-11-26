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

// Package muslce implements interaction with the MUSCLE multiple alignment tool.
// MUSCLE is available from http://www.drive5.com/muscle/
package muscle

import (
	"code.google.com/p/biogo.external"
	"os/exec"
	"text/template"
	"time"
)

type Log struct {
	File   string
	Append bool
}
type Muscle struct {
	// Usage: muscle -in <inputfile> -out <outputfile>
	//
	// Common options (for a complete list please see the User Guide):
	//
	//     -in <inputfile>    Input file in FASTA format (default stdin)
	//     -out <outputfile>  Output alignment in FASTA format (default stdout)
	//     -diags             Find diagonals (faster for similar sequences)
	//     -maxiters <n>      Maximum number of iterations (integer, default 16)
	//     -maxhours <h>      Maximum time to iterate in hours (default no limit)
	//     -html              Write output in HTML format (default FASTA)
	//     -msf               Write output in GCG MSF format (default FASTA)
	//     -clw               Write output in CLUSTALW format (default FASTA)
	//     -clwstrict         As -clw, with 'CLUSTAL W (1.81)' header
	//     -log[a] <logfile>  Log to file (append if -loga, overwrite if -log)
	//     -quiet             Do not write progress messages to stderr
	//     -version           Display version information and exit
	//
	// Without refinement (very fast, avg accuracy similar to T-Coffee): -maxiters 2
	// Fastest possible (amino acids): -maxiters 1 -diags -sv -distance1 kbit20_3
	// Fastest possible (nucleotides): -maxiters 1 -diags
	Cmd string `buildarg:"{{if .}}{{.}}{{else}}muscle{{end}}"` // muscle

	// Files:
	InFile  string `buildarg:"{{with .}}-in \"{{.}}\"{{end}}"`                              // -in <inputfile>
	OutFile string `buildarg:"{{with .}}-out \"{{.}}\"{{end}}"`                             // -out <outputfile>
	Log     Log    `buildarg:"{{if .File}}-log{{if .Append}}a{{end}} \"{{.File}}\"{{end}}"` // -log[a] <logfile>
	Quiet   bool   `buildarg:"{{if .}}-quiet{{end}}"`                                       // -quiet

	// Formatting:
	Html          bool `buildarg:"{{if .}}-html{{end}}"`      // -html
	Msf           bool `buildarg:"{{if .}}-msf{{end}}"`       // -msf
	Clustal       bool `buildarg:"{{if .}}-clw{{end}}"`       // -clw
	ClustalStrict bool `buildarg:"{{if .}}-clwstrict{{end}}"` // -clwstrict

	// Common options:
	FindDiagonals bool          `buildarg:"{{if .}}-diags{{end}}"`                  // -diags
	MaxIterations int           `buildarg:"{{with .}}-maxiters {{.}}{{end}}"`       // -maxiters <n>
	MaxDuration   time.Duration `buildarg:"{{with .}}-maxhours {{hours .}}{{end}}"` // -maxhours <h>

	// Other value options (see MUSCLE user guide):
	// Gleaned from user guide - may not reflect reality.
	AnchorSpacing   int     `buildarg:"{{with .}}-anchorspacing {{.}}{{end}}"`   // -anchorspacing <n>
	Center          float64 `buildarg:"{{with .}}-center {{.}}{{end}}"`          // -center <f.>
	Cluster1        string  `buildarg:"{{with .}}-cluster1 \"{{.}}\"{{end}}"`    // -cluster1 "upgma|upgma|neighborjoining"
	Cluster2        string  `buildarg:"{{with .}}-cluster2 \"{{.}}\"{{end}}"`    // -cluster2 "upgma|upgma|neighborjoining"
	ClustalOut      string  `buildarg:"{{with .}}-clwout \"{{.}}\"{{end}}"`      // -clwout <file>
	DiagonalBreak   int     `buildarg:"{{with .}}-diagbreak {{.}}{{end}}"`       // -diagbreak <n>
	DiagonalLength  int     `buildarg:"{{with .}}-diaglength {{.}}{{end}}"`      // -diaglength <n>
	DiagonalMargin  int     `buildarg:"{{with .}}-diagmargin {{.}}{{end}}"`      // -diagmargin <n>
	Distance1       string  `buildarg:"{{with .}}-distance1 \"{{.}}\"{{end}}"`   // -distance1 "kmer6_6|kmer20_3|kmer20_4|kbit20_3|kmer4_6"
	Distance2       string  `buildarg:"{{with .}}-distance2 \"{{.}}\"{{end}}"`   // -distance2 "pctid_kimura|pctid_log"
	FastaOut        string  `buildarg:"{{with .}}-fastaout \"{{.}}\"{{end}}"`    // -fastaout <file>
	GapOpen         float64 `buildarg:"{{with .}}-gapopen {{.}}{{end}}"`         // -gapopen <f.>
	GapExtend       float64 `buildarg:"{{with .}}-gapextend {{.}}{{end}}"`       // -gapextend <f.>
	HydroWindow     int     `buildarg:"{{with .}}-hydro {{.}}{{end}}"`           // -hydro <n>
	HydroFactor     float64 `buildarg:"{{with .}}-hydrofactor {{.}}{{end}}"`     // -hydrofactor <f.>
	In1             string  `buildarg:"{{with .}}-in1 \"{{.}}\"{{end}}"`         // -in1 <file>
	In2             string  `buildarg:"{{with .}}-in2 \"{{.}}\"{{end}}"`         // -in2 <file>
	Matrix          string  `buildarg:"{{with .}}-matrix \"{{.}}\"{{end}}"`      // -matrix <file>
	MaxTrees        int     `buildarg:"{{with .}}-maxtrees {{.}}{{end}}"`        // -maxtrees <n>
	MinBestColScore float64 `buildarg:"{{with .}}-minbestcolscore {{.}}{{end}}"` // -minbestcolscore <f.>
	MinSmoothScore  float64 `buildarg:"{{with .}}-minsmoothscore {{.}}{{end}}"`  // -minsmoothscore <f.>
	MsaOut          string  `buildarg:"{{with .}}-msaout \"{{.}}\"{{end}}"`      // -msaout <file>
	ObjectiveScore  string  `buildarg:"{{with .}}-objscore \"{{.}}\"{{end}}"`    // -objscore "sp|ps|dp|xp|spf|spm"
	PhyInterOut     string  `buildarg:"{{with .}}-phyiout \"{{.}}\"{{end}}"`     // -phyiout <file>
	PhySequenOut    string  `buildarg:"{{with .}}-physout \"{{.}}\"{{end}}"`     // -physout <file>
	RefineWindow    int     `buildarg:"{{with .}}-refinewindow {{.}}{{end}}"`    // -refinewindow <n>
	Root1           string  `buildarg:"{{with .}}-root1 \"{{.}}\"{{end}}"`       // -root1 "pseudo|midlongestspan|minavgleafdist"
	Root2           string  `buildarg:"{{with .}}-root2 \"{{.}}\"{{end}}"`       // -root2 "pseudo|midlongestspan|minavgleafdist"
	ScoreFile       string  `buildarg:"{{with .}}-scorefile \"{{.}}\"{{end}}"`   // -scorefile <file>
	SeqType         string  `buildarg:"{{with .}}-seqtype \"{{.}}\"{{end}}"`     // -seqtype "protein|nucleo|auto"
	SmoothScoreCeil float64 `buildarg:"{{with .}}-smoothscoreceil {{.}}{{end}}"` // -smoothscoreceil <f.>
	SmoothWindow    int     `buildarg:"{{with .}}-smoothwindow {{.}}{{end}}"`    // -smoothwindow <n>
	SpScore         string  `buildarg:"{{with .}}-spscore \"{{.}}\"{{end}}"`     // -spscore <file>
	Tree1           string  `buildarg:"{{with .}}-tree1 \"{{.}}\"{{end}}"`       // -tree1 <file>
	Tree2           string  `buildarg:"{{with .}}-tree2 \"{{.}}\"{{end}}"`       // -tree2 <file>
	UseTree         string  `buildarg:"{{with .}}-usetree \"{{.}}\"{{end}}"`     // -usetree <file>
	Weight1         string  `buildarg:"{{with .}}-weight1 \"{{.}}\"{{end}}"`     // -weight1 "none|henikoff|henikoffpb|gsc|clustalw|threeway"
	Weight2         string  `buildarg:"{{with .}}-weight2 \"{{.}}\"{{end}}"`     // -weight2 "none|henikoff|henikoffpb|gsc|clustalw|threeway"

	// Other flag options (see MUSCLE user guide):
	// Gleaned from user guide - may not reflect reality.
	Anchors        bool `buildarg:"{{if .}}-anchors{{end}}"`   // -anchors
	Brenner        bool `buildarg:"{{if .}}-brenner{{end}}"`   // -brenner
	Cluster        bool `buildarg:"{{if .}}-cluster{{end}}"`   // -cluster
	Dimer          bool `buildarg:"{{if .}}-dimer{{end}}"`     // -dimer
	Core           bool `buildarg:"{{if .}}-core{{end}}"`      // -core
	Diags1         bool `buildarg:"{{if .}}-diags1{{end}}"`    // -diags1
	Diags2         bool `buildarg:"{{if .}}-diags2{{end}}"`    // -diags2
	Fasta          bool `buildarg:"{{if .}}-fasta{{end}}"`     // -fasta
	Group          bool `buildarg:"{{if .}}-group{{end}}"`     // -group
	LogExpectation bool `buildarg:"{{if .}}-le{{end}}"`        // -le
	NoAnchors      bool `buildarg:"{{if .}}-noanchors{{end}}"` // -noanchors
	NoCore         bool `buildarg:"{{if .}}-nocore{{end}}"`    // -nocore
	PhylipInter    bool `buildarg:"{{if .}}-phyi{{end}}"`      // -phyi
	PhylipSequen   bool `buildarg:"{{if .}}-phys{{end}}"`      // -phys
	Profile        bool `buildarg:"{{if .}}-profile{{end}}"`   // -profile
	Refine         bool `buildarg:"{{if .}}-refine{{end}}"`    // -refine
	RefineByWindow bool `buildarg:"{{if .}}-refinew{{end}}"`   // -refinew
	SumOfPairsProt bool `buildarg:"{{if .}}-sp{{end}}"`        // -sp
	PPScore        bool `buildarg:"{{if .}}-ppscore{{end}}"`   // -ppscore
	SumOfPairsNuc  bool `buildarg:"{{if .}}-spn{{end}}"`       // -spn
	SumOfPairsProf bool `buildarg:"{{if .}}-sv{{end}}"`        // -sv
	verbose        bool `buildarg:"{{if .}}-verbose{{end}}"`   // -verbose
}

func hours(d time.Duration) float64 {
	return d.Hours()
}

func (m Muscle) BuildCommand() (*exec.Cmd, error) {
	cl, err := external.Build(m, template.FuncMap{"hours": hours})
	if err != nil {
		return nil, err
	}
	return exec.Command(cl[0], cl[1:]...), nil
}
