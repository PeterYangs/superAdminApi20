package migrate_2021_08_18_155238_update_access_table

import "github.com/PeterYangs/superAdminCore/v2/migrate"

func Up() {

	migrate.Table("access", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_18_155238_update_access_table"

		createMigrate.Integer("admin_id").Unsigned().Default(0).Comment("管理员id")

	})

}

func Down() {

	migrate.Table("access", func(createMigrate *migrate.Migrate) {

	})

}
