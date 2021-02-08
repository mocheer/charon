package tables

// User struct
type User struct {
	Name     string `gorm:"unique_index;not null" json:"name"`
	Password string `gorm:"not null" json:"password"`
	Nickname string
	Role     int // -1 未初始化，0 游客 ， 1 管理员
}

// TableName 设置表名
func (User) TableName() string {
	return "pipal.user"
}
