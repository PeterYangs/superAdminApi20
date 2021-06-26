package migrate_2021_06_26_111540_create_task_table

import "gin-web/migrate"

func Up() {

	migrate.Create("task", func(createMigrate *migrate.Migrate) {

		createMigrate.Name = "migrate_2021_06_26_111540_create_task_table"

		createMigrate.BigIncrements("id")

	})

}

func Down() {

	migrate.DropIfExists("task")

}
