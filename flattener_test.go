package flatjson

import (
	"strconv"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestWriteValues(t *testing.T) {
	for i, tc := range []struct {
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
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sb := &strings.Builder{}
			assert.NoError(t, NewFlattener(sb).WriteValues([]byte(tc.json)))
			assert.Equal(t, tc.expected, sb.String())
		})
	}
}
