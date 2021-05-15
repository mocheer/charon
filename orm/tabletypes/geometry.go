package tabletypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Geometry json.RawMessage

// Value return json value, implement driver.Valuer interface
func (g Geometry) Value() (driver.Value, error) {
	if len(g) == 0 {
		return nil, nil
	}
	bytes, err := json.RawMessage(g).MarshalJSON()
	return string(bytes), err
}

//
func (g *Geometry) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal Geometry value:", value))
	}
	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*g = Geometry(result)
	return err
}

func (g Geometry) MarshalJSON() ([]byte, error) {
	return json.RawMessage(g).MarshalJSON()
}

// UnmarshalJSON to deserialize []byte
func (g *Geometry) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*g = Geometry(result)
	return err
}

func (g Geometry) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{GeometryQueryClause{Field: f}}
}

type GeometryQueryClause struct {
	Field *schema.Field
}

func (gqc GeometryQueryClause) Name() string {
	return ""
}

func (gqc GeometryQueryClause) Build(clause.Builder) {
}

func (gqc GeometryQueryClause) MergeClause(*clause.Clause) {
}

func (gqc GeometryQueryClause) ModifyStatement(stmt *gorm.Statement) {

	if stmt.Selects == nil {
		selects := []clause.Column{}
		for _, name := range stmt.Schema.DBNames {
			raw := false
			if name == gqc.Field.DBName {
				name = fmt.Sprintf("ST_ASGeojson(%s,4) as %s", gqc.Field.DBName, gqc.Field.DBName)
				raw = true
			}
			selects = append(selects, clause.Column{Name: name, Raw: raw})
		}
		stmt.AddClause(clause.Select{
			Columns: selects,
		})
	}

}
