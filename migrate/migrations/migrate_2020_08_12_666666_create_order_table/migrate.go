package migrate_2020_08_12_666666_create_order_table

import "gin-web/migrate"

func Up() {

	migrate.Create("order", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2020_08_12_666666_create_order_table"

		createMigrate.BigIncrements("id")

		createMigrate.Integer("name").Unsigned().Nullable()

	})

}

func Down() {

	migrate.DropIfExists("order")

}
