package rank

type Field string

type FieldConfig struct {
	// Weight is the relative weight of the field.
	// The field that has the bulk of the content should be 1.
	Weight float64

	// B is the strength of length normalization.
	// For 0, no normalization is performed.
	// For 1, results ar scaled to the average document length.
	B float64
}
