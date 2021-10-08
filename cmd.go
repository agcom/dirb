package main

import (
	"encoding/json"
	"fmt"
	"github.com/agcom/dirb/jsn"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var dir *jsn.Dir
var pretty = false

// Usage: dirb command
func cmd() {
	if len(remArgs) == 0 {
		cmdNon()
	} else {
		pArg0 := remArgs[0]
		remArgs = remArgs[1:]
		switch pArg0 {
		case "init":
			cmdInit()
		case "create", "new", "add":
			cmdNew()
		case "get", "read":
			cmdGet()
		case "update", "up", "patch", "pch":
			cmdUp()
		case "overwrite", "ow", "replace":
			cmdOver()
		case "remove", "rm":
			cmdRm()
		case "help":
			cmdHelp()
		case "grep", "search":
			cmdGrep()
		case "ls", "list":
			cmdLs()
		case "find":
			cmdFind()
		default:
			cmdUnk(pArg0)
		}
	}
}

// Usage: dirb init [-d path]
func cmdInit() {
	if !checkInit() {
		os.Exit(2)
	}

	d := dir.BinDir().Dir()
	err := os.MkdirAll(d, 0775)
	if err != nil {
		errorf("failed to create directory %q (or one of its parents); %v", d, err)
	}
}

func checkInit() bool {
	fail := false

	// Check args
	if len(remArgs) != 0 {
		fail = true
		errorf("unexpected argument(s): %v", remArgs)
	}

	// Check flags

	d := "."
	foundD := false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if foundD {
				// Already found
				fail = true
				errorr("multiple directory flags")
			} else {
				foundD = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)

	return !fail
}

// Usage: dirb create json [-d path]
func cmdNew() {
	if !checkNew() {
		os.Exit(2)
	}

	s := remArgs[0]

	var jo map[string]interface{}
	var err error
	if s == "-" {
		// Read from stdin
		jo, err = jsn.ReaderToJsnObj(os.Stdin)
	} else {
		jo, err = jsn.StrToJsnObj(s)
	}
	if err != nil {
		fatalMultiErr(err)
	}

	name, err := newJsnGenName(dir, jo)
	if err != nil {
		// This command should never fail; unexpected error.
		fatalMultiErr(err)
	} else {
		fmt.Println(name)
	}
}

func checkNew() bool {
	fail := false

	// Check args
	if len(remArgs) != 1 {
		fail = true
		if len(remArgs) == 0 {
			errorr("no argument")
		} else {
			errorf("unexpected argument(s): %v", remArgs[1:])
		}
	}

	// Check flags

	d := "."
	foundD := false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if foundD {
				// Already found
				fail = true
				errorr("multiple directory flags")
			} else {
				foundD = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)

	return !fail
}

// Usage: dirb get name [-d path] [-p [bool]]
func cmdGet() {
	if !checkGet() {
		os.Exit(2)
	}

	name := remArgs[0]

	jo, err := dir.GetObj(name)
	if err != nil {
		fatalMultiErr(err)
	}

	s, err := jsnObjToStrTabIndent(jo, pretty)
	if err != nil {
		fatalMultiErr(err)
	} else {
		if s[len(s)-1] == '\n' {
			fmt.Print(s)
		} else {
			fmt.Println(s)
		}
	}
}

func checkGet() bool {
	fail := false

	// Check args
	if len(remArgs) != 1 {
		fail = true
		if len(remArgs) == 0 {
			errorr("no argument")
		} else {
			errorf("unexpected argument(s): %v", remArgs[1:])
		}
	}

	// Check flags

	d := "."
	foundD := false

	pl := false
	foundP := false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if foundD {
				// Already found
				fail = true
				errorr("multiple \"directory\" flags")
			} else {
				foundD = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		case "p", "pretty":
			if foundP {
				// Already found
				fail = true
				errorr("multiple \"pretty\" flags")
			} else {
				foundP = true
				rmFlag(i)
				if f.HasVal {
					ps := f.Val
					// 1 | 0 | t | f | T | F | true | false | TRUE | FALSE | True | False
					switch ps {
					case "1", "t", "T", "true", "TRUE", "True":
						pl = true
					case "0", "f", "F", "false", "FALSE", "False":
						pl = false
					default:
						fail = true
						errorf("invalid boolean value %q", ps)
					}
				} else {
					pl = true
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)
	pretty = pl

	return !fail
}

// Usage: dirb up name json [-d path]
func cmdUp() {
	if !checkUp() {
		os.Exit(2)
	}

	name := remArgs[0]
	s := remArgs[1]

	var jo map[string]interface{}
	var err error
	if s == "-" {
		// Read from stdin
		jo, err = jsn.ReaderToJsnObj(os.Stdin)
	} else {
		jo, err = jsn.StrToJsnObj(s)
	}
	if err != nil {
		fatalMultiErr(err)
	}

	err = dir.Up(name, jo)
	if err != nil {
		fatalMultiErr(err)
	}
}

func checkUp() bool {
	fail := false

	// Check args
	if len(remArgs) != 2 {
		fail = true
		if len(remArgs) < 2 {
			if len(remArgs) == 0 {
				errorr("no argument")
			} else {
				errorf("not enough arguments (only %d): %v", len(remArgs), remArgs)
			}
		} else {
			errorf("unexpected argument(s): %v", remArgs[2:])
		}
	}

	// Check flags

	d := "."
	foundD := false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if foundD {
				// Already found
				fail = true
				errorr("multiple \"directory\" flags")
			} else {
				foundD = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)

	return !fail
}

// Usage: dirb over name json [-d path]
func cmdOver() {
	if !checkOver() {
		os.Exit(2)
	}

	name := remArgs[0]
	s := remArgs[1]

	var jo map[string]interface{}
	var err error
	if s == "-" {
		// Read from stdin
		jo, err = jsn.ReaderToJsnObj(os.Stdin)
	} else {
		jo, err = jsn.StrToJsnObj(s)
	}
	if err != nil {
		fatalMultiErr(err)
	}

	err = dir.Over(name, jo)
	if err != nil {
		fatalMultiErr(err)
	}
}

func checkOver() bool {
	fail := false

	// Check args
	if len(remArgs) != 2 {
		fail = true
		if len(remArgs) < 2 {
			if len(remArgs) == 0 {
				errorr("no argument")
			} else {
				errorf("not enough arguments (only %d): %v", len(remArgs), remArgs)
			}
		} else {
			errorf("unexpected argument(s): %v", remArgs[2:])
		}
	}

	// Check flags

	d := "."
	foundD := false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if foundD {
				// Already found
				fail = true
				errorr("multiple \"directory\" flags")
			} else {
				foundD = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)

	return !fail
}

// Usage: dirb rm name [-d path]
func cmdRm() {
	if !checkRm() {
		os.Exit(2)
	}

	name := remArgs[0]

	err := dir.Rm(name)
	if err != nil {
		fatalMultiErr(err)
	}
}

func checkRm() bool {
	fail := false

	// Check args
	if len(remArgs) != 1 {
		fail = true
		if len(remArgs) == 0 {
			errorr("no argument")
		} else {
			errorf("unexpected argument(s): %v", remArgs[1:])
		}
	}

	// Check flags

	d := "."
	foundD := false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if foundD {
				// Already found
				fail = true
				errorr("multiple \"directory\" flags")
			} else {
				foundD = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)

	return !fail
}

// Usage: dirb
func cmdNon() {
	fatalc(2, "no command given")
}

// Usage: dirb help command
func cmdHelp() {
	fatal("not yet implemented")
}

// Usage: dirb unknown-command
func cmdUnk(unkCmd string) {
	fatalfc(2, "unknown command %q", unkCmd)
}

// Usage: dirb find l op r [-l [bool]] [-r [bool]] [-d path]
func cmdFind() {
	if !checkFind() {
		os.Exit(2)
	}

	lops := remArgs[0]
	ops := remArgs[1]
	rops := remArgs[2]

	remArgs = remArgs[3:]

	op, err := opFunc(ops)
	if err != nil {
		fatalc(2, err)
	}

	var lop, rop interface{}

	if leftOperandIsFieldRef {
		lop = parseFieldRef(lops)
	} else {
		lop, err = jsn.StrToJsn(lops)
		if err != nil {
			lop = lops
		}
	}

	if rightOperandIsFieldRef {
		rop = parseFieldRef(rops)
	} else {
		rop, err = jsn.StrToJsn(rops)
		if err != nil {
			rop = rops
		}
	}

	find(lop, rop, op)
}

type jsnObjName struct {
	name string
	jo   map[string]interface{}
}

func find(lop, rop interface{}, op func(interface{}, interface{}) bool) {
	fail := false
	ns, err := dir.All()
	if err != nil {
		fail = true
		multiErr(err)
	}

	jons := make([]*jsnObjName, 0, len(ns))
	for _, n := range ns {
		jo, err := dir.GetObj(n)
		if err != nil {
			errorr(err)
		} else {
			jons = append(jons, &jsnObjName{n, jo})
		}
	}

loop:
	for _, jon := range jons {
		n := jon.name
		jo := jon.jo
		var l, r interface{}

		if leftOperandIsFieldRef {
			var nest interface{} = jo
			for _, k := range lop.(fieldRef) {
				if k == "root" {
					continue
				} else {
					if nestJo, ok := nest.(map[string]interface{}); ok {
						if nestJoK, ok := nestJo[k]; ok {
							nest = nestJoK
						} else {
							continue loop
						}
					} else {
						continue loop
					}
				}
			}
			l = nest
		} else {
			l = lop
		}

		if rightOperandIsFieldRef {
			var nest interface{} = jo
			for _, k := range rop.(fieldRef) {
				if k == "root" {
					continue
				} else {
					if nestJo, ok := nest.(map[string]interface{}); ok {
						if nestJoK, ok := nestJo[k]; ok {
							nest = nestJoK
						} else {
							continue loop
						}
					} else {
						continue loop
					}
				}
			}
			r = nest
		} else {
			r = rop
		}

		if op(l, r) {
			fmt.Println(n)
		}
	}

	if fail {
		os.Exit(1)
	}
}

func opFunc(op string) (func(interface{}, interface{}) bool, error) {
	switch op {
	case "<":
		return func(jl interface{}, jr interface{}) bool {
			// Both jl & jr should be of a primitive json type.
			_, jlook := jl.(map[string]interface{})
			_, jlaok := jl.([]interface{})
			_, jrook := jl.(map[string]interface{})
			_, jraok := jl.([]interface{})

			if jlook || jlaok || jrook || jraok {
				return false
			} else {
				return fmt.Sprint(jl) < fmt.Sprint(jr)
			}
		}, nil
	case "<=":
		return func(jl interface{}, jr interface{}) bool {
			// Both jl & jr should be of a primitive json type.
			_, jlook := jl.(map[string]interface{})
			_, jlaok := jl.([]interface{})
			_, jrook := jl.(map[string]interface{})
			_, jraok := jl.([]interface{})

			if jlook || jlaok || jrook || jraok {
				return false
			} else {
				return fmt.Sprint(jl) <= fmt.Sprint(jr)
			}
		}, nil
	case ">":
		return func(jl interface{}, jr interface{}) bool {
			// Both jl & jr should be of a primitive json type.
			_, jlook := jl.(map[string]interface{})
			_, jlaok := jl.([]interface{})
			_, jrook := jl.(map[string]interface{})
			_, jraok := jl.([]interface{})

			if jlook || jlaok || jrook || jraok {
				return false
			} else {
				return fmt.Sprint(jl) > fmt.Sprint(jr)
			}
		}, nil
	case ">=":
		return func(jl interface{}, jr interface{}) bool {
			// Both jl & jr should be of a primitive json type.
			_, jlook := jl.(map[string]interface{})
			_, jlaok := jl.([]interface{})
			_, jrook := jl.(map[string]interface{})
			_, jraok := jl.([]interface{})

			if jlook || jlaok || jrook || jraok {
				return false
			} else {
				return fmt.Sprint(jl) >= fmt.Sprint(jr)
			}
		}, nil
	case "==":
		return func(jl interface{}, jr interface{}) bool {
			return jsnEq(jl, jr)
		}, nil
	case "!=":
		return func(jl interface{}, jr interface{}) bool {
			return !jsnEq(jl, jr)
		}, nil
	case "in":
		return opIn, nil
	case "!in":
		return func(jl interface{}, jr interface{}) bool {
			return !opIn(jl, jr)
		}, nil
	}

	return nil, fmt.Errorf("unknown operator %q", op)
}

func opIn(jl interface{}, jr interface{}) bool {
	if jljo, ok := jl.(map[string]interface{}); ok {
		if jrjo, ok := jr.(map[string]interface{}); ok {
			// Json object in json object
			return objInObj(jljo, jrjo)
		} else if jra, ok := jr.([]interface{}); ok {
			// Json object in array
			return objInArr(jljo, jra)
		} else {
			// Json object in a primitive
			return false
		}
	} else if jla, ok := jl.([]interface{}); ok {
		if jrjo, ok := jr.(map[string]interface{}); ok {
			// Json array in json object
			return arrInObj(jla, jrjo)
		} else if jra, ok := jr.([]interface{}); ok {
			// Array in array
			return arrInArr(jla, jra)
		} else {
			// Array in a primitive
			return false
		}
	} else {
		if jrjo, ok := jr.(map[string]interface{}); ok {
			// A primitive in json object
			return primInObj(jl, jrjo)
		} else if jra, ok := jr.([]interface{}); ok {
			// A primitive in array
			return primInArr(jl, jra)
		} else {
			// A primitive in a primitive
			return primInPrim(jl, jr)
		}
	}
}

func primInPrim(jl interface{}, jr interface{}) bool {
	jls := fmt.Sprint(jl)
	jrs := fmt.Sprint(jr)

	return strings.Contains(jrs, jls)
}

func valInArr(jl interface{}, jar []interface{}) bool {
	for _, jarv := range jar {
		if jsnEq(jl, jarv) {
			return true
		}
	}

	return false
}

func primInArr(jl interface{}, jar []interface{}) bool {
	return valInArr(jl, jar)
}

func primInObj(jl interface{}, jor map[string]interface{}) bool {
	return valInObj(jl, jor)
}

func arrInArr(jal []interface{}, jar []interface{}) bool {
	return valInArr(jal, jar)
}

func arrInObj(jal []interface{}, jor map[string]interface{}) bool {
	return valInObj(jal, jor)
}

func objInArr(jol map[string]interface{}, jar []interface{}) bool {
	return valInArr(jol, jar)
}

func objInObj(jol, jor map[string]interface{}) bool {
	return valInObj(jol, jor)
}

func valInObj(jl interface{}, jor map[string]interface{}) bool {
	for _, vr := range jor {
		if jsnEq(vr, jl) {
			return true
		}
	}

	return false
}

func jsnEq(j1, j2 interface{}) bool {
	if reflect.TypeOf(j1) != reflect.TypeOf(j2) {
		return false
	}

	if j1 == j2 {
		return true
	}

	switch x := j1.(type) {
	case map[string]interface{}:
		y := j2.(map[string]interface{})

		if len(x) != len(y) {
			return false
		}

		for k, v1 := range x {
			v2 := y[k]

			if (v1 == nil) != (v2 == nil) {
				return false
			}

			if !jsnEq(v1, v2) {
				return false
			}
		}

		return true
	case []interface{}:
		y := j2.([]interface{})

		if len(x) != len(y) {
			return false
		}

		var matches int
		flagged := make([]bool, len(y))
		for _, v1 := range x {
			for i, v2 := range y {
				if jsnEq(v1, v2) && !flagged[i] {
					matches++
					flagged[i] = true

					break
				}
			}
		}

		return matches == len(x)
	default:
		return j1 == j2
	}
}

func checkExactRemArgs(i int) error {
	if i < 0 {
		panic(fmt.Sprintf("negative number of args %d", i))
	}

	l := len(remArgs)
	if l != i {
		if l < i {
			if l == 0 {
				return fmt.Errorf("no argument")
			} else {
				return fmt.Errorf("not enough arguments (only %d): %s", l, strings.Join(remArgs, " "))
			}
		} else {
			if l == 1 {
				// i is 0
				return fmt.Errorf("unexpected argument: %v", remArgs[0])
			} else {
				return fmt.Errorf("unexpected arguments: %s", strings.Join(remArgs, " "))
			}
		}
	}

	return nil
}

func parseBoolVal(s string) (bool, error) {
	// 1 | 0 | t | f | T | F | true | false | TRUE | FALSE | True | False
	switch s {
	case "1", "t", "T", "true", "TRUE", "True":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value %q", s)
	}
}

type fieldRef []string

var fieldRefRegex *regexp.Regexp

func parseFieldRef(s string) fieldRef {
	if fieldRefRegex == nil {
		fieldRefRegex = regexp.MustCompile(`([^\\])(\.)`)
	}

	ref := make([]string, 0, 1)
	iss := fieldRefRegex.FindAllStringSubmatchIndex(s, -1)
	if len(iss) == 0 {
		ref = append(ref, s)
	} else {
		for _, is := range iss {
			_ = is[4]
			end := is[5]
			ref = append(ref, s[:end-1])
			s = s[end:]
		}
	}

	return ref
}

var leftOperandIsFieldRef, rightOperandIsFieldRef bool

func checkFind() bool {
	fail := false

	// Check args
	err := checkExactRemArgs(3)
	if err != nil {
		fail = false
		errorr(err)
	}

	// Check flags

	d, df := ".", false
	pp, pf := false, false
	l, lf := false, false
	r, rf := false, false

	for i, f := range flags {
		switch f.Name {
		case "d", "directory":
			if df {
				// Already found
				fail = true
				errorr("multiple \"directory\" flags")
			} else {
				df = true
				rmFlag(i)
				if f.HasVal {
					d = f.Val
				} else {
					fail = true
					errorr("no value assigned to a \"directory\" flag")
				}
			}
		case "p", "pretty":
			if pf {
				// Already found
				fail = true
				errorr("multiple \"pretty\" flags")
			} else {
				pf = true
				rmFlag(i)
				if f.HasVal {
					ps := f.Val
					// 1 | 0 | t | f | T | F | true | false | TRUE | FALSE | True | False
					switch ps {
					case "1", "t", "T", "true", "TRUE", "True":
						pp = true
					case "0", "f", "F", "false", "FALSE", "False":
						pp = false
					default:
						fail = true
						errorf("invalid boolean value %q", ps)
					}
				} else {
					pp = true
				}
			}
		case "l", "left-operand-is-field-reference":
			if lf {
				// Already found
				fail = true
				errorr("multiple \"left-operand-is-field-reference\" flags")
			} else {
				lf = true
				rmFlag(i)
				if f.HasVal {
					ls := f.Val
					var err error
					l, err = parseBoolVal(ls)
					if err != nil {
						fail = true
						errorr(err)
					}
				} else {
					l = true
				}
			}
		case "r", "right-operand-is-field-reference":
			if rf {
				// Already found
				fail = true
				errorr("multiple \"right-operand-is-field-reference\" flags")
			} else {
				rf = true
				rmFlag(i)
				if f.HasVal {
					rs := f.Val
					var err error
					r, err = parseBoolVal(rs)
					if err != nil {
						fail = true
						errorr(err)
					}
				} else {
					r = true
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	dir = jsn.NewDir(d)
	pretty = pp
	leftOperandIsFieldRef = l
	rightOperandIsFieldRef = r

	return !fail
}

// Usage: dirb grep regex [-d path]
func cmdGrep() {
	fatal("not yet implemented")
}

// Usage: dirb ls [-d path]
func cmdLs() {
	fatal("not yet implemented")
}

func rmFlag(i int) {
	flags = append(flags[:i], flags[i+1:]...)
}

func jsnObjToStrTabIndent(jo map[string]interface{}, tabIndent bool) (string, error) {
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	if tabIndent {
		enc.SetIndent("", "\t")
	}
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(jo)
		_ = w.CloseWithError(err)
	}()
	b, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("failed to encode into a json string; %w", err)
	}

	return string(b), nil
}
