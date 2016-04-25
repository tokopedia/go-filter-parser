package filter

import (
	"testing"
)

type FilterAny struct {
	AnyString      String
	AnyStringSlice StringSlice
	AnyBool        Bool
	AnyFloat       Float64
	AnyFloatSlice  Float64Slice
	AnyFloatRange  Float64Range
}

func (f *FilterAny) FilterMap() FilterMap {
	return FilterMap{
		&f.AnyString:      "any_string",
		&f.AnyStringSlice: "any_string_slice",
		&f.AnyBool:        "any_bool",
		&f.AnyFloat:       "any_float",
		&f.AnyFloatSlice:  "any_float_slice",
		&f.AnyFloatRange:  "any_float_range",
	}
}

func TestParse(t *testing.T) {
	fq := `any_string==hello world\;foo\;bar;any_string_slice==hello\:world:foo:bar;any_bool!=True;any_float==12345.6789;any_float_slice!=1:2.3:4:5.6:7890;any_float_range!=5000..10000`

	f := new(FilterAny)
	if err := Parse(fq, f); err != nil {
		t.Error(err)
	}

	if f.AnyString.Value != `hello world;foo;bar` {
		t.Errorf("invalid output for %T type, expecting %q, got %q.", f.AnyString, `hello world;foo;bar`, f.AnyString.Value)
	}

	if len(f.AnyStringSlice.Value) != 3 {
		t.Errorf("invalid output for %T type, expecting %q, got %q.", f.AnyStringSlice, `[hello:world foo bar]`, f.AnyStringSlice.Value)
	}

	if f.AnyBool.Value != true {
		t.Errorf("invalid output for %T type, expecting \"%v\", got \"%v\".", f.AnyBool, true, f.AnyBool.Value)
	}

	if f.AnyFloat.Value != 12345.6789 {
		t.Errorf("invalid output for %T type, expecting \"%f\", got \"%f\".", f.AnyFloat, 12345.6789, f.AnyFloat.Value)
	}

	if len(f.AnyFloatSlice.Value) != 5 {
		t.Errorf("invalid output for %T type, expecting \"%v\", got \"%v\".", f.AnyFloatSlice, []float64{1, 2.3, 4, 5.6, 7890}, f.AnyFloatSlice.Value)
	}

	if f.AnyFloatRange.Value[0] != 5000 || f.AnyFloatRange.Value[1] != 10000 || f.AnyFloatRange.Value[0] > f.AnyFloatRange.Value[1] {
		t.Errorf("invalid output for %T type, expecting \"%v\", got \"%v\".", f.AnyFloatRange, [2]float64{5000, 10000}, f.AnyFloatRange.Value)
	}
}
