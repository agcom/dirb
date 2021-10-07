package sflag

import "fmt"

type Flag struct {
	Name   string
	Val    string
	HasVal bool
}

func Parse(args []string) ([]*Flag, []string, error) {
	return ParseNoNxtArgVal(args, nil)
}

func ParseNoNxtArgVal(args []string, noNxtArgVal []string) ([]*Flag, []string, error) {
	noNxtArgValM := make(map[string]interface{}, len(noNxtArgVal))
	for _, e := range noNxtArgVal {
		noNxtArgValM[e] = nil
	}
	fs := make([]*Flag, len(args)/2)[:0]
	cArgs := make([]string, len(args))
	copy(cArgs, args)
	args = cArgs

	s := session{args: args, i: 0}
	for {
		f, err := s.parseOne(noNxtArgValM)
		if f != nil {
			fs = append(fs, f)
			continue
		} else if err == nil { // No more to parse
			return fs, s.args, nil
		} else {
			return fs, s.args, err
		}
	}
}

type session struct {
	args []string
	i    int
}

func remove(ss []string, i int) []string {
	return append(ss[:i], ss[i+1:]...)
}

func (s *session) parseOne(noNxtArgValM map[string]interface{}) (*Flag, error) {
	var arg string
	var numMinuses int
	for {
		if len(s.args[s.i:]) == 0 {
			return nil, nil
		}

		arg = s.args[s.i:][0]
		if len(arg) < 2 || arg[0] != '-' { // Not a flag arg
			s.i++
			continue
		} else {
			numMinuses = 1
			if arg[1] == '-' {
				numMinuses++
				if len(arg) == 2 { // -- terminates the flags
					s.args = remove(s.args, s.i)
					return nil, nil
				}
			}
			break
		}
	}

	nv := arg[numMinuses:]
	if nv[0] == '-' {
		return nil, fmt.Errorf("invalid flag %q; its name starts with '-'", nv)
	} else if nv[0] == '=' {
		return nil, fmt.Errorf("invalid flag %q; its name starts with '='", nv)
	}

	// Hit; search for its value if any.
	s.args = remove(s.args, s.i)
	hasVal := false
	name := nv
	val := ""
	for i := 1; i < len(nv); i++ {
		if nv[i] == '=' {
			hasVal = true
			val = nv[i+1:]
			name = nv[:i]
			break
		}
	}

	if _, ok := noNxtArgValM[name]; !hasVal && !ok {
		// Read the next arg as its value
		if len(s.args[s.i:]) > 0 {
			hasVal = true
			val = s.args[s.i]
			s.args = remove(s.args, s.i)
		}
	}

	return &Flag{name, val, hasVal}, nil
}
