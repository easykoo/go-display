package handler

import (
	"github.com/go-martini/martini"

	"github.com/easykoo/go-display/middleware"
)

func Index(ctx *middleware.Context) {
	ctx.Redirect("/admin/welcome")
}

func LangHandler(ctx *middleware.Context, params martini.Params) {
	lang := params["lang"]
	ctx.S.Set("Lang", lang)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}
