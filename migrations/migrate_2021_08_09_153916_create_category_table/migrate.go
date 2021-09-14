package migrate_2021_08_09_153916_create_category_table

import "github.com/PeterYangs/superAdminCore/migrate"

func Up() {

	migrate.Create("category", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_09_153916_create_category_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		// createMigrate.Timestamp("deleted_at").Nullable()

		createMigrate.Integer("pid").Unsigned().Comment("父类id")

		createMigrate.Integer("lv").Unsigned().Default(1).Comment("层级")

		createMigrate.String("title", 255).Comment("标题")

		createMigrate.String("img", 255).Default("").Comment("图片")

		createMigrate.Integer("sort").Comment("排序,越小越靠前")

		createMigrate.String("path", 255).Comment("层级路径，逗号分隔")

		createMigrate.Unique("title", "pid")

	})

}

func Down() {

	migrate.DropIfExists("category")

}
