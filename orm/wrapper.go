package orm

import (
	"encoding/json"
	"strings"

	"github.com/mocheer/charon/global"
	"github.com/mocheer/charon/model"
	"github.com/mocheer/pluto/fn"
	"github.com/mocheer/pluto/ref"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

type Wrapper struct {
	Ctx *gorm.DB
}

// Model 设置数据模型
func (o *Wrapper) Model(name string) *Wrapper {
	o.Ctx = o.Ctx.Model(model.NewTableStruct(name))
	return o
}

// Select 字段选择
func (o *Wrapper) Select(data string) *Wrapper {
	if data != "" {
		// 不止是字段选择，还是字段重命名，且支持函数调用
		o.Ctx.Select(strings.Split(data, ","))
	}
	return o
}

// Joins
func (o *Wrapper) Joins(data []string) *Wrapper {
	for _, joinStr := range data {
		o.Ctx.Joins(joinStr)
	}
	return o
}

// Where
func (o *Wrapper) Where(data string) *Wrapper {
	r := gjson.Parse(data)
	switch r.Type {
	case gjson.JSON:
		o.Ctx.Where(r.Value())
	default:
		o.Ctx.Where(data)
	}
	return o
}

//
func (o *Wrapper) Not(data string) *Wrapper {
	r := gjson.Parse(data)
	switch r.Type {
	case gjson.JSON:
		o.Ctx.Not(r.Value())
	default:
		o.Ctx.Not(data)
	}
	return o
}

//
func (o *Wrapper) Order(data string) *Wrapper {
	if data != "" {
		o.Ctx.Order(data)
	}
	return o
}

func (o *Wrapper) Limit(data int) *Wrapper {
	if data != 0 {
		o.Ctx.Limit(data)
	}

	return o
}

func (o *Wrapper) Offset(data int) *Wrapper {
	if data != 0 {
		o.Ctx.Offset(data)
	}
	return o
}

//
func (o *Wrapper) Query(args *SelectArgs) interface{} {
	o.Model(args.Name).Select(args.Select).Joins(args.Joins).Where(args.Where).Not(args.Not).Order(args.Order).Limit(args.Limit).Offset(args.Offset)
	//
	entity := model.NewTableStruct(args.Name)
	switch args.Mode {
	case "first":
		o.Ctx.First(entity)
		return entity
	case "last":
		o.Ctx.Last(entity)
		return entity
	case "find":
		var result = model.NewTableStructArray(args.Name)
		o.Ctx.Find(result)
		return result
	default:
		var result []map[string]interface{}
		o.ScanIntoMap(&result)
		return result
	}
}

func (o *Wrapper) NewModelEntity(data []byte) interface{} {
	entity := ref.New(o.Ctx.Statement.Model)
	err := json.Unmarshal([]byte(data), entity)
	if err != nil {
		panic(err)
	}
	return entity
}

// Create 创建
func (o *Wrapper) Create(data []byte) *Wrapper {
	o.Ctx = o.Ctx.Create(o.NewModelEntity(data)) //RowsAffected 才能生效
	return o
}

// Update 修改
func (o *Wrapper) Updates(data []byte) *Wrapper {
	o.Ctx = o.Ctx.Updates(o.NewModelEntity(data)) //RowsAffected 才能生效
	return o
}

// Raw
func (o *Wrapper) Raw(sql string, values ...interface{}) *Wrapper {
	o.Ctx = o.Ctx.Raw(sql, values...)
	return o
}

// Success 是否执行成功
func (o *Wrapper) Success() bool {
	return o.Ctx.RowsAffected > 0
}

// ScanIntoMap 扫描数据
func (o *Wrapper) ScanIntoMap(data *[]map[string]interface{}) {
	rows, err := o.Ctx.Rows()
	result := *data
	if err == nil {
		for rows.Next() {
			rowData := map[string]interface{}{}
			err = global.DB.ScanRows(rows, &rowData)
			//
			if err == nil {
				columnTypes, _ := rows.ColumnTypes()
				for _, columnType := range columnTypes {
					// gorm-postgresql 底层的 DataType 并没有识别JSON，然后转成[]byte，而是用的string
					if columnType.DatabaseTypeName() == "JSON" {
						name := columnType.Name()
						if fn.IsString(rowData[name]) {
							var data interface{}
							err := json.Unmarshal([]byte(rowData[name].(string)), &data)
							if err == nil {
								rowData[name] = data
							}
						}
					}
				}
			}
			//
			result = append(result, rowData)
		}
	}
	*data = result
}
