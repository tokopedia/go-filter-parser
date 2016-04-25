package filter

type String struct {
	Operator int
	Value    string
}

type StringSlice struct {
	Operator int
	Value    []string
}

type Bool struct {
	Operator int
	Value    bool
}

type Float64 struct {
	Operator int
	Value    float64
}

type Float64Slice struct {
	Operator int
	Value    []float64
}

type Float64Range struct {
	Operator int
	Value    [2]float64
}
