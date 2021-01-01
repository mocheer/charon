package tables

import "gorm.io/datatypes"

// Petiole 页面配置
type Petiole struct {
	Stipule string         `json:"stipule"`
	Name    string         `json:"name"`
	Dep     string         `json:"dep"`
	Guard   bool           `json:"guard"`
	Props   datatypes.JSON `json:"props"`
	Options datatypes.JSON `json:"options"`
	Data    string         `json:"data"` // Data 字段的原始类型是byte，这里用string，相当于convert(data,'utf-8') => 目前发现直接基于string转换，有特殊字符的情况下修改失败
	Remark  string         `json:"remark"`
}

// TableName 设置表名
func (Petiole) TableName() string {
	return "pipal.petiole"
}
