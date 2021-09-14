package artisan

import (
	"github.com/PeterYangs/superAdminCore/artisan"
	"superadmin/artisan/seeds"
)

func Load() []artisan.Artisan {

	return []artisan.Artisan{
		new(seeds.Seeds),
	}
}
