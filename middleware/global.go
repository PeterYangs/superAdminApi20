package middleware

import (
	"github.com/PeterYangs/superAdminCore/contextPlus"
	"github.com/PeterYangs/superAdminCore/middleware/session"
	"superadmin/middleware/accessLog"
)

func Load() []contextPlus.HandlerFunc {

	return []contextPlus.HandlerFunc{

		session.StartSession,
		accessLog.AccessLog,
	}
}
