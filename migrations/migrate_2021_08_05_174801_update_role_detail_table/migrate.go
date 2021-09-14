package migrate_2021_08_05_174801_update_role_detail_table

import "github.com/PeterYangs/superAdminCore/migrate"

func Up() {

	migrate.Table("role_detail", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_05_174801_update_role_detail_table"

		createMigrate.Integer("role_id").Comment("角色id,-1是超级管理员").Change()

	})

}

func Down() {

}
