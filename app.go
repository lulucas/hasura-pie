package pie

import (
	"context"
	"encoding/json"
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
	internalEcho *echo.Echo
	externalEcho *echo.Echo
	logger       Logger
	modules      []Module
	actions      map[string]Handler
	events       map[string]Handler
	builder      *di.Builder
	container    di.Container
	opt          *option
}

func NewApp() *App {
	logger := NewLogger()
	builder, _ := di.NewBuilder()
	opt := &option{}
	if err := envconfig.InitWithPrefix(opt, "APP_"); err != nil {
		logger.Fatalf("Load env error, %s", err.Error())
	}
	return &App{
		internalEcho: echo.New(),
		actions:      map[string]Handler{},
		events:       map[string]Handler{},
		logger:       logger,
		builder:      builder,
		opt:          opt,
	}
}

type option struct {
	InternalPort int `envconfig:"default=3000"`
	ExternalPort int `envconfig:"default=8000"`
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

	a.internalEcho.HidePort = true
	a.internalEcho.HideBanner = true
	a.internalEcho.Use(middleware.CORS())

	// DO NOT EXPOSE INTERNAL ECHO PORT TO PUBLIC NETWORK !!!
	// Handle hasura actions
	a.internalEcho.POST("/actions", func(c echo.Context) error {
		act := Action{}
		if err := c.Bind(&act); err != nil {
			return c.JSON(http.StatusBadRequest, hasuraErrorResponse(err))
		}
		if h, ok := a.actions[act.Action.Name]; ok {
			a.logger.WithField("core", "action").Infof("Call %s", act.Action.Name)
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

	// Handle hasura events
	a.internalEcho.POST("/events", func(c echo.Context) error {
		rawEvt := RawEvent{}
		if err := c.Bind(&rawEvt); err != nil {
			return c.JSON(http.StatusBadRequest, hasuraErrorResponse(err))
		}
		if h, ok := a.events[rawEvt.Trigger.Name]; ok {
			a.logger.WithField("core", "event").Infof("Trigger %s, %s", rawEvt.Trigger.Name, rawEvt.Id)
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, sessionName, Session{
				UserId: uuid.FromStringOrNil(rawEvt.Event.SessionVariables.XHasuraUserID),
				Role:   rawEvt.Event.SessionVariables.XHasuraRole,
			})
			evt := &Event{
				Id:        rawEvt.Id,
				CreatedAt: rawEvt.CreatedAt,
				Table:     rawEvt.Table,
				Op:        rawEvt.Event.Op,
				Old:       rawEvt.Event.Data.Old,
				New:       rawEvt.Event.Data.New,
			}
			payload, err := json.Marshal(evt)
			if err != nil {
				return err
			}
			res, err := h.Invoke(ctx, payload)
			if err != nil {
				return c.JSON(http.StatusBadRequest, hasuraErrorResponse(err))
			}
			return c.JSONBlob(http.StatusOK, res)
		} else {
			return errors.Errorf("RawEvent %s not found", rawEvt.Trigger.Name)
		}
	})

	// 打印监听端口
	if !IsProduction() {
		a.logger.WithField("core", "http").Warnf("Listen internal http on http://127.0.0.1:%d", a.opt.InternalPort)
	} else {
		a.logger.WithField("core", "http").Warnf("Listen internal http on http://0.0.0.0:%d", a.opt.InternalPort)
	}

	// 启动监听
	if err := a.internalEcho.Start(fmt.Sprintf(":%d", a.opt.InternalPort)); err != nil {
		a.logger.WithField("core", "http").Fatalf("Listen error %s", err.Error())
	}
}
