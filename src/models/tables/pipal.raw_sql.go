package tables

// RawSQL 应用配置
type RawSQL struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// TableName 设置表名
func (RawSQL) TableName() string {
	return "pipal.raw_sql"
}
