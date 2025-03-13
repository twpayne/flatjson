package flatjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
)

var (
	unquotedPropertyNameRegexp = regexp.MustCompile(`^[A-Za-z_][0-9A-Za-z_]*$`)

	keywords = map[string]bool{
		"false": true,
		"null":  true,
		"true":  true,
	}
)

// A Flattener converts JSON into flat JSON.
type Flattener struct {
	w      io.Writer
	prefix string
	suffix string
}

// A FlattenerOption sets an option on a Flattener.
type FlattenerOption func(*Flattener)

func propertyAccessor(path, name string) string {
	if unquotedPropertyNameRegexp.MatchString(name) && !keywords[name] {
		if path == "" {
			return name
		}
		return path + "." + name
	}
	return fmt.Sprintf("%s[%q]", path, name)
}

// NewFlattener returns a new Flattener that writes to w.
func NewFlattener(w io.Writer, options ...FlattenerOption) *Flattener {
	f := &Flattener{
		w:      w,
		prefix: "root",
		suffix: ";\n",
	}
	for _, option := range options {
		option(f)
	}
	return f
}

func (f *Flattener) writeArrayValues(path string, array []interface{}) error {
	if _, err := fmt.Fprintf(f.w, "%s = []%s", path, f.suffix); err != nil {
		return err
	}
	for i, value := range array {
		if err := f.writeValuesHelper(path+"["+strconv.Itoa(i)+"]", value); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flattener) writeObjectValues(path string, object map[string]interface{}) error {
	if _, err := fmt.Fprintf(f.w, "%s = {}%s", path, f.suffix); err != nil {
		return err
	}
	properties := make([]string, 0, len(object))
	for property := range object {
		properties = append(properties, property)
	}
	sort.Strings(properties)
	for _, property := range properties {
		if err := f.writeValuesHelper(propertyAccessor(path, property), object[property]); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flattener) writeValuesHelper(path string, value interface{}) error {
	switch value := value.(type) {
	case []interface{}:
		return f.writeArrayValues(path, value)
	case map[string]interface{}:
		return f.writeObjectValues(path, value)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		_, err = f.w.Write([]byte(path + " = " + string(data) + f.suffix))
		return err
	}
}

// WriteValues decodes JSON from data and writes it.
func (f *Flattener) WriteValues(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	var value interface{}
	if err := d.Decode(&value); err != nil {
		return err
	}
	return f.writeValuesHelper(f.prefix, value)
}

// WithPrefix sets the prefix on a Flattener.
func WithPrefix(prefix string) FlattenerOption {
	return func(f *Flattener) {
		f.prefix = prefix
	}
}

// WithSuffix sets the suffix on a Flattener.
func WithSuffix(suffix string) FlattenerOption {
	return func(f *Flattener) {
		f.suffix = suffix
	}
}
