package flatjson

import (
	"fmt"
	"io"
)

// A Deepener converts flat JSON into a JSON object.
type Deepener struct{}

// NewDeepener returns a new Deepener.
func NewDeepener() *Deepener {
	return &Deepener{}
}

// MergeValues merges values read from r into root.
func (d *Deepener) MergeValues(root interface{}, r io.Reader) (interface{}, error) {
	assignments, err := newParser(r).parseAssignments()
	if err != nil {
		return nil, err
	}
	for _, assignment := range assignments {
		root = d.recursiveMerge(root, assignment.properties, assignment.value)
	}
	return root, nil
}

func (d *Deepener) recursiveMerge(root interface{}, properties []interface{}, value interface{}) interface{} {
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
			array[property] = d.recursiveMerge(array[property], properties[1:], value)
		} else {
			array = append(array, d.recursiveMerge(nil, properties[1:], value))
		}
		return array
	case string:
		object, ok := root.(map[string]interface{})
		if !ok || object == nil {
			object = make(map[string]interface{})
		}
		object[property] = d.recursiveMerge(object[property], properties[1:], value)
		return object
	default:
		panic(fmt.Sprintf("unexpected property %v (%T)", property, property))
	}
}
