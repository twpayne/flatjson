package flatjson

import (
	"bytes"
	"testing"
)

func TestWriteValues(t *testing.T) {
	for _, tc := range []struct {
		prefix string
		json   string
		want   string
	}{
		{
			json: `{}`,
			want: "root = {};\n",
		},
		{
			json: `[]`,
			want: "root = [];\n",
		},
		{
			json: `0`,
			want: "root = 0;\n",
		},
		{
			json: `""`,
			want: "root = \"\";\n",
		},
		{
			json: `true`,
			want: "root = true;\n",
		},
		{
			json: `false`,
			want: "root = false;\n",
		},
		{
			json: `null`,
			want: "root = null;\n",
		},
		{
			json: `[1,2,3]`,
			want: "root = [];\nroot[0] = 1;\nroot[1] = 2;\nroot[2] = 3;\n",
		},
		{
			json: `{"a":{"b":"c"}}`,
			want: "root = {};\nroot.a = {};\nroot.a.b = \"c\";\n",
		},
		{
			json: `{"a.b":"c"}`,
			want: "root = {};\nroot[\"a.b\"] = \"c\";\n",
		},
		{
			json: `{"false":false}`,
			want: "root = {};\nroot[\"false\"] = false;\n",
		},
		{
			json: `{"null":false}`,
			want: "root = {};\nroot[\"null\"] = false;\n",
		},
		{
			json: `{"true":false}`,
			want: "root = {};\nroot[\"true\"] = false;\n",
		},
	} {
		b := &bytes.Buffer{}
		if err := NewFlattener(b).WriteValues([]byte(tc.json)); err != nil {
			t.Errorf("WriteValues(%q, %q) == %v, want <nil>", tc.prefix, tc.json, err)
		}
		if got := b.String(); got != tc.want {
			t.Errorf("WriteValues(%q, %q) wrote:\n%s\nwant:%s", tc.prefix, tc.json, got, tc.want)
		}
	}
}
