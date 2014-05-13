package handler

import (
	. "github.com/easykoo/go-display/common"
	"github.com/easykoo/go-display/middleware"
	"github.com/easykoo/go-display/model"
)

func WelcomeHandler(ctx *middleware.Context) {
	ctx.HTML(200, "admin/welcome", ctx)
}

func SettingsHandler(ctx *middleware.Context, settings model.Settings) {
	if ctx.R.Method == "POST" {
		err := settings.Update()
		PanicIf(err)
		dbSettings := model.GetSettings()
		ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.change.success"))
		ctx.S.Set("Settings", dbSettings)
	}
	user := &model.User{}
	users, err := user.SelectAll()
	PanicIf(err)
	ctx.Set("Users", users)

	ctx.HTML(200, "admin/settings", ctx)
}
