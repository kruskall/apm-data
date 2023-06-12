package modelpb

// Labels wraps a map[string]string or map[string][]string with utility
// methods.
type Labels map[string]*LabelValue

// Set sets the label k to value v. If there existed a label in l with the same
// key, it will be replaced and its Global field will be set to false.
func (l Labels) Set(k string, v string) {
	l[k] = &LabelValue{Value: v}
}

// SetSlice sets the label k to value v. If there existed a label in l with the
// same key, it will be replaced and its Global field will be set to false.
func (l Labels) SetSlice(k string, v []string) {
	l[k] = &LabelValue{Values: v}
}

// Clone creates a deep copy of Labels.
func (l Labels) Clone() Labels {
	cp := make(Labels)
	for k, v := range l {
		to := LabelValue{Global: v.Global, Value: v.Value}
		if len(v.Values) > 0 {
			to.Values = make([]string, len(v.Values))
			copy(to.Values, v.Values)
		}
		cp[k] = &to
	}
	return cp
}

// NumericLabels wraps a map[string]float64 or map[string][]float64 with utility
// methods.
type NumericLabels map[string]*NumericLabelValue

// Set sets the label k to value v. If there existed a label in l with the same
// key, it will be replaced and its Global field will be set to false.
func (l NumericLabels) Set(k string, v float64) {
	l[k] = &NumericLabelValue{Value: v}
}

// SetSlice sets the label k to value v. If there existed a label in l with the
// same key, it will be replaced and its Global field will be set to false.
func (l NumericLabels) SetSlice(k string, v []float64) {
	l[k] = &NumericLabelValue{Values: v}
}

// Clone creates a deep copy of NumericLabels.
func (l NumericLabels) Clone() NumericLabels {
	cp := make(NumericLabels)
	for k, v := range l {
		to := NumericLabelValue{Global: v.Global, Value: v.Value}
		if len(v.Values) > 0 {
			to.Values = make([]float64, len(v.Values))
			copy(to.Values, v.Values)
		}
		cp[k] = &to
	}
	return cp
}
