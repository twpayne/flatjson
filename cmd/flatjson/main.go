// flatjson flattens JSON files and generates diffs of flattened JSON files.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/pflag"

	"github.com/twpayne/flatjson"
)

var (
	context = pflag.Int("context", 3, "context")
	diff    = pflag.Bool("diff", false, "diff")
	prefix  = pflag.String("prefix", "root", "prefix.")
	suffix  = pflag.String("suffix", ";\n", "suffix.")
	reverse = pflag.Bool("reverse", false, "reverse")
)

// mergeValuesFromFile reads flat JSON from the file named filename and merges
// it into root.
func mergeValuesFromFile(d *flatjson.Deepener, root interface{}, filename string) (interface{}, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return d.MergeValues(root, r)
}

// runDiff writes the diff of the flat writes of the two files specified on the
// command line.
func runDiff() error {
	if len(pflag.Args()) != 2 {
		return errors.New("-diff requires exactly two filenames")
	}
	text := make([]string, 0, pflag.NArg())
	for _, arg := range pflag.Args() {
		sb := &strings.Builder{}
		f := flatjson.NewFlattener(sb, flatjson.WithPrefix(*prefix), flatjson.WithSuffix(*suffix))
		if err := writeValuesFromFile(f, arg); err != nil {
			return err
		}
		text = append(text, sb.String())
	}
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(text[0]),
		B:        difflib.SplitLines(text[1]),
		FromFile: pflag.Arg(0),
		ToFile:   pflag.Arg(1),
		Context:  *context,
	}
	return difflib.WriteUnifiedDiff(os.Stdout, diff)
}

// runForward flat writes the JSON in each file specified on the command line.
// If no files are specified then the JSON is read from stdin.
func runForward() error {
	f := flatjson.NewFlattener(os.Stdout, flatjson.WithPrefix(*prefix), flatjson.WithSuffix(*suffix))
	if len(pflag.Args()) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		return f.WriteValues(data)
	}
	for _, arg := range pflag.Args() {
		if err := writeValuesFromFile(f, arg); err != nil {
			return err
		}
	}
	return nil
}

// runReverse reads flat JSON from each file specified on the command line and
// writes the resulting JSON to stdout. If no files are specified then the flat
// JSON is read from stdin.
func runReverse() error {
	d := flatjson.NewDeepener()
	var root interface{}
	if len(pflag.Args()) == 0 {
		var err error
		root, err = d.MergeValues(root, os.Stdin)
		if err != nil {
			return err
		}
	} else {
		for _, arg := range pflag.Args() {
			var err error
			root, err = mergeValuesFromFile(d, root, arg)
			if err != nil {
				return err
			}
		}
	}
	return json.NewEncoder(os.Stdout).Encode(root)
}

// writeValuesFromFile reads JSON from the file named filename and writes the
// flattened values to f.
func writeValuesFromFile(f *flatjson.Flattener, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return f.WriteValues(data)
}

func run() error {
	if *diff && *reverse {
		return errors.New("cannot use --diff with --reverse")
	}
	switch {
	case *diff:
		return runDiff()
	case *reverse:
		return runReverse()
	default:
		return runForward()
	}
}

func main() {
	pflag.Parse()
	if err := run(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
