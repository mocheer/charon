package orm

// SelectArgs 构建查询语句的参数
type SelectArgs struct {
	Name   string   `query:"name"`   // 表名
	Mode   string   `query:"mode"`   // 数据返回模式 take | first | last | find | default
	Where  string   `query:"where"`  // 查询 a=1 || {a:1}  => where a=1
	Not    string   `query:"not"`    // 查询 a=1 || {a:1}  => where not a=1
	Select string   `query:"select"` // 字段 a as b, sum(a) as b
	Limit  int      `query:"limit"`  // 限制数据量
	Offset int      `query:"offset"` // 偏移位置
	Order  string   `query:"order"`  // 排序
	Joins  []string `query:"joins"`  // 多表关联查询
}
