package flatjson

import (
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestFlatJSON(t *testing.T) {
	for i, tc := range []struct {
		v  interface{}
		fj string
	}{
		{
			v:  nil,
			fj: "root = null;\n",
		},
		{
			v:  true,
			fj: "root = true;\n",
		},
		{
			v:  false,
			fj: "root = false;\n",
		},
		{
			v:  0,
			fj: "root = 0;\n",
		},
		{
			v:  []interface{}{},
			fj: "root = [];\n",
		},
		{
			v:  map[string]interface{}{},
			fj: "root = {};\n",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gotFJ, err := Marshal(tc.v)
			assert.NoError(t, err)
			assert.Equal(t, tc.fj, string(gotFJ))
			actualV := tc.v
			assert.NoError(t, Unmarshal([]byte(tc.fj), &actualV))
			assert.Equal(t, tc.v, actualV)
		})
	}
}
