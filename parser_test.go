package flatjson

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAssignment(t *testing.T) {
	for _, tc := range []struct {
		s                  string
		expectedAssignment *assignment
		expectErr          bool
	}{
		{
			s: "root = 0;\n",
			expectedAssignment: &assignment{
				identifier: "root",
				value:      json.Number("0"),
			},
		},
		{
			s: "root = [];\n",
			expectedAssignment: &assignment{
				identifier: "root",
				value:      []interface{}{},
			},
		},
		{
			s: "root[0] = true;\n",
			expectedAssignment: &assignment{
				identifier: "root",
				properties: []interface{}{0},
				value:      true,
			},
		},
		{
			s: "root = {};\n",
			expectedAssignment: &assignment{
				identifier: "root",
				value:      map[string]interface{}{},
			},
		},
		{
			s: "root.prop = {};\n",
			expectedAssignment: &assignment{
				identifier: "root",
				properties: []interface{}{"prop"},
				value:      map[string]interface{}{},
			},
		},
		{
			s: "root[\"prop\"] = null;\n",
			expectedAssignment: &assignment{
				identifier: "root",
				properties: []interface{}{"prop"},
				value:      nil,
			},
		},
		{
			s: "// comment\nroot = 0;\n",
			expectedAssignment: &assignment{
				identifier: "root",
				value:      json.Number("0"),
			},
		},
		{
			s: "root = 0; // comment\n",
			expectedAssignment: &assignment{
				identifier: "root",
				value:      json.Number("0"),
			},
		},
		{
			s: "/*comment*/root/*comment*/=/*comment*/0/*comment*/;/*comment*/\n",
			expectedAssignment: &assignment{
				identifier: "root",
				value:      json.Number("0"),
			},
		},
	} {
		assignment, err := newParser(bytes.NewBufferString(tc.s)).parseAssignment()
		if tc.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAssignment, assignment)
		}
	}
}

func TestParsePropertyAccess(t *testing.T) {
	for _, tc := range []struct {
		s                string
		expectedProperty interface{}
		expectErr        bool
	}{
		{s: "", expectErr: true},
		{s: ".", expectErr: true},
		{s: ".a", expectedProperty: "a"},
		{s: ".aB0_", expectedProperty: "aB0_"},
		{s: ".0", expectErr: true},
		{s: "[0]", expectedProperty: 0},
		{s: "[123]", expectedProperty: 123},
		{s: "[0", expectErr: true},
		{s: "[a", expectErr: true},
		{s: "[\"\"]", expectedProperty: ""},
		{s: "[\"a\"]", expectedProperty: "a"},
	} {
		actualProperty, err := newParser(bytes.NewBufferString(tc.s)).parsePropertyAccess()
		if tc.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedProperty, actualProperty)
		}
	}
}

func TestParseValue(t *testing.T) {
	for _, tc := range []struct {
		s             string
		expectedValue interface{}
		expectErr     bool
	}{
		{s: "0", expectedValue: json.Number("0")},
		{s: "false", expectedValue: false},
		{s: "true", expectedValue: true},
		{s: "null", expectedValue: nil},
		{s: "\"\"", expectedValue: ""},
		{s: "\"/*comment*/\"", expectedValue: "/*comment*/"},
		{s: "\"//comment\n\"", expectedValue: "//comment\n"},
		{s: "[]", expectedValue: []interface{}{}},
		{s: "{}", expectedValue: map[string]interface{}{}},
	} {
		actualValue, err := newParser(bytes.NewBufferString(tc.s)).parseValue()
		if tc.expectErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, actualValue)
		}
	}
}
