// Copyright ©2011-2012 Dan Kortschak <dan.kortschak@adelaide.edu.au>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package external

import (
	"fmt"
)

func ExampleBuild_1() {
	// samtools sort [-n] [-m maxMem] <in.bam> <out.prefix>
	type SamToolsSort struct {
		Name      string
		Comment   string
		Cmd       string   `buildarg:"{{if .}}{{.}}{{else}}samtools{{end}}"` // samtools
		SubCmd    struct{} `buildarg:"sort"`                                 // sort
		SortNames bool     `buildarg:"{{if .}}-n{{end}}"`                    // [-n]
		MaxMem    int      `buildarg:"{{if .}}-m{{split}}{{.}}{{end}}"`      // [-m maxMem]
		InFile    string   `buildarg:"{{.}}"`                                // "<in.bam>"
		OutFile   string   `buildarg:"{{.}}"`                                // "<out.prefix>"
		CommandBuilder
	}

	s := SamToolsSort{
		Name:      "Sort",
		SortNames: true,
		MaxMem:    1e8,
		InFile:    "infile",
		OutFile:   "outfile",
	}

	args, err := Build(s)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", args)
	}
	// Output:
	// []string{"samtools", "sort", "-n", "-m", "100000000", "infile", "outfile"}
}

func ExampleBuild_2() {
	// samtools merge [-h inh.sam] [-n] <out.bam> <in1.bam> <in2.bam> [...]
	type SamToolsMerge struct {
		Name       string
		Comment    string
		Cmd        string   `buildarg:"{{if .}}{{.}}{{else}}samtools{{end}}"` // samtools
		SubCmd     struct{} `buildarg:"merge"`                                // merge
		HeaderFile string   `buildarg:"{{if .}}-h{{split}}{{.}}{{end}}"`      // [-h inh.sam]
		SortNames  bool     `buildarg:"{{if .}}-n{{end}}"`                    // [-n]
		OutFile    string   `buildarg:"{{.}}"`                                // <out.bam>
		InFiles    []string `buildarg:"{{args .}}"`                           // <in.bam>...
		CommandBuilder
	}

	s := &SamToolsMerge{
		Name:       "Merge",
		Cmd:        "samtools",
		HeaderFile: "header",
		InFiles:    []string{"infile1", "infile2"},
		OutFile:    "outfile",
	}

	args, err := Build(s)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", args)
	}
	// Output:
	// []string{"samtools", "merge", "-h", "header", "outfile", "infile1", "infile2"}
}

func ExampleBuild_3() {
	// sed [-n] [-e <exp>]... [-f <file>]... [--follow-symlinks] [-i[suf]] [-l <len>] [--posix] [-r] [-s] [-s] <in>... > <out>
	type InPlace struct {
		Yes bool
		Suf string
	}
	type Sed struct {
		Name       string
		Comment    string
		Cmd        string   `buildarg:"{{if .}}{{.}}{{else}}sed{{end}}"`                    // sed
		Quiet      bool     `buildarg:"{{if .}}-n{{end}}"`                                  // [-n]
		Script     []string `buildarg:"{{if .}}{{mprintf \"-e\x00'%s'\" . | args}}{{end}}"` // [-e '<exp>']...
		ScriptFile []string `buildarg:"{{if .}}{{mprintf \"-f\x00%s\" . | args}}{{end}}"`   // [-f "<file>"]...
		Follow     bool     `buildarg:"{{if .}}--follow-symlinks{{end}}"`                   // [--follow-symlinks]
		InPlace    InPlace  `buildarg:"{{if .Yes}}-i{{with .Suf}}{{.}}{{end}}{{end}}"`      // [-i[suf]]
		WrapAt     int      `buildarg:"{{if .}}-l{{split}}{{.}}{{end}}"`                    // [-l <len>]
		Posix      bool     `buildarg:"{{if .}}--posix{{end}}"`                             // [--posix]
		ExtendRE   bool     `buildarg:"{{if .}}-r{{end}}"`                                  // [-r]
		Separate   bool     `buildarg:"{{if .}}-s{{end}}"`                                  // [-s]
		Unbuffered bool     `buildarg:"{{if .}}-u{{end}}"`                                  // [-u]
		InFiles    []string `buildarg:"{{args .}}"`                                         // "<in>"...
		CommandBuilder
	}

	s := &Sed{
		Name:    "Sed",
		Cmd:     "sed",
		WrapAt:  76,
		Script:  []string{`s/\<hi\>/lo/g`, `s/\<left\>/right/g`},
		InPlace: InPlace{true, "bottomright"},
		InFiles: []string{"infile1", "infile2"},
	}

	args, err := Build(s)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", args)
	}
	// Output:
	// []string{"sed", "-e", "'s/\\<hi\\>/lo/g'", "-e", "'s/\\<left\\>/right/g'", "-ibottomright", "-l", "76", "infile1", "infile2"}
}

func ExampleBuild_4() {
	// bowtie [options]* <ebwt> {-1 <m1> -2 <m2> | --12 <r> | <s>} [<hit>]
	type Bowtie struct {
		Name     string
		Comment  string
		Cmd      string   `buildarg:"{{if .}}{{.}}{{else}}bowtie{{end}}"`  // bowtie
		Index    string   `buildarg:"{{.}}"`                               // <ebwt>
		One      []string `buildarg:"{{if .}}-1{{join \",\" .}}{{end}}"`   // -1 <m1>
		Two      []string `buildarg:"{{if .}}-2{{join \",\" .}}{{end}}"`   // -2 <m2>
		Mixed    []string `buildarg:"{{if .}}--12{{join \",\" .}}{{end}}"` // --12 <r>
		Unpaired []string `buildarg:"{{if .}}{{join \",\" .}}{{end}}"`     // <s>
		OutFile  string   `buildarg:"{{if .}}{{.}}{{end}}"`                // <hit>
		CommandBuilder
	}

	b := &Bowtie{
		Name:     "Bowtie",
		Cmd:      "bowtie",
		Index:    "ebwt",
		Unpaired: []string{"a.fa", "b.fa", "c.fa", "d.fa", "e.fa"},
		OutFile:  "oufile",
	}

	args, err := Build(b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", args)
	}
	// Output:
	// []string{"bowtie", "ebwt", "a.fa,b.fa,c.fa,d.fa,e.fa", "oufile"}
}
