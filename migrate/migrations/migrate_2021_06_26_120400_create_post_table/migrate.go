package migrate_2021_06_26_120400_create_post_table

import "gin-web/migrate"

func Up() {

	migrate.Create("post", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_06_26_120400_create_post_table"

		createMigrate.BigIncrements("id")

		createMigrate.Integer("time").Unsigned()

	})

}

func Down() {

	migrate.DropIfExists("post")

}
