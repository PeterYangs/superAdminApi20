package migrate_2021_08_04_143652_create_role_detail_table

import "gin-web/migrate"

func Up() {

	migrate.Create("role_detail", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_04_143652_create_role_detail_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		// createMigrate.Timestamp("deleted_at").Nullable()

		createMigrate.Integer("admin_id").Unsigned().Comment("管理员id")

		createMigrate.Integer("role_id").Unsigned().Comment("角色id")

		createMigrate.Unique("admin_id", "role_id")

	})

}

func Down() {

	migrate.DropIfExists("role_detail")

}
