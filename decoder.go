package flatjson

import (
	"bytes"
	"fmt"
	"io"
)

// A Decoder decodes flatjson.
type Decoder struct {
	r io.Reader
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

// Decode decodes a value, merging it into root.
func (d *Decoder) Decode(root interface{}) (interface{}, error) {
	assignments, err := newParser(d.r).parseAssignments()
	if err != nil {
		return nil, err
	}
	for _, assignment := range assignments {
		root = recursiveMerge(root, assignment.properties, assignment.value)
	}
	return root, nil
}

func recursiveMerge(root interface{}, properties []interface{}, value interface{}) interface{} {
	if len(properties) == 0 {
		return value
	}
	switch property := properties[0].(type) {
	case int:
		array, ok := root.([]interface{})
		if !ok {
			array = []interface{}{}
		}
		if property < len(array) {
			array[property] = recursiveMerge(array[property], properties[1:], value)
		} else {
			array = append(array, recursiveMerge(nil, properties[1:], value))
		}
		return array
	case string:
		object, ok := root.(map[string]interface{})
		if !ok || object == nil {
			object = make(map[string]interface{})
		}
		object[property] = recursiveMerge(object[property], properties[1:], value)
		return object
	default:
		panic(fmt.Sprintf("unexpected property %v (%T)", property, property))
	}
}

// Unmarshal parses the flatjson-encoded data and stores the result in the
// value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	_, err := NewDecoder(bytes.NewBuffer(data)).Decode(v)
	return err
}
