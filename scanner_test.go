package flatjson

import (
	"bytes"
	"testing"
)

func TestScanComment(t *testing.T) {
	for _, tc := range []struct {
		s       string
		wantErr bool
	}{
		{s: "//"},
		{s: "//\n"},
		{s: "/**/"},
		{s: "/*\n*/"},
		{s: "/", wantErr: true},
		{s: "/*", wantErr: true},
		{s: "/*/", wantErr: true},
	} {
		tok, lit := newScanner(bytes.NewBufferString(tc.s)).scanComment()
		if tc.wantErr {
			if tok != tokenIllegal || lit != tc.s {
				t.Errorf("newScanner(%q).scanComment() == %s, %q, want tokenIllegal, %q", tc.s, tok, lit, tc.s)
			}
		} else {
			if tok != tokenComment || lit != tc.s {
				t.Errorf("newScanner(%q).scanComment() == %s, %q, want tokenComment, %q", tc.s, tok, lit, tc.s)
			}
		}
	}
}

func TestScanNumber(t *testing.T) {
	for _, tc := range []struct {
		s       string
		wantErr bool
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
		{s: "0e", wantErr: true},
		{s: "0e0"},
		{s: "0E0"},
		{s: "0e+0"},
		{s: "0e-0"},
		{s: "-1.23e+45"},
	} {
		tok, lit := newScanner(bytes.NewBufferString(tc.s)).scanNumber()
		if tc.wantErr {
			if tok != tokenIllegal || lit != "" {
				t.Errorf("newScanner(%q).scanNumber() == %s, %q, want tokenIllegal, \"\"", tc.s, tok, lit)
			}
		} else {
			if tok != tokenNumber || lit != tc.s {
				t.Errorf("newScanner(%q).scanNumber() == %s, %q, want tokenNumber, %q", tc.s, tok, lit, tc.s)
			}
		}
	}
}
