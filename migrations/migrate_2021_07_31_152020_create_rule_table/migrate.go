package migrate_2021_07_31_152020_create_rule_table

import "github.com/PeterYangs/superAdminCore/v2/migrate"

func Up() {

	migrate.Create("rule", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_07_31_152020_create_rule_table"

		createMigrate.BigIncrements("id")

		createMigrate.Timestamp("created_at").Nullable()

		createMigrate.Timestamp("updated_at").Nullable()

		createMigrate.Timestamp("deleted_at").Nullable()

		createMigrate.String("title", 255).Comment("规则描述")

		createMigrate.String("rule", 255).Comment("规则")

		createMigrate.String("group_name", 255).Comment("分组名称")

		//createMigrate.String("test", 255).Default(migrate.Null).Comment("测试")

	})

}

func Down() {

	migrate.DropIfExists("rule")

}
