package routes

import (
	"github.com/PeterYangs/superAdminCore/route"
	"superadmin/controller"
)

func Routes(r route.Group) {

	r.Registered(route.GET, "/index", controller.Index)

}
