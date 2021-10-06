package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/agcom/bs/bins"
	"github.com/agcom/bs/jsns"
	"io"
	"os"
	"path/filepath"
)

type repo struct {
	jsnObjs *jsns.ObjJsns
	jsns    *jsns.Bins
	bins    *bins.Dir
}

var rootDir = "."
var pretty = true

var repoLkp = make(map[string]*repo, 3)
var repos = make([]*repo, 0, 3)

var exe = filepath.Base(os.Args[0])

func newEntity(name string, aliases []string) {
	bs := bins.NewDir(filepath.Join(rootDir, name))
	js := jsns.NewBins(bs)
	jos := jsns.NewObjJsns(jsns.JsnExtMid(js))

	repo := &repo{
		jsnObjs: jos,
		jsns:    js,
		bins:    bs,
	}

	repos = append(repos, repo)

	repoLkp[name] = repo

	for _, alias := range aliases {
		repoLkp[alias] = repo
	}
}

const bk = "bk"
const slr = "slr"
const byr = "byr"

func init() {
	newEntity("books", []string{"book", bk, "bks"})
	newEntity("sellers", []string{"seller", slr, "slrs"})
	newEntity("buyers", []string{"buyer", byr, "byrs"})
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		nonCmd()
	} else {
		switch args[0] {
		case "init":
			initCmd(args[1:])
		case "create", "new", "add":
			newCmd(args[1:])
		case "get", "read":
			getCmd(args[1:])
		case "help":
			helpCmd(args[1:])
		case "overwrite", "ow", "replace":
			owCmd(args[1:])
		case "remove", "rm":
			rmCmd(args[1:])
		case "update", "up", "patch", "pch":
			upCmd(args[1:])
		default:
			unkCmd(args[0], args[1:])
		}
	}
}

func initCmd(args []string) {
	err := checkNoArgs(args)
	if err != nil {
		Fatal(err)
	}

	fail := false
	for _, repo := range repos {
		dir := repo.bins.Dir()
		err := os.MkdirAll(dir, 0770)
		if err != nil {
			fail = true
			Errorf("failed to make directory \"%s\" (or one of its parent directories) ", dir)
		}
	}

	if fail {
		os.Exit(1)
	}
}

func newCmd(args []string) {
	err := checkExactArgs(args, 2)
	if err != nil {
		Fatal(err)
	}

	entityName := args[0]
	s := args[1]

	r, ok := repoLkp[entityName]
	if !ok {
		FatalfCode(2, "entity \"%s\" isn't defined", entityName)
	}

	var jo jsns.JsnObj
	if s == "-" {
		jo, err = readerToJsnObj(os.Stdin)
	} else {
		jo, err = strToJsnObj(s)
	}
	if err != nil {
		Fatal(err)
	}

	name, err := jsns.NewJsnObjGenName(r.jsnObjs, jo)
	if err != nil {
		// TODO: translate errors.
		Fatal(err)
	}

	fmt.Println(name)
}

func unkCmd(cmd string, args []string) {
	Fatalf("unknown command \"%s\"", cmd)
}

func nonCmd() {
	FatalfCode(2, "no command given; use \"%s help\" for more info.", exe)
}

func getCmd(args []string) {
	err := checkExactArgs(args, 2)
	if err != nil {
		Fatal(err)
	}

	entityName := args[0]
	name := args[1]

	r, ok := repoLkp[entityName]
	if !ok {
		FatalfCode(2, "entity \"%s\" isn't defined", entityName)
	}

	jo, err := r.jsnObjs.Get(name)
	if err != nil {
		// TODO: translate errors.
		Fatal(err)
	}

	str, err := jsnObjToStr(jo, pretty)
	if err != nil {
		Fatal(err)
	}

	if str[len(str)-1] == '\n' {
		fmt.Print(str)
	} else {
		fmt.Println(str)
	}
}

func helpCmd(args []string) {

}

func owCmd(args []string) {
	err := checkExactArgs(args, 3)
	if err != nil {
		Fatal(err)
	}

	entityName := args[0]
	name := args[1]
	s := args[2]

	r, ok := repoLkp[entityName]
	if !ok {
		FatalfCode(2, "entity \"%s\" isn't defined", entityName)
	}

	var jo jsns.JsnObj
	if s == "-" {
		jo, err = readerToJsnObj(os.Stdin)
	} else {
		jo, err = strToJsnObj(s)
	}
	if err != nil {
		Fatal(err)
	}

	err = r.jsnObjs.Ow(name, jo)
	if err != nil {
		// TODO: translate errors.
		Fatal(err)
	}
}

func rmCmd(args []string) {
	err := checkExactArgs(args, 2)
	if err != nil {
		Fatal(err)
	}

	entityName := args[0]
	name := args[1]

	r, ok := repoLkp[entityName]
	if !ok {
		FatalfCode(2, "entity \"%s\" isn't defined", entityName)
	}

	err = r.jsnObjs.Rm(name)
	if err != nil {
		// TODO: translate errors.
		Fatal(err)
	}
}

func upCmd(args []string) {
	err := checkExactArgs(args, 3)
	if err != nil {
		Fatal(err)
	}

	entityName := args[0]
	name := args[1]
	s := args[2]

	r, ok := repoLkp[entityName]
	if !ok {
		FatalfCode(2, "entity \"%s\" isn't defined", entityName)
	}

	var jo jsns.JsnObj
	if s == "-" {
		jo, err = readerToJsnObj(os.Stdin)
	} else {
		jo, err = strToJsnObj(s)
	}
	if err != nil {
		Fatal(err)
	}

	err = jsns.WeakUpJsnObj(r.jsnObjs, name, jo)
	if err != nil {
		// TODO: translate error.
		Fatal(err)
	}
}

func checkExactArgs(args []string, l int) error {
	if l == 0 {
		return checkNoArgs(args)
	} else if l < 0 {
		panic("negative number of arguments")
	} else {
		if len(args) != l {
			var s string
			if l == 1 {
				s = "arg"
			} else {
				s = "args"
			}
			return fmt.Errorf("accepts %d %s, but received %d", l, s, len(args))
		} else {
			return nil
		}
	}
}

func checkNoArgs(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("accepts no args, but received %d", len(args))
	}

	return nil
}

var jsnButNotJsnObj = errors.New("the json string is not a json object string")

func strToJsnObj(s string) (jsns.JsnObj, error) {
	return bytesToJsnObj([]byte(s))
}

func bytesToJsnObj(b []byte) (jsns.JsnObj, error) {
	var j jsns.Jsn
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, fmt.Errorf("invalid json string %v; %w", string(b), err)
	}
	if jo, ok := j.(jsns.JsnObj); ok {
		return jo, nil
	} else {
		return nil, jsnButNotJsnObj
	}
}

func jsnObjToStr(j jsns.JsnObj, indent bool) (string, error) {
	r, w := io.Pipe()
	enc := json.NewEncoder(w)
	if indent {
		enc.SetIndent("", "\t")
	}
	enc.SetEscapeHTML(false)
	go func() {
		err := enc.Encode(j)
		_ = w.CloseWithError(err)
	}()
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func readerToJsnObj(r io.Reader) (jsns.JsnObj, error) {
	dec := json.NewDecoder(r)
	var j jsns.Jsn
	err := dec.Decode(&j)
	if err != nil {
		return nil, err
	}
	if jo, ok := j.(jsns.JsnObj); ok {
		return jo, nil
	} else {
		return nil, jsnButNotJsnObj
	}
}
