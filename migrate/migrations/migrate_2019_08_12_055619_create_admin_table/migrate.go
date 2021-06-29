package migrate_2019_08_12_055619_create_admin_table

import "gin-web/migrate"

func Up() {

	migrate.Create("admin", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2019_08_12_055619_create_admin_table"

		//主键
		createMigrate.BigIncrements("id")

		//int
		createMigrate.Integer("user_id").Unsigned().Nullable().Default(0).Unique().Comment("用户id")

		//varchar
		createMigrate.String("title", 255).Default("").Unsigned().Comment("标题")

		//text
		createMigrate.Text("content").Comment("内容")

		//索引
		createMigrate.Unique("user_id", "title")

	})

}

func Down() {

	migrate.DropIfExists("admin")

}
