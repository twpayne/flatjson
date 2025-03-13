package flatjson

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestScanComment(t *testing.T) {
	for i, tc := range []struct {
		s         string
		expectErr bool
	}{
		{s: "//"},
		{s: "//\n"},
		{s: "/**/"},
		{s: "/*\n*/"},
		{s: "/", expectErr: true},
		{s: "/*", expectErr: true},
		{s: "/*/", expectErr: true},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tok, lit := newScanner(bytes.NewBufferString(tc.s)).scanComment()
			if tc.expectErr {
				assert.Equal(t, tokenIllegal, tok)
				assert.Equal(t, tc.s, lit)
			} else {
				assert.Equal(t, tokenComment, tok)
				assert.Equal(t, tc.s, lit)
			}
		})
	}
}

func TestScanNumber(t *testing.T) {
	for i, tc := range []struct {
		s         string
		expectErr bool
	}{
		{s: "0"},
		{s: "-1"},
		{s: "-123"},
		{s: "1"},
		{s: "123"},
		{s: "0.1"},
		{s: "0.123"},
		{s: "-0.1"},
		{s: "-0.123"},
		{s: "0e", expectErr: true},
		{s: "0e0"},
		{s: "0E0"},
		{s: "0e+0"},
		{s: "0e-0"},
		{s: "-1.23e+45"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tok, lit := newScanner(bytes.NewBufferString(tc.s)).scanNumber()
			if tc.expectErr {
				assert.Equal(t, tokenIllegal, tok)
				assert.Equal(t, "", lit)
			} else {
				assert.Equal(t, tokenNumber, tok)
				assert.Equal(t, tc.s, lit)
			}
		})
	}
}
