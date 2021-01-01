package models

//postgresql type -> go type
var PgTypeMap = map[string]string{
		"int4":      "int32",
		"int8":      "int64",
		"float4":    "float32",
		"float8":    "float64",
		"double":    "float64",
		"varchar":   "string",
		"boolean":   "bool",
		"timestamp": "time.Time",
		"date":      "time.Time",
}
