package migrate_2021_08_17_165541_update_menu_table

import "github.com/PeterYangs/superAdminCore/v2/migrate"

func Up() {

	migrate.Table("menu", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_17_165541_update_menu_table"

		createMigrate.String("rule", 255).Default("").Comment("对应的权限，用于显示和隐藏不同角色的菜单显示")

	})

}

func Down() {

	migrate.Table("menu", func(createMigrate *migrate.Migrate) {

		createMigrate.DropColumn("rule")

	})

}
