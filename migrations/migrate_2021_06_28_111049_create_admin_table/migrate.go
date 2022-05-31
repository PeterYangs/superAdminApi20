package migrate_2021_06_28_111049_create_admin_table

import "github.com/PeterYangs/superAdminCore/v2/migrate"

func Up() {

	migrate.Create("admin", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_06_28_111049_create_admin_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		createMigrate.String("username", 255).Unique().Comment("用户名")

		createMigrate.String("password", 255).Comment("密码")

		createMigrate.String("email", 255).Unique().Comment("邮箱")

		createMigrate.Integer("status").Default(1).Comment("状态")

		createMigrate.Timestamp("deleted_at").Nullable()

	})

}

func Down() {

	migrate.DropIfExists("admin")

}
