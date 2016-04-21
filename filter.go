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

type String struct {
	Operator int
	Value    string
}

type Bool struct {
	Operator int
	Value    bool
}

type Float64 struct {
	Operator int
	Value    float64
}

type Range struct {
	Operator int
	Value    [2]float64
}

type Filter interface {
	FilterMap() FilterMap
}

type FilterMap map[interface{}]string

func Parse(fs string, f Filter) error {
	tm, err := mapFilter(fs)
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
			if tm[v].Operator != OperatorEqualTo && tm[v].Operator != OperatorNotEqualTo {
				return fmt.Errorf("unsupported operator for \"%T\" data type.", k)
			}
			k.(*String).Operator = tm[v].Operator
			k.(*String).Value = tm[v].Value
		case *Bool:
			if tm[v].Operator != OperatorEqualTo && tm[v].Operator != OperatorNotEqualTo {
				return fmt.Errorf("unsupported operator for \"%T\" data type.", k)
			}
			bv := false
			if strings.ToLower(tm[v].Value) == "true" {
				bv = true
			}
			k.(*Bool).Operator = tm[v].Operator
			k.(*Bool).Value = bv
		case *Float64:
			nv, err := strconv.ParseFloat(tm[v].Value, 64)
			if err != nil {
				return fmt.Errorf("unable to parse %q for data type \"%T\"", tm[v].Value, k)
			}
			k.(*Float64).Operator = tm[v].Operator
			k.(*Float64).Value = nv
		case *Range:
			if tm[v].Operator != OperatorEqualTo && tm[v].Operator != OperatorNotEqualTo {
				return fmt.Errorf("unsupported operator for \"%T\" data type.", k)
			}
			td := strings.Split(tm[v].Value, "..")
			if len(td) < 2 {
				return fmt.Errorf("unable to parse %q for data type \"%T\"", tm[v].Value, k)
			}
			rv1, err := strconv.ParseFloat(td[0], 64)
			if err != nil {
				return fmt.Errorf("unable to parse %q for data type \"%T\"", td[0], k)
			}
			rv2, err := strconv.ParseFloat(td[1], 64)
			if err != nil {
				return fmt.Errorf("unable to parse %q for data type \"%T\"", td[1], k)
			}
			if rv1 > rv2 {
				return fmt.Errorf("\"%T\" data type should be in ascending order", k)
			}
			k.(*Range).Operator = tm[v].Operator
			k.(*Range).Value = [2]float64{rv1, rv2}
		default:
			return fmt.Errorf("unsupported data type \"%T\".", k)
		}
	}
	return nil
}

func mapFilter(fs string) (map[string]String, error) {
	fa := strings.Split(fs, ";")
	tm := map[string]String{}
	for _, v := range fa {
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
