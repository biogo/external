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

// Package external allows uniform interaction with external tools based on a config struct.
package external

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
	"text/template"
)

// CommandBuilder is an interface that assembles a set of command line arguments, and creates
// an *exec.Cmd that can run the command. The method BuildCommand is responsible for handling 
// set up of redirections and parameter sanity checking if required. 
type CommandBuilder interface {
	BuildCommand() (*exec.Cmd, error)
}

// mprintf applies Sprintf with the provided format to each element of slice. It returns an
// error if slice is not a slice or an array or a pointer to either of these types.
func mprintf(format string, slice interface{}) (f []string, err error) {
	v := reflect.ValueOf(slice)
	if kind := v.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		l := v.Len()
		f = make([]string, l)
		for i := 0; i < l; i++ {
			f[i] = fmt.Sprintf(format, v.Index(i).Interface())
		}
	default:
		return nil, errors.New("not a slice or array type")
	}

	return
}

// quote wraps in quotes an item or each element of a slice. The returned value is either
// a string or a slice of strings. Maps are not handled; Quote will panic if given a map to
// process.
func quote(s interface{}) interface{} {
	rv := reflect.ValueOf(s)
	switch rv.Kind() {
	case reflect.Slice:
		q := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			q[i] = fmt.Sprintf("%q", rv.Index(i).Interface())
		}
		return q
	case reflect.Map:
		panic("external: map quoting not handled")
	default:
		return fmt.Sprintf("%q", s)
	}

	panic("cannot reach")
}

// join calls strings.Join with the parameter order reversed to allow use in a template pipeline.
func join(sep string, a []string) string { return strings.Join(a, sep) }

// splitargs is an alias to Join with sep equal to the split tag.
func splitargs(a []string) string { return strings.Join(a, split()) }

// split includes the split tag, "\x00".
func split() string { return string(0) }

// Build builds a set of command line args from cb, which must be a struct. cb's fields
// are inspected for struct tags "buildarg" key. The value for buildarg tag should be a valid
// text template. Build applies executes the template using the value of the field or each
// element of the value of the field if the field is a slice or an array.
// An argument split tag, "\x00", is used to denote separation of elements of the args array
// within any single parameter specification. Template functions can be provided via funcs.
//
// Four convenience functions are provided:
//  quote
//	Wraps each element of a slice of strings in quotes.
//  mprintf
//	Applies fmt.Sprintf to each element of a slice, given a format string.
//  join
//	Calls strings.Join with parameter order reversed to allow pipelining.
//  args
//	Joins a slice of strings with the split tag.
//  split
//	Includes a split tag in a pipeline.
func Build(cb CommandBuilder, funcs ...template.FuncMap) (args []string, err error) {
	v := reflect.ValueOf(cb)
	if kind := v.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("external: not a struct")
	}
	n := v.NumField()
	t := v.Type()
	b := &bytes.Buffer{}
	for i := 0; i < n; i++ {
		tf := t.Field(i)
		if tf.PkgPath != "" {
			continue
		}
		tag := tf.Tag.Get("buildarg")
		if tag != "" {
			tmpl := template.New(tf.Name)
			tmpl.Funcs(template.FuncMap{
				"join":    join,
				"args":    splitargs,
				"split":   split,
				"quote":   quote,
				"mprintf": mprintf,
			})
			for _, fn := range funcs {
				tmpl.Funcs(fn)
			}

			template.Must(tmpl.Parse(tag))
			err = tmpl.Execute(b, v.Field(i).Interface())
			if err != nil {
				return args, err
			}
			if b.Len() > 0 {
				for _, arg := range strings.Split(b.String(), string(0)) {
					if len(arg) > 0 {
						args = append(args, arg)
					}
				}
			}
			b.Reset()
		}
	}

	return
}
