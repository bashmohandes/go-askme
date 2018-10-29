package middleware

import "github.com/bashmohandes/go-askme/web/framework"

type authMiddleware struct {
}

func (a *authMiddleware) Run(cxt framework.Context) bool {
	return true
}
