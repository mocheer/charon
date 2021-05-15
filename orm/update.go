package orm

// UpdateArgs 构建查询语句的参数
type UpdateArgs struct {
	Name  string `query:"name"`  // 表名
	Where string `query:"where"` // 查询 a=1 || {a:1}  => where a=1
}
