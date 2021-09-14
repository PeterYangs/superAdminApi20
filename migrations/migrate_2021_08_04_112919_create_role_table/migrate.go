package migrate_2021_08_04_112919_create_role_table

import "github.com/PeterYangs/superAdminCore/migrate"

func Up() {

	migrate.Create("role", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_04_112919_create_role_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		// createMigrate.Timestamp("deleted_at").Nullable()

		createMigrate.String("title", 255).Comment("描述")

		createMigrate.String("rules", 1000).Comment("规则id")

	})

}

func Down() {

	migrate.DropIfExists("role")

}
