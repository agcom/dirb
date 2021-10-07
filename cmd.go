package main

import (
	"fmt"
	"github.com/agcom/bs/jsn"
	"os"
)

var rootDir jsn.Dir
var pretty = false

// Usage: bs <command>
func cmd() {
	if len(args) == 0 {
		cmdNon()
	} else {
		pArg0 := pArgs[0]
		pArgs = pArgs[1:]
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
		default:
			cmdUnk(pArg0)
		}
	}
}

// Usage: bs init [-d <path>]
func cmdInit() {
	if !checkInit() {
		os.Exit(2)
	}

	regDefEnts()

	for _, entDir := range entDirs {
		d := entDir.BinDir().Dir()
		err := os.MkdirAll(d, 0775)
		if err != nil {
			errorf("failed to make directory %q (or one of its parent directories); %v", d, err)
		}
	}
}

func checkInit() bool {
	fail := false

	// Check positional args
	if len(pArgs) != 0 {
		fail = true
		errorf("unexpected positional argument(s): %v", pArgs)
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

	rootDir = *jsn.NewDir(d)

	return !fail
}

// Usage: bs create <entity> <json> [-d <path>]
func cmdNew() {
	if !checkNew() {
		os.Exit(2)
	}

	regDefEnts()

	ent := pArgs[0]
	s := pArgs[1]

	entDir := entDir(ent)
	if entDir == nil {
		fatalf("entity %q isn't defined", ent)
	}

	var jo map[string]interface{}
	var err error
	if s == "-" {
		// Read from stdin
		jo, err = readerToJsnObj(os.Stdin)
	} else {
		jo, err = strToJsnObj(s)
	}
	if err != nil {
		fatalErr(err)
	}

	name, err := newJsnGenName(entDir, jo)
	if err != nil {
		// This command should never fail; unexpected error.
		fatalErr(err)
	} else {
		fmt.Println(name)
	}
}

func checkNew() bool {
	fail := false

	// Check positional args
	if len(pArgs) != 2 {
		fail = true
		if len(pArgs) < 2 {
			if len(pArgs) == 0 {
				errorr("no positional argument")
			} else if len(pArgs) == 1 {
				errorf("not enough positional arguments (only 1): %v", pArgs)
			}
		} else {
			errorf("unexpected positional argument(s): %v", pArgs[2:])
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

	rootDir = *jsn.NewDir(d)

	return !fail
}

// Usage: bs get <entity> <name> [-d <path>] [-p [<bool>]]
func cmdGet() {
	if !checkGet() {
		os.Exit(2)
	}

	regDefEnts()

	ent := pArgs[0]
	name := pArgs[1]

	entDir := entDir(ent)
	if entDir == nil {
		fatalf("entity %q isn't defined", ent)
	}

	jo, err := entDir.GetObj(name)
	if err != nil {
		fatalErr(err)
	}

	s, err := jsnObjToStr(jo, pretty)
	if err != nil {
		fatalErr(err)
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

	// Check positional args
	if len(pArgs) != 2 {
		fail = true
		if len(pArgs) < 2 {
			if len(pArgs) == 0 {
				errorr("no positional argument")
			} else if len(pArgs) == 1 {
				errorf("not enough positional arguments (only 1): %v", pArgs)
			}
		} else {
			errorf("unexpected positional argument(s): %v", pArgs[2:])
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
					fail = true
					errorr("no value assigned to \"directory\" flag")
				}
			}
		default:
			fail = true
			errorf("unexpected flag %q", f.Name)
		}
	}

	rootDir = *jsn.NewDir(d)
	pretty = pl

	return !fail
}

// Usage: bs up <entity> <name> <json> [-d <path>]
func cmdUp() {
	if !checkUp() {
		os.Exit(2)
	}

	regDefEnts()

	ent := pArgs[0]
	name := pArgs[1]
	s := pArgs[2]

	entDir := entDir(ent)
	if entDir == nil {
		fatalf("entity %q isn't defined", ent)
	}

	var jo map[string]interface{}
	var err error
	if s == "-" {
		// Read from stdin
		jo, err = readerToJsnObj(os.Stdin)
	} else {
		jo, err = strToJsnObj(s)
	}
	if err != nil {
		fatalErr(err)
	}

	err = entDir.Up(name, jo)
	if err != nil {
		fatalErr(err)
	}
}

func checkUp() bool {
	fail := false

	// Check positional args
	if len(pArgs) != 3 {
		fail = true
		if len(pArgs) < 3 {
			if len(pArgs) == 0 {
				errorr("no positional argument")
			} else {
				errorf("not enough positional arguments (only %d): %v", len(pArgs), pArgs)
			}
		} else {
			errorf("unexpected positional argument(s): %v", pArgs[2:])
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

	rootDir = *jsn.NewDir(d)

	return !fail
}

// Usage: bs over <entity> <name> <json> [-d <path>]
func cmdOver() {
	if !checkOver() {
		os.Exit(2)
	}

	regDefEnts()

	ent := pArgs[0]
	name := pArgs[1]
	s := pArgs[2]

	entDir := entDir(ent)
	if entDir == nil {
		fatalf("entity %q isn't defined", ent)
	}

	var jo map[string]interface{}
	var err error
	if s == "-" {
		// Read from stdin
		jo, err = readerToJsnObj(os.Stdin)
	} else {
		jo, err = strToJsnObj(s)
	}
	if err != nil {
		fatalErr(err)
	}

	err = entDir.Over(name, jo)
	if err != nil {
		fatalErr(err)
	}
}

func checkOver() bool {
	fail := false

	// Check positional args
	if len(pArgs) != 3 {
		fail = true
		if len(pArgs) < 3 {
			if len(pArgs) == 0 {
				errorr("no positional argument")
			} else {
				errorf("not enough positional arguments (only %d): %v", len(pArgs), pArgs)
			}
		} else {
			errorf("unexpected positional argument(s): %v", pArgs[2:])
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

	rootDir = *jsn.NewDir(d)

	return !fail
}

// Usage: bs rm <entity> <name> [-d <path>]
func cmdRm() {
	if !checkRm() {
		os.Exit(2)
	}

	regDefEnts()

	ent := pArgs[0]
	name := pArgs[1]

	entDir := entDir(ent)
	if entDir == nil {
		fatalf("entity %q isn't defined", ent)
	}

	err := entDir.Rm(name)
	if err != nil {
		fatalErr(err)
	}
}

func checkRm() bool {
	fail := false

	// Check positional args
	if len(pArgs) != 2 {
		fail = true
		if len(pArgs) < 2 {
			if len(pArgs) == 0 {
				errorr("no positional argument")
			} else if len(pArgs) == 1 {
				errorf("not enough positional arguments (only 1): %v", pArgs)
			}
		} else {
			errorf("unexpected positional argument(s): %v", pArgs[2:])
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

	rootDir = *jsn.NewDir(d)

	return !fail
}

// Usage: bs
func cmdNon() {
	fatalCode(2, "no command given")
}

// Usage: bs help [<command>]
func cmdHelp() {
	fatal("not yet implemented")
}

// Usage: bs <unknown-command>
func cmdUnk(cmd string) {
	fatalfCode(2, "unknown command %q", cmd)
}

func rmFlag(i int) {
	flags = append(flags[:i], flags[i+1:]...)
}
