package handler

import (
	"github.com/easykoo/binding"

	. "github.com/easykoo/go-display/common"
	"github.com/easykoo/go-display/middleware"
	"github.com/easykoo/go-display/model"
)

func LogoutHandler(ctx *middleware.Context) {
	ctx.S.Set("SignedUser", nil)
	ctx.Redirect("/index")
}

func LoginHandler(ctx *middleware.Context, formErr binding.Errors, loginUser model.UserLoginForm) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		password := Md5(loginUser.Password)
		user := &model.User{Username: loginUser.Username, Password: password}
		if !ctx.HasError() {
			if has, err := user.Exist(); has {
				PanicIf(err)
				if user.Locked {
					ctx.Set("User", user)
					ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.invalid.username.or.password"))
					ctx.HTML(200, "user/login", ctx)
					return
				}
				ctx.S.Set("SignedUser", user)
				Log.Info(user.Username, " login")
				ctx.Redirect("/admin/welcome")
			} else {
				ctx.Set("User", user)
				ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.invalid.username.or.password"))
				ctx.HTML(200, "user/login", ctx)
			}
		} else {
			ctx.HTML(200, "user/login", ctx)
		}
	default:
		ctx.HTML(200, "user/login", ctx)
	}
}

func RegisterHandler(ctx *middleware.Context, formErr binding.Errors, user model.UserRegisterForm) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		if !ctx.HasError() {
			dbUser := model.User{Username: user.Username, Password: user.Password, Email: user.Email}

			if exist, err := dbUser.ExistUsername(); exist {
				PanicIf(err)
				ctx.AddFieldError("username", Translate(ctx.S.Get("Lang").(string), "message.error.already.exists"))
			}

			if exist, err := dbUser.ExistEmail(); exist {
				PanicIf(err)
				ctx.AddFieldError("email", Translate(ctx.S.Get("Lang").(string), "message.error.already.exists"))
			}

			if !ctx.HasError() {
				dbUser.Password = Md5(user.Password)
				err := dbUser.Insert()
				PanicIf(err)
				ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.register.success"))
			} else {
				ctx.Set("User", user)
			}
			ctx.HTML(200, "user/register", ctx)
		} else {
			ctx.Set("User", user)
			ctx.HTML(200, "user/register", ctx)
		}
	default:
		ctx.HTML(200, "user/register", ctx)
	}
}

func ProfileHandler(ctx *middleware.Context, formErr binding.Errors, user model.User) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		if !ctx.HasError() {
			err := user.Update()
			PanicIf(err)
			dbUser, err := user.GetUserById(user.Id)
			PanicIf(err)
			ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.change.success"))
			ctx.S.Set("SignedUser", dbUser)
		}
		ctx.HTML(200, "profile/profile", ctx)
	default:
		ctx.HTML(200, "profile/profile", ctx)
	}
}

func PasswordHandler(ctx *middleware.Context, formErr binding.Errors, password model.Password) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		if !ctx.HasError() {
			if password.CurrentPassword == password.ConfirmPassword {
				ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.password.not.changed"))
			} else {
				user := &model.User{Id: password.Id}
				dbUser, err := user.GetUserById(user.Id)
				PanicIf(err)
				if dbUser.Password == Md5(password.CurrentPassword) {
					dbUser.Password = Md5(password.ConfirmPassword)
					err := dbUser.Update()
					PanicIf(err)
					ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.change.success"))
				} else {
					ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.wrong.password"))
				}
			}
		}
	default:
	}
	ctx.HTML(200, "profile/password", ctx)
}

func CheckEmail(ctx *middleware.Context) {
	if user := ctx.S.Get("SignedUser"); user.(model.User).Email != ctx.R.Form["email"][0] {
		test := &model.User{Email: ctx.R.Form["email"][0]}
		if exist, _ := test.ExistEmail(); exist {
			ctx.JSON(200, Translate(ctx.S.Get("Lang").(string), "message.error.already.exists"))
			return
		}
	}
	ctx.JSON(200, true)
}
