package migrate_2019_08_12_666666_create_user_table

import "gin-web/migrate"

func Up() {

	migrate.Create("user", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2019_08_12_666666_create_user_table"

		createMigrate.BigIncrements("id")

		createMigrate.Integer("name").Unsigned().Nullable()

	})

}

func Down() {

	migrate.DropIfExists("user")

}
