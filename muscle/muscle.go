// Copyright ©2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package muscle implements interaction with the MUSCLE multiple alignment tool.
// MUSCLE is available from http://www.drive5.com/muscle/
package muscle

import (
	"os/exec"
	"text/template"
	"time"

	"github.com/biogo/external"
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
	InFile  string `buildarg:"{{if .}}-in{{split}}{{.}}{{end}}"`                                // -in <inputfile>
	OutFile string `buildarg:"{{if .}}-out{{split}}{{.}}{{end}}"`                               // -out <outputfile>
	Log     Log    `buildarg:"{{if .File}}-log{{if .Append}}a{{end}}{{split}}{{.File}}{{end}}"` // -log[a] <logfile>
	Quiet   bool   `buildarg:"{{if .}}-quiet{{end}}"`                                           // -quiet

	// Formatting:
	Html          bool `buildarg:"{{if .}}-html{{end}}"`      // -html
	Msf           bool `buildarg:"{{if .}}-msf{{end}}"`       // -msf
	Clustal       bool `buildarg:"{{if .}}-clw{{end}}"`       // -clw
	ClustalStrict bool `buildarg:"{{if .}}-clwstrict{{end}}"` // -clwstrict

	// Common options:
	FindDiagonals bool          `buildarg:"{{if .}}-diags{{end}}"`                        // -diags
	MaxIterations int           `buildarg:"{{if .}}-maxiters{{split}}{{.}}{{end}}"`       // -maxiters <n>
	MaxDuration   time.Duration `buildarg:"{{if .}}-maxhours{{split}}{{hours .}}{{end}}"` // -maxhours <h>

	// Other value options (see MUSCLE user guide):
	// Gleaned from user guide - may not reflect reality.
	AnchorSpacing   int     `buildarg:"{{if .}}-anchorspacing{{split}}{{.}}{{end}}"`   // -anchorspacing <n>
	Center          float64 `buildarg:"{{if .}}-center{{split}}{{.}}{{end}}"`          // -center <f.>
	Cluster1        string  `buildarg:"{{if .}}-cluster1{{split}}{{.}}{{end}}"`        // -cluster1 "upgma|upgma|neighborjoining"
	Cluster2        string  `buildarg:"{{if .}}-cluster2{{split}}{{.}}{{end}}"`        // -cluster2 "upgma|upgma|neighborjoining"
	ClustalOut      string  `buildarg:"{{if .}}-clwout{{split}}{{.}}{{end}}"`          // -clwout <file>
	DiagonalBreak   int     `buildarg:"{{if .}}-diagbreak{{split}}{{.}}{{end}}"`       // -diagbreak <n>
	DiagonalLength  int     `buildarg:"{{if .}}-diaglength{{split}}{{.}}{{end}}"`      // -diaglength <n>
	DiagonalMargin  int     `buildarg:"{{if .}}-diagmargin{{split}}{{.}}{{end}}"`      // -diagmargin <n>
	Distance1       string  `buildarg:"{{if .}}-distance1{{split}}{{.}}{{end}}"`       // -distance1 "kmer6_6|kmer20_3|kmer20_4|kbit20_3|kmer4_6"
	Distance2       string  `buildarg:"{{if .}}-distance2{{split}}{{.}}{{end}}"`       // -distance2 "pctid_kimura|pctid_log"
	FastaOut        string  `buildarg:"{{if .}}-fastaout{{split}}{{.}}{{end}}"`        // -fastaout <file>
	GapOpen         float64 `buildarg:"{{if .}}-gapopen{{split}}{{.}}{{end}}"`         // -gapopen <f.>
	GapExtend       float64 `buildarg:"{{if .}}-gapextend{{split}}{{.}}{{end}}"`       // -gapextend <f.>
	HydroWindow     int     `buildarg:"{{if .}}-hydro{{split}}{{.}}{{end}}"`           // -hydro <n>
	HydroFactor     float64 `buildarg:"{{if .}}-hydrofactor{{split}}{{.}}{{end}}"`     // -hydrofactor <f.>
	In1             string  `buildarg:"{{if .}}-in1{{split}}{{.}}{{end}}"`             // -in1 <file>
	In2             string  `buildarg:"{{if .}}-in2{{split}}{{.}}{{end}}"`             // -in2 <file>
	Matrix          string  `buildarg:"{{if .}}-matrix{{split}}{{.}}{{end}}"`          // -matrix <file>
	MaxTrees        int     `buildarg:"{{if .}}-maxtrees{{split}}{{.}}{{end}}"`        // -maxtrees <n>
	MinBestColScore float64 `buildarg:"{{if .}}-minbestcolscore{{split}}{{.}}{{end}}"` // -minbestcolscore <f.>
	MinSmoothScore  float64 `buildarg:"{{if .}}-minsmoothscore{{split}}{{.}}{{end}}"`  // -minsmoothscore <f.>
	MsaOut          string  `buildarg:"{{if .}}-msaout{{split}}{{.}}{{end}}"`          // -msaout <file>
	ObjectiveScore  string  `buildarg:"{{if .}}-objscore{{split}}{{.}}{{end}}"`        // -objscore "sp|ps|dp|xp|spf|spm"
	PhyInterOut     string  `buildarg:"{{if .}}-phyiout{{split}}{{.}}{{end}}"`         // -phyiout <file>
	PhySequenOut    string  `buildarg:"{{if .}}-physout{{split}}{{.}}{{end}}"`         // -physout <file>
	RefineWindow    int     `buildarg:"{{if .}}-refinewindow{{split}}{{.}}{{end}}"`    // -refinewindow <n>
	Root1           string  `buildarg:"{{if .}}-root1{{split}}{{.}}{{end}}"`           // -root1 "pseudo|midlongestspan|minavgleafdist"
	Root2           string  `buildarg:"{{if .}}-root2{{split}}{{.}}{{end}}"`           // -root2 "pseudo|midlongestspan|minavgleafdist"
	ScoreFile       string  `buildarg:"{{if .}}-scorefile{{split}}{{.}}{{end}}"`       // -scorefile <file>
	SeqType         string  `buildarg:"{{if .}}-seqtype{{split}}{{.}}{{end}}"`         // -seqtype "protein|nucleo|auto"
	SmoothScoreCeil float64 `buildarg:"{{if .}}-smoothscoreceil{{split}}{{.}}{{end}}"` // -smoothscoreceil <f.>
	SmoothWindow    int     `buildarg:"{{if .}}-smoothwindow{{split}}{{.}}{{end}}"`    // -smoothwindow <n>
	SpScore         string  `buildarg:"{{if .}}-spscore{{split}}{{.}}{{end}}"`         // -spscore <file>
	Tree1           string  `buildarg:"{{if .}}-tree1{{split}}{{.}}{{end}}"`           // -tree1 <file>
	Tree2           string  `buildarg:"{{if .}}-tree2{{split}}{{.}}{{end}}"`           // -tree2 <file>
	UseTree         string  `buildarg:"{{if .}}-usetree{{split}}{{.}}{{end}}"`         // -usetree <file>
	Weight1         string  `buildarg:"{{if .}}-weight1{{split}}{{.}}{{end}}"`         // -weight1 "none|henikoff|henikoffpb|gsc|clustalw|threeway"
	Weight2         string  `buildarg:"{{if .}}-weight2{{split}}{{.}}{{end}}"`         // -weight2 "none|henikoff|henikoffpb|gsc|clustalw|threeway"

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
	Verbose        bool `buildarg:"{{if .}}-verbose{{end}}"`   // -verbose
}

func hours(d time.Duration) float64 {
	return d.Hours()
}

func (m Muscle) BuildCommand() (*exec.Cmd, error) {
	cl := external.Must(external.Build(m, template.FuncMap{"hours": hours}))
	return exec.Command(cl[0], cl[1:]...), nil
}
