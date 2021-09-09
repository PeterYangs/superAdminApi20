package migrate_2021_09_09_094754_create_file_table

import "gin-web/migrate"

func Up() {

	migrate.Create("file", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_09_09_094754_create_file_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		// createMigrate.Timestamp("deleted_at").Nullable()

		createMigrate.String("path", 255).Comment("文件地址")

		createMigrate.String("name", 255).Comment("文件名称")

		createMigrate.BigInteger("size").Unsigned().Comment("文件大小，单位字节")

		createMigrate.Integer("admin_id").Unsigned().Comment("上传管理员")

	})

}

func Down() {

	migrate.DropIfExists("file")

}
