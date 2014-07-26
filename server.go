package main

import (
	"github.com/easykoo/binding"
	"github.com/easykoo/sessions"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"

	. "github.com/easykoo/go-display/auth"
	. "github.com/easykoo/go-display/common"
	"github.com/easykoo/go-display/handler"
	"github.com/easykoo/go-display/middleware"
	"github.com/easykoo/go-display/model"

	"encoding/gob"
	"html/template"
	"os"
	"time"
)

func init() {
	SetConfig()
	SetLog()
	Log.Debug("server initializing...")
}

func newMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(middleware.GetLogger())
	m.Map(model.SetEngine())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	gob.Register(model.User{})
	gob.Register(model.Settings{})

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))

	m.Use(render.Renderer(render.Options{
		Directory:  "view",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		Funcs: []template.FuncMap{
			{
				"formatTime": func(args ...interface{}) string {
					return args[0].(time.Time).Format("Jan _2 15:04")
				},
				"cnFormatTime": func(args ...interface{}) string {
					return args[0].(time.Time).Format("2006-01-02 15:04")
				},
				"mdToHtml": func(args ...interface{}) template.HTML {
					return template.HTML(string(blackfriday.MarkdownBasic([]byte(args[0].(string)))))
				},
				"unescaped": func(args ...interface{}) template.HTML {
					return template.HTML(args[0].(string))
				},
				"equal": func(args ...interface{}) bool {
					return args[0] == args[1]
				},
				"tsl": func(lang string, format string) string {
					return Translate(lang, format)
				},
				"tslf": func(lang string, format string, args ...interface{}) string {
					return Translatef(lang, format, args...)
				},
				"privilege": func(user interface{}, module int) bool {
					if user == nil {
						return false
					}
					return CheckPermission(user, module)
				},
				"plus": func(args ...int) int {
					var result int
					for _, val := range args {
						result += val
					}
					return result
				},
			},
		},
	}))

	m.Use(middleware.InitContext())

	return &martini.ClassicMartini{m, r}
}

func main() {
	m := newMartini()

	m.Get("/", handler.Index)
	m.Get("/index", handler.Index)
	m.Get("/language/change/:lang", handler.LangHandler)

	m.Group("/user", func(r martini.Router) {
		r.Any("/all", AuthRequest(Module_User), handler.AllUserHandler)
		r.Any("/logout", handler.LogoutHandler)
		r.Any("/login", binding.Form(model.UserLoginForm{}), handler.LoginHandler)
		r.Any("/register", binding.Form(model.UserRegisterForm{}), handler.RegisterHandler)
		r.Any("/delete", AuthRequest(Module_User), handler.DeleteUsers)
		r.Any("/delete/:id", AuthRequest(Module_User), handler.DeleteUser)
		r.Any("/role", AuthRequest(Module_User), handler.SetRole)
		r.Any("/ban/:id", AuthRequest(Module_User), handler.BanUser)
		r.Any("/lift/:id", AuthRequest(Module_User), handler.LiftUser)
		r.Any("/loginx", binding.Form(model.UserLoginForm{}), handler.LoginxHandler)
	})

	m.Group("/profile", func(r martini.Router) {
		r.Any("/profile", AuthRequest(SignInRequired), binding.Form(model.User{}), handler.ProfileHandler)
		r.Any("/preferences", AuthRequest(SignInRequired), handler.PreferencesHandler)
		r.Any("/password", AuthRequest(SignInRequired), binding.Form(model.Password{}), handler.PasswordHandler)
		r.Any("/checkEmail", AuthRequest(SignInRequired), binding.Form(model.User{}), handler.CheckEmail)
	})

	m.Group("/admin", func(r martini.Router) {
		r.Get("/welcome", AuthRequest(SignInRequired), handler.WelcomeHandler)
		r.Any("/settings", AuthRequest(Module_Admin), binding.Form(model.Settings{}), handler.SettingsHandler)
	})

	m.Group("/picture", func(r martini.Router) {
		r.Any("", AuthRequest(Module_Pad), handler.Picture)
		r.Any("/upload", AuthRequest(Module_Pad), handler.Upload)
		r.Any("/new", AuthRequest(Module_Pad), binding.Bind(model.Picture{}), handler.NewPicture)
		r.Any("/delete", AuthRequest(Module_Pad), handler.DeletePictures)
		r.Any("/delete/:id", AuthRequest(Module_Pad), handler.DeletePicture)
		r.Any("/edit/:id", AuthRequest(Module_Pad), handler.EditPictureView)
		r.Any("/edit", AuthRequest(Module_Pad), binding.Bind(model.Picture{}), handler.EditPicture)
	})

	m.Group("/pad", func(r martini.Router) {
		r.Any("", AuthRequest(Module_Picture), handler.Pad)
		r.Any("/new", AuthRequest(Module_Picture), binding.Bind(model.Pad{}), handler.NewPad)
		r.Any("/delete", AuthRequest(Module_Picture), handler.DeletePads)
		r.Any("/delete/:id", AuthRequest(Module_Picture), handler.DeletePad)
		r.Any("/edit/:id", AuthRequest(Module_Picture), handler.EditPadView)
		r.Any("/edit", AuthRequest(Module_Picture), binding.Bind(model.Pad{}), handler.EditPad)
		r.Any("/choosePicture", AuthRequest(Module_Picture), handler.ChoosePicture)
		r.Any("/:name/info", handler.Info)
	})

	Log.Info("server is started...")
	os.Setenv("PORT", Cfg.MustValue("", "http_port", "3000"))
	m.Run()
}
