package middleware

import (
	"github.com/PeterYangs/superAdminCore/v2/contextPlus"
	"github.com/PeterYangs/superAdminCore/v2/middleware/session"
	"superadmin/middleware/accessLog"
)

func Load() []contextPlus.HandlerFunc {

	return []contextPlus.HandlerFunc{

		session.StartSession,
		accessLog.AccessLog,
	}
}
