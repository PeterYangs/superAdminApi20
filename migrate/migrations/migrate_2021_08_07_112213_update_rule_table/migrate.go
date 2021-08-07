package migrate_2021_08_07_112213_update_rule_table

import "gin-web/migrate"

func Up() {

	migrate.Table("rule", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_08_07_112213_update_rule_table"

		createMigrate.Unique("rule")
	})

}

func Down() {

}
