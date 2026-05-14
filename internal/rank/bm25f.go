package rank

type BM25F struct {
	fieldConfigs []FieldConfig
	k1           float64
}

// todo: pull field out into file.
