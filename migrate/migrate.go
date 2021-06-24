package migrate

type migrate struct {
}

type create struct {
	Table string
}

func (m migrate) Create(table string, callback func(create)) {

	callback(create{
		Table: table,
	})
}

// BigIncrements 主键字段
func (c create) BigIncrements(id string) {

}
