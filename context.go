package pie

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	uuid "github.com/satori/go.uuid"
	"github.com/vrischmann/envconfig"
)

const (
	sessionName = "_session"
)

type BeforeCreatedContext interface {
	// 加载环境变量
	LoadFromEnv(opt interface{})
	// 添加注入对象
	Add(def ...di.Def)
}

type CreatedContext interface {
	// 获取Session
	GetSession(ctx context.Context) Session
	// 获取注入对象
	Get(name string) interface{}
	// 日志
	Logger() Logger
	// Http
	Http() *echo.Echo
	// 处理动作
	HandleAction(name string, handler interface{})
	// 处理事件
	HandleEvent(name string, handler interface{})
}

type Session struct {
	// 用户编号
	UserId uuid.UUID
	// 用户角色
	Role string
}

type moduleContext struct {
	app    *App
	module string
	logger Logger
}

func newModuleContext(app *App, module string) *moduleContext {
	return &moduleContext{
		app:    app,
		module: module,
		logger: app.logger.WithField("module", module),
	}
}

func (c *moduleContext) GetSession(ctx context.Context) Session {
	return ctx.Value(sessionName).(Session)
}

func (c *moduleContext) LoadFromEnv(opt interface{}) {
	if err := envconfig.InitWithPrefix(opt, c.module); err != nil {
		c.logger.Fatalf("Load env error, %s", err.Error())
	}
	c.logger.Debugf("Option value after loaded from env, %+v", opt)
}

func (c *moduleContext) Add(def ...di.Def) {
	c.app.addDef(def...)
}

func (c *moduleContext) Get(name string) interface{} {
	return c.app.container.Get(name)
}

func (c *moduleContext) Logger() Logger {
	return c.logger
}

func (c *moduleContext) Http() *echo.Echo {
	return c.app.externalEcho
}

// Timer task TODO
func (c *moduleContext) HandleCron() {

}

// Hasura action
func (c *moduleContext) HandleAction(name string, handler interface{}) {
	c.logger.Infof("Action handler is added: %s", name)
	c.app.addActionHandler(name, NewHandler(handler))
}

// Hasura event
func (c *moduleContext) HandleEvent(name string, handler interface{}) {
	c.logger.Infof("RawEvent handler is added: %s", name)
	c.app.addEventHandler(name, NewHandler(handler))
}
