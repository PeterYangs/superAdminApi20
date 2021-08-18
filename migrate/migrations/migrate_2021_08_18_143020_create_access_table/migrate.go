package migrate_2021_08_18_143020_create_access_table

import "gin-web/migrate"

func Up() {

	migrate.Create("access", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_18_143020_create_access_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		createMigrate.String("ip", 255).Comment("ip")

		createMigrate.Text("url").Comment("请求路径")

		createMigrate.Text("params").Comment("参数")

		// createMigrate.Timestamp("deleted_at").Nullable()

	})

}

func Down() {

	migrate.DropIfExists("access")

}
