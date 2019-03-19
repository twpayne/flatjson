package flatjson

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteValue(t *testing.T) {
	for _, tc := range []struct {
		json     string
		expected string
	}{
		{
			json:     `{}`,
			expected: "root = {};\n",
		},
		{
			json:     `[]`,
			expected: "root = [];\n",
		},
		{
			json:     `0`,
			expected: "root = 0;\n",
		},
		{
			json:     `""`,
			expected: "root = \"\";\n",
		},
		{
			json:     `true`,
			expected: "root = true;\n",
		},
		{
			json:     `false`,
			expected: "root = false;\n",
		},
		{
			json:     `null`,
			expected: "root = null;\n",
		},
		{
			json:     `[1,2,3]`,
			expected: "root = [];\nroot[0] = 1;\nroot[1] = 2;\nroot[2] = 3;\n",
		},
		{
			json:     `{"a":{"b":"c"}}`,
			expected: "root = {};\nroot.a = {};\nroot.a.b = \"c\";\n",
		},
		{
			json:     `{"a.b":"c"}`,
			expected: "root = {};\nroot[\"a.b\"] = \"c\";\n",
		},
		{
			json:     `{"false":false}`,
			expected: "root = {};\nroot[\"false\"] = false;\n",
		},
		{
			json:     `{"null":false}`,
			expected: "root = {};\nroot[\"null\"] = false;\n",
		},
		{
			json:     `{"true":false}`,
			expected: "root = {};\nroot[\"true\"] = false;\n",
		},
	} {
		w := &bytes.Buffer{}
		var value interface{}
		assert.NoError(t, json.Unmarshal([]byte(tc.json), &value))
		assert.NoError(t, NewEncoder(w).Encode(value))
		assert.Equal(t, tc.expected, w.String())
	}
}

func TestFlattenDecoder(t *testing.T) {
	for _, tc := range []struct {
		s         string
		prefix    string
		useNumber bool
		expected  string
	}{
		{
			s:         "0",
			prefix:    "root",
			useNumber: true,
			expected:  "root = 0;\n",
		},
		{
			s:         "0.1",
			prefix:    "root",
			useNumber: true,
			expected:  "root = 0.1;\n",
		},
		{
			s:         "0.1",
			prefix:    "root",
			useNumber: true,
			expected:  "root = 0.1;\n",
		},
		{
			s:        "true",
			prefix:   "root",
			expected: "root = true;\n",
		},
		{
			s:        "false",
			prefix:   "root",
			expected: "root = false;\n",
		},
		{
			s:        "\"\"",
			prefix:   "root",
			expected: "root = \"\";\n",
		},
		{
			s:        "null",
			prefix:   "root",
			expected: "root = null;\n",
		},
		{
			s:        "[]",
			prefix:   "root",
			expected: "root = [];\n",
		},
		{
			s:        "[true]",
			prefix:   "root",
			expected: "root = [];\nroot[0] = true;\n",
		},
		{
			s:        "[[]]",
			prefix:   "root",
			expected: "root = [];\nroot[0] = [];\n",
		},
		{
			s:        "{}",
			prefix:   "root",
			expected: "root = {};\n",
		},
		{
			s:         "{\"prop\":0}",
			prefix:    "root",
			useNumber: true,
			expected:  "root = {};\nroot.prop = 0;\n",
		},
		{
			s:         "{\"quoted.prop\":0}",
			prefix:    "root",
			useNumber: true,
			expected:  "root = {};\nroot[\"quoted.prop\"] = 0;\n",
		},
	} {
		d := json.NewDecoder(bytes.NewBufferString(tc.s))
		if tc.useNumber {
			d.UseNumber()
		}
		w := &bytes.Buffer{}
		e := NewEncoder(w, EncoderPrefix(tc.prefix))
		assert.NoError(t, e.Transcode(d))
		assert.Equal(t, tc.expected, w.String())
	}
}
