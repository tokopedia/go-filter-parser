package filter

import (
	"testing"
)

type FilterAny struct {
	AnyString String
	AnyBool   Bool
	AnyFloat  Float64
	AnyRange  Range
}

func (f *FilterAny) FilterMap() FilterMap {
	return FilterMap{
		&f.AnyString: "any_string",
		&f.AnyBool:   "any_bool",
		&f.AnyFloat:  "any_float",
		&f.AnyRange:  "any_range",
	}
}

func TestParse(t *testing.T) {
	fq := `any_string==hello world\;foo\;bar;any_bool!=True;any_float==12345.6789;any_range!=5000..10000`

	f := new(FilterAny)
	if err := Parse(fq, f); err != nil {
		t.Error(err)
	}

	if f.AnyString.Value != `hello world;foo;bar` {
		t.Errorf("invalid output for %T type, expecting %q, got %q", f.AnyString, f.AnyString.Value, `hello world;foo;bar`)
	}

	if f.AnyBool.Value != true {
		t.Errorf("invalid output for %T type, expecting \"%v\", got \"%v\"", f.AnyBool, f.AnyBool.Value, true)
	}

	if f.AnyFloat.Value != 12345.6789 {
		t.Errorf("invalid output for %T type, expecting \"%f\", got \"%f\"", f.AnyFloat, f.AnyFloat.Value, 12345.6789)
	}

	if f.AnyRange.Value[0] != 5000 || f.AnyRange.Value[1] != 10000 {
		t.Errorf("invalid output for %T type, expecting \"%f\" and \"%f\", got \"%f\" and \"%f\"", f.AnyRange, f.AnyRange.Value[0], f.AnyRange.Value[1], 5000, 10000)
	}
}
