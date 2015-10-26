// Copyright ©2011-2012 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// mprintf applies Sprintf with the provided format to each element of an array, slice or map, or
// pointer to any of these, otherwise if returns the fmt.Sprintf representation of the underlying
// value with the given format.
func mprintf(format string, value interface{}) interface{} {
	rv := reflect.ValueOf(value)
	if kind := rv.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		q := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			q[i] = fmt.Sprintf(format, rv.Index(i).Interface())
		}
		return q
	case reflect.Map:
		q := make([]string, rv.Len())
		for i, k := range rv.MapKeys() {
			q[i] = fmt.Sprintf(format, rv.MapIndex(k).Interface())
		}
		return q
	default:
		return fmt.Sprintf(format, rv.Interface())
	}

	panic("cannot reach")
}

// quote wraps in quotes an item or each element of an array, slice or map by calling mprintf with
// "%q" as the format.
func quote(value interface{}) interface{} { return mprintf("%q", value) }

// join performs the genric equivalent of a call to strings.Join with the parameter order
// reversed to allow use in a template pipeline.
func join(sep string, a interface{}) string {
	rv := reflect.ValueOf(a)
	if kind := rv.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		cs := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			cs[i] = fmt.Sprint(rv.Index(i))
		}
		return strings.Join(cs, sep)
	case reflect.Map:
		cs := make([]string, rv.Len())
		for i, k := range rv.MapKeys() {
			cs[i] = fmt.Sprint(rv.MapIndex(k))
		}
		return strings.Join(cs, sep)
	}
	return fmt.Sprint(a)
}

// splitargs is an alias to join with sep equal to the split tag.
func splitargs(a interface{}) string { return join(split(), a) }

// split includes the split tag, "\x00".
func split() string { return string(0) }

// Build builds a set of command line args from cb, which must be a struct. cb's fields
// are inspected for struct tags "buildarg" key. The value for buildarg tag should be a valid
// text template. Build applies executes the template using the value of the field or each
// element of the value of the field if the field is a slice or an array.
// An argument split tag, "\x00", can be used to denote separation of elements of the args array
// within any single parameter specification. Template functions can be provided via funcs.
//
// Four convenience functions are provided:
//  args
//	Joins %v representation of elements of an array, slice or map, or reference to any of
//	these, using split tag as a separator. Otherwise it returns the %v representation of the
//	underlying value.
//  join
//	Joins %v representation of elements of an array, slice or map, or reference to any of
//	these, using the the value of the first argument as a separator. Otherwise it returns the
//	%v representation of the underlying value.
//  mprintf
//	Applies fmt.Sprintf, given a format string, to a value or each element of an array, slice
//	or map, or reference to any of these.
//  quote
//	Wraps in quotes a value or each element of an array, slice or map, or reference to any
//	of these.
//  split
//	Includes a split tag in a pipeline.
//
//  Note that args, join, mprintf and quote will return randomly ordered arguments if a map is used
//  as a template input.
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
		if tf.PkgPath != "" && !tf.Anonymous {
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

// Must is a helper that wraps a call to a function returning ([]string, error)
// and panics if the error is non-nil.
func Must(args []string, err error) []string {
	if err != nil {
		panic(fmt.Sprintf("external: failed to build args list: %v: built %#v", err, args))
	}
	return args
}
