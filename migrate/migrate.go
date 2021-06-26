package migrate

import (
	"fmt"
	"gin-web/database"
	"gin-web/model"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

var batch = 1

func init() {

	//加载配置文件
	err := godotenv.Load("./.env")
	if err != nil {
		panic("配置文件加载失败")
	}

	var migrations model.Migrations

	re := database.GetDb().Order("id desc").First(&migrations)

	if re.Error == nil {

		batch = migrations.Batch + 1

		//batch=1
	}

}

type Tag int

const (
	CREATE Tag = 0x00000
	UPDATE Tag = 0x00001
)

type Types string

const (
	Int     Types = "int"
	Varchar Types = "varchar"
)

func (t Types) ToString() string {

	return string(t)
}

type Migrate struct {
	Tag    Tag
	Table  string
	fields []*field
	Name   string
}

type field struct {
	column       string //字段名称
	isPrimaryKey bool   //主键
	isUnsigned   bool   //无符号
	isNullable   bool
	types        Types //数据类型
	length       int
}

func Create(table string, callback func(*Migrate)) {

	m := &Migrate{
		Table: table,
		Tag:   CREATE,
	}

	defer func() {

		run(m)

	}()

	callback(m)

}

func DropIfExists(table string) {

	database.GetDb().Exec("drop table if exists `" + table + "`")

}

// BigIncrements 主键字段
func (c *Migrate) BigIncrements(column string) {

	c.fields = append(c.fields, &field{column: column, isPrimaryKey: true})
}

// Integer int
func (c *Migrate) Integer(column string) *field {

	f := &field{column: column, types: Int, length: 10}

	c.fields = append(c.fields, f)

	return f
}

// Unsigned 无符号
func (f *field) Unsigned() *field {

	f.isUnsigned = true

	return f
}

func (f *field) Nullable() *field {

	f.isNullable = true

	return f
}

func run(m *Migrate) {

	isFind := database.GetDb().Where("migration = ?", m.Name).First(&model.Migrations{})

	//已存在的迁移不执行
	if isFind.Error == nil {

		fmt.Println("find:" + m.Name)

		return
	}

	checkMigrationsTable()

	//batch := 1

	if m.Tag == CREATE {

		sql := "CREATE TABLE `" + m.Table + "` (" +
			"`" + getPrimaryKey(m) + "` int(10) unsigned NOT NULL AUTO_INCREMENT," +
			getColumn(m) +
			"PRIMARY KEY (`" + getPrimaryKey(m) + "`)" +
			") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"

		t := database.GetDb().Exec(sql)

		fmt.Println(t.Error)

		database.GetDb().Create(&model.Migrations{
			Migration: m.Name,
			Batch:     batch,
		})

	}

}

func getPrimaryKey(m *Migrate) string {

	id := ""

	for _, f := range m.fields {

		if f.isPrimaryKey {

			id = f.column
		}
	}

	if id == "" {

		panic("主键不能为空")
	}

	return id
}

func getColumn(m *Migrate) string {

	str := ""

	for _, f := range m.fields {

		if f.isPrimaryKey {

			continue

		}

		str += "`" + f.column + "` " + f.types.ToString() + "(" + cast.ToString(f.length) + ") "

		if f.isUnsigned {

			str += " unsigned "
		}

		if !f.isNullable {

			str += " NOT NULL "
		}

		str += ","

	}

	return str
}

// CheckMigrationsTable 检查数据迁移表是否存在
func checkMigrationsTable() {

	database.GetDb().Exec("CREATE TABLE IF NOT EXISTS `migrations` (`id` int(10) unsigned NOT NULL AUTO_INCREMENT,  `migration` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,  `batch` int(11) NOT NULL,  PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

}
