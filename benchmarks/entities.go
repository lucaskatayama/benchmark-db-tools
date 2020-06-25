package benchmarks

type Model struct {
	ID    int    `dbq:"id" gorm:"column:id" db:"id"`
	Name  string `dbq:"name" gorm:"column:name" db:"name"`
	Email string `dbq:"email" gorm:"column:email" db:"email"`
}

// Recommended by dbq
func (m *Model) ScanFast() []interface{} {
	return []interface{}{&m.ID, &m.Name, &m.Email}
}

// Required by gorm
func (Model) TableName() string {
	return "tests"
}
