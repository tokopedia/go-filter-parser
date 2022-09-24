package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	rxFilter = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9_]*?)(==|!=|>=?|<=?)(.+)$`)
)

const (
	OperatorEqualTo int = iota
	OperatorNotEqualTo
	OperatorGreaterThan
	OperatorLessThan
	OperatorGreaterThanOrEqualTo
	OperatorLessThanOrEqualTo
)

var operatorText map[int]string

func init() {
	operatorText = map[int]string{
		OperatorEqualTo:              "==",
		OperatorNotEqualTo:           "!=",
		OperatorGreaterThan:          ">",
		OperatorLessThan:             "<",
		OperatorGreaterThanOrEqualTo: ">=",
		OperatorLessThanOrEqualTo:    "<=",
	}
}

func OperatorText(op int) string {
	if v, ok := operatorText[op]; ok {
		return v
	}
	return ""
}

type Filter interface {
	FilterMap() FilterMap
}

type FilterMap map[interface{}]string

func Parse(fs, sep string, f Filter) error {
	if fs == "" {
		return nil
	}

	tm, err := mapFilter(fs, sep)
	if err != nil {
		return err
	}

	fm := f.FilterMap()
	for k, v := range fm {
		if _, ok := tm[v]; !ok {
			continue
		}

		switch k.(type) {
		case *String:
			if err := mapString(k.(*String), tm[v]); err != nil {
				return err
			}
		case *StringSlice:
			if err := mapStringSlice(k.(*StringSlice), tm[v]); err != nil {
				return err
			}
		case *Bool:
			if err := mapBool(k.(*Bool), tm[v]); err != nil {
				return err
			}
		case *Float64:
			if err := mapFloat64(k.(*Float64), tm[v]); err != nil {
				return err
			}
		case *Float64Slice:
			if err := mapFloat64Slice(k.(*Float64Slice), tm[v]); err != nil {
				return err
			}
		case *Float64Range:
			if err := mapFloat64Range(k.(*Float64Range), tm[v]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported data type \"%T\".", k)
		}
	}
	return nil
}

func mapFilter(fs, sep string) (map[string]String, error) {
	fs = strings.Replace(fs, "\\;", "{{FILTER_SEMICOLON}}", -1)
	fa := strings.Split(fs, sep)
	tm := map[string]String{}
	for _, v := range fa {
		v = strings.Replace(v, "{{FILTER_SEMICOLON}}", ";", -1)
		kv := rxFilter.FindAllStringSubmatch(v, -1)
		if len(kv) < 1 {
			return nil, fmt.Errorf("unable to parse %q.", v)
		}
		op := 0
		switch kv[0][2] {
		case "==":
			op = OperatorEqualTo
		case "!=":
			op = OperatorNotEqualTo
		case ">":
			op = OperatorGreaterThan
		case "<":
			op = OperatorLessThan
		case ">=":
			op = OperatorGreaterThanOrEqualTo
		case "<=":
			op = OperatorLessThanOrEqualTo
		default:
			return nil, fmt.Errorf("unsupported operator.")
		}
		tm[kv[0][1]] = String{op, kv[0][3]}
	}
	return tm, nil
}

func mapString(r *String, m String) error {
	if m.Operator != OperatorEqualTo && m.Operator != OperatorNotEqualTo {
		return fmt.Errorf("unsupported operator for \"%T\" data type.", r)
	}
	r.Operator = m.Operator
	r.Value = m.Value
	return nil
}

func mapStringSlice(r *StringSlice, m String) error {
	if m.Operator != OperatorEqualTo && m.Operator != OperatorNotEqualTo {
		return fmt.Errorf("unsupported operator for \"%T\" data type.", r)
	}

	mv := strings.Replace(m.Value, "\\:", "{{FILTER_COLON}}", -1)
	ma := strings.Split(mv, ":")
	for k, v := range ma {
		ma[k] = strings.Replace(v, "{{FILTER_COLON}}", ":", -1)
	}

	r.Operator = m.Operator
	r.Value = ma
	return nil
}

func mapBool(r *Bool, m String) error {
	if m.Operator != OperatorEqualTo && m.Operator != OperatorNotEqualTo {
		return fmt.Errorf("unsupported operator for \"%T\" data type.", r)
	}
	bv := false
	if strings.ToLower(strings.Trim(m.Value, " ")) == "true" {
		bv = true
	}
	r.Operator = m.Operator
	r.Value = bv
	return nil
}

func mapFloat64(r *Float64, m String) error {
	nv, err := strconv.ParseFloat(m.Value, 64)
	if err != nil {
		return fmt.Errorf("unable to parse %q for data type \"%T\"", m.Value, r)
	}
	r.Operator = m.Operator
	r.Value = nv
	return nil
}

func mapFloat64Slice(r *Float64Slice, m String) error {
	if m.Operator != OperatorEqualTo && m.Operator != OperatorNotEqualTo {
		return fmt.Errorf("unsupported operator for \"%T\" data type.", r)
	}

	td := strings.Split(m.Value, ":")
	fa := make([]float64, len(td))
	for k, v := range td {
		tf, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("unable to parse %q for data type \"%T\"", v, r)
		}
		fa[k] = tf
	}
	r.Operator = m.Operator
	r.Value = fa
	return nil
}

func mapFloat64Range(r *Float64Range, m String) error {
	if m.Operator != OperatorEqualTo && m.Operator != OperatorNotEqualTo {
		return fmt.Errorf("unsupported operator for \"%T\" data type.", r)
	}
	td := strings.Split(m.Value, "..")
	if len(td) < 2 {
		return fmt.Errorf("unable to parse %q for data type \"%T\"", m.Value, r)
	}
	rv1, err := strconv.ParseFloat(td[0], 64)
	if err != nil {
		return fmt.Errorf("unable to parse %q for data type \"%T\"", td[0], r)
	}
	rv2, err := strconv.ParseFloat(td[1], 64)
	if err != nil {
		return fmt.Errorf("unable to parse %q for data type \"%T\"", td[1], r)
	}
	if rv1 > rv2 {
		return fmt.Errorf("\"%T\" data type should be in ascending order", r)
	}
	r.Operator = m.Operator
	r.Value = [2]float64{rv1, rv2}
	return nil
}
