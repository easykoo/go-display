package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-display/common"
	"github.com/easykoo/go-display/middleware"
	"github.com/easykoo/go-display/model"

	"encoding/json"
)

func AllUserHandler(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		user := new(model.User)
		user.SetPageActive(true)
		user.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		user.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		user.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		users, total, err := user.SearchByPage()
		PanicIf(err)
		ctx.Set("aaData", users)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "user/allUser", ctx)
	}
}

func DeleteUser(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	user := new(model.User)
	user.Id = ParseInt(id)
	err := user.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeleteUsers(ctx *middleware.Context) {
	users := ctx.R.FormValue("Users")
	var res []int
	json.Unmarshal([]byte(users), &res)
	user := new(model.User)
	err := user.DeleteUsers(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func SetRole(ctx *middleware.Context) {
	id := ctx.R.PostFormValue("Id")
	roleId := ctx.R.PostFormValue("RoleId")
	version := ctx.R.PostFormValue("Version")
	user := new(model.User)
	user.Id = ParseInt(id)
	user.Role.Id = ParseInt(roleId)
	user.Version = ParseInt(version)
	err := user.SetRole()
	PanicIf(err)
	ctx.Set("success", true)
	Log.Debug("User: ", user.Id, " roleId set to ", roleId)
	ctx.JSON(200, ctx.Response)
}

func BanUser(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	user := new(model.User)
	user.Id = ParseInt(id)
	err := user.SetLock(true)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func LiftUser(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	user := new(model.User)
	user.Id = ParseInt(id)
	err := user.SetLock(false)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func PreferencesHandler(ctx *middleware.Context) {
	ctx.HTML(200, "profile/preferences", ctx)
}
