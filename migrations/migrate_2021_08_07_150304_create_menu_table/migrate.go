package migrate_2021_08_07_150304_create_menu_table

import "github.com/PeterYangs/superAdminCore/migrate"

func Up() {

	migrate.Create("menu", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_07_150304_create_menu_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		// createMigrate.Timestamp("deleted_at").Nullable()

		createMigrate.Integer("pid").Unsigned().Comment("父级菜单id")

		createMigrate.String("title", 255).Comment("菜单名称")

		createMigrate.String("path", 255).Comment("路径")

		createMigrate.Integer("sort").Unsigned().Default(100).Comment("排序,越小越靠前")

	})

}

func Down() {

	migrate.DropIfExists("menu")

}
