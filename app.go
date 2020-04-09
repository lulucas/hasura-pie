package pie

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/sarulabs/di"
	uuid "github.com/satori/go.uuid"
	"github.com/vrischmann/envconfig"
	"net/http"
)

type App struct {
	e         *echo.Echo
	logger    Logger
	modules   []Module
	actions   map[string]Handler
	events    map[string]Handler
	builder   *di.Builder
	container di.Container
	opt       *option
}

func NewApp() *App {
	logger := NewLogger()
	builder, _ := di.NewBuilder()
	opt := &option{}
	if err := envconfig.InitWithPrefix(opt, "APP_"); err != nil {
		logger.Fatalf("Load env error, %s", err.Error())
	}
	return &App{
		e:       echo.New(),
		actions: map[string]Handler{},
		events:  map[string]Handler{},
		logger:  logger,
		builder: builder,
		opt:     opt,
	}
}

type option struct {
	Port int `envconfig:"default=3000"`
}

func (a *App) AddModule(module ...Module) {
	a.modules = append(a.modules, module...)
}

func (a *App) addActionHandler(name string, handler Handler) {
	a.actions[name] = handler
}

func (a *App) addEventHandler(name string, handler Handler) {
	a.events[name] = handler
}

func (a *App) addDef(def ...di.Def) {
	err := a.builder.Add(def...)
	if err != nil {
		a.logger.Fatalf("Add def to builder error, %s", err.Error())
	}
}

func (a *App) Start() {
	// 准备阶段
	for _, m := range a.modules {
		m.BeforeCreated(newModuleContext(a, getModuleName(m)))
	}

	// 构建注入容器
	a.container = a.builder.Build()

	// 完成阶段
	for _, m := range a.modules {
		m.Created(newModuleContext(a, getModuleName(m)))
	}

	a.e.HidePort = true
	a.e.HideBanner = true
	a.e.Use(middleware.CORS())

	// 处理动作
	// 信任来自Action Body Session数据，没有校验机制
	// 因此接口不能暴露在公网
	a.e.POST("/actions", func(c echo.Context) error {
		act := Action{}
		if err := c.Bind(&act); err != nil {
			return c.JSON(http.StatusBadRequest, hasuraErrorResponse(err))
		}
		if h, ok := a.actions[act.Action.Name]; ok {
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, sessionName, Session{
				UserId: uuid.FromStringOrNil(act.SessionVariables.XHasuraUserId),
				Role:   act.SessionVariables.XHasuraRole,
			})
			res, err := h.Invoke(ctx, act.Input)
			if err != nil {
				return c.JSON(http.StatusBadRequest, hasuraErrorResponse(err))
			}
			return c.JSONBlob(http.StatusOK, res)
		} else {
			return errors.Errorf("Action %s not found", act.Action.Name)
		}
	})

	// 处理事件

	// 打印监听端口
	if !IsProduction() {
		a.logger.WithField("core", "http").Warnf("Listen http on http://127.0.0.1:%d", a.opt.Port)
	} else {
		a.logger.WithField("core", "http").Warnf("Listen http on http://0.0.0.0:%d", a.opt.Port)
	}

	// 启动监听
	if err := a.e.Start(fmt.Sprintf(":%d", a.opt.Port)); err != nil {
		a.logger.WithField("core", "http").Fatalf("Listen error %s", err.Error())
	}
}
