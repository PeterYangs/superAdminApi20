package artisan

import (
	"github.com/PeterYangs/superAdminCore/v2/artisan"
	"superadmin/artisan/seeds"
)

func Load() []artisan.Artisan {

	return []artisan.Artisan{
		new(seeds.Seeds),
	}
}
