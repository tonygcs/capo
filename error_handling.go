package capo

import (
	"net/http"
)

func ErrorHandling(ctx *Context) {
	ctxErr := ctx.Err()
	if ctxErr != nil {
		switch err := ctxErr.(type) {
		case *ServerError:
			ctx.Write(err)
		default:
			ctx.SetStatus(http.StatusInternalServerError)
			ctx.Write(NewServerError(InternalServerErrorCode, err))
		}
	}
}
