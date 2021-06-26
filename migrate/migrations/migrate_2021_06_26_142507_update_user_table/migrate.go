package migrate_2021_06_26_142507_update_user_table

import "gin-web/migrate"

func Up() {

	migrate.Table("user", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_06_26_142507_update_user_table"

		createMigrate.String("title", 255)
		createMigrate.String("keyword", 255)

	})

}

func Down() {

}
