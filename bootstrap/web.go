package bootstrap

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/websocket"
	"time"
)

type Configurator func(*Bootstrapper)

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time

	Sessions *sessions.Sessions
}

func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b)
	}

	return b
}

func (b *Bootstrapper) SetupViews(viewsDir string) {
	b.RegisterView(iris.HTML(viewsDir, ".html").Layout("shared/layout.html"))
}

func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName,
		Expires:  expires,
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

func (b *Bootstrapper) SetupWebsockets(endpoint string, handler websocket.ConnHandler) {
	ws := websocket.New(websocket.DefaultGorillaUpgrader, handler)
	b.Get(endpoint, websocket.Handler(ws))
}

func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		ctx.JSON(err)
		return
	})
}

const (
	StaticAssets = "public/"
	Favicon      = "favicon.ico"
)

func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	//b.SetupSessions(30*24*time.Hour,
	//	[]byte("QXRraUkPr7459Te00aQ6JCiToueeeI83zY8jn8Np7O2mHrGmWrkJgkdOiDZ0hVrv"),
	//	[]byte("Z6Q34tCOiDNkHIUxObZXNgmH2QQhcmynmHUHenLawOqm8YNPYt9MC1RZW0YYslYG"),
	//)

	b.SetupViews("app/web/views")
	b.SetupErrorHandlers()
	b.Favicon(StaticAssets + Favicon)
	b.HandleDir("/public", iris.Dir(StaticAssets))

	b.Use(recover.New())
	b.Use(logger.New())

	return b
}

func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}
