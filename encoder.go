package flatjson

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// A Encoder flattens JSON.
type Encoder struct {
	w      io.Writer
	prefix string
	suffix string
}

// A EncoderOption modifies a Encoder.
type EncoderOption func(*Encoder)

// NewEncoder returns a new Encoder.
func NewEncoder(w io.Writer, options ...EncoderOption) *Encoder {
	e := &Encoder{
		w:      w,
		prefix: "root",
		suffix: ";\n",
	}
	for _, option := range options {
		option(e)
	}
	return e
}

// Transcode reads JSON tokens from d and writes their flatjson
// representation.
func (e *Encoder) Transcode(d *json.Decoder) error {
	return e.writeDecoderHelper(d, e.prefix)
}

func (e *Encoder) writeDecoderHelper(d *json.Decoder, prefix string) error {
	token, err := d.Token()
	if err != nil {
		return err
	}
	switch token := token.(type) {
	case bool:
		_, err := e.w.Write([]byte(prefix + " = " + strconv.FormatBool(token) + e.suffix))
		return err
	case float64:
		_, err := e.w.Write([]byte(prefix + " = " + strconv.FormatFloat(token, 'e', -1, 64) + e.suffix))
		return err
	case json.Delim:
		switch token {
		case '[':
			index := 0
			if _, err := e.w.Write([]byte(prefix + " = []" + e.suffix)); err != nil {
				return err
			}
			for d.More() {
				err := e.writeDecoderHelper(d, prefix+"["+strconv.Itoa(index)+"]")
				if err != nil {
					return err
				}
				index++
			}
			_, err := d.Token() // ']'
			return err
		case '{':
			if _, err := e.w.Write([]byte(prefix + " = {}" + e.suffix)); err != nil {
				return err
			}
			for d.More() {
				propertyToken, err := d.Token()
				if err != nil {
					return err
				}
				property, ok := propertyToken.(string)
				if !ok {
					return fmt.Errorf("expected a string, got %v", propertyToken)
				}
				if err := e.writeDecoderHelper(d, propertyAccessor(prefix, property)); err != nil {
					return err
				}
			}
			_, err := d.Token() // '}'
			return err
		default:
			return fmt.Errorf("unexpected delimiter: %v", token)
		}
	case json.Number:
		_, err := e.w.Write([]byte(prefix + " = " + token.String() + e.suffix))
		return err
	case string:
		_, err := e.w.Write([]byte(prefix + " = " + strconv.QuoteToASCII(token) + e.suffix))
		return err
	case nil:
		_, err := e.w.Write([]byte(prefix + " = null" + e.suffix))
		return err
	default:
		return fmt.Errorf("unknown token: %v", token)
	}
}

// Encode encodes value.
func (e *Encoder) Encode(value interface{}) error {
	return e.encodeHelper(e.prefix, value)
}

func (e *Encoder) encodeArray(prefix string, array []interface{}) error {
	if _, err := e.w.Write([]byte(prefix + " = []" + e.suffix)); err != nil {
		return err
	}
	for i, value := range array {
		if err := e.encodeHelper(prefix+"["+strconv.Itoa(i)+"]", value); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeObject(prefix string, object map[string]interface{}) error {
	if _, err := e.w.Write([]byte(prefix + " = {}" + e.suffix)); err != nil {
		return err
	}
	properties := make([]string, 0, len(object))
	for property := range object {
		properties = append(properties, property)
	}
	sort.Strings(properties)
	for _, property := range properties {
		if err := e.encodeHelper(propertyAccessor(prefix, property), object[property]); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeHelper(prefix string, value interface{}) error {
	switch value := value.(type) {
	case []interface{}:
		return e.encodeArray(prefix, value)
	case map[string]interface{}:
		return e.encodeObject(prefix, value)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		_, err = e.w.Write([]byte(prefix + " = " + string(data) + e.suffix))
		return err
	}
}

// Marshal returns the flatjson encoding of v.
func Marshal(v interface{}, options ...EncoderOption) ([]byte, error) {
	sb := &strings.Builder{}
	if err := NewEncoder(sb, options...).Encode(v); err != nil {
		return nil, err
	}
	return []byte(sb.String()), nil
}

// EncoderPrefix sets the prefix.
func EncoderPrefix(prefix string) EncoderOption {
	return func(e *Encoder) {
		e.prefix = prefix
	}
}

// EncoderSuffix sets the suffix.
func EncoderSuffix(suffix string) EncoderOption {
	return func(e *Encoder) {
		e.suffix = suffix
	}
}
