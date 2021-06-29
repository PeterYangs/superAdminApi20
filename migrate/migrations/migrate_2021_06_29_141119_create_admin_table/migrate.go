package migrate_2021_06_29_141119_create_admin_table

import "gin-web/migrate"

func Up() {

	migrate.Create("admin", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_06_29_141119_create_admin_table"

		createMigrate.BigIncrements("id")

		createMigrate.String("name", 255).Comment("姓名")

	})

}

func Down() {

	migrate.DropIfExists("admin")

}
