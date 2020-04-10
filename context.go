package pie

import (
	"context"
	"encoding/json"
	"github.com/Rican7/conjson"
	"github.com/Rican7/conjson/transform"
	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sarulabs/di"
	uuid "github.com/satori/go.uuid"
	"github.com/vrischmann/envconfig"
	"reflect"
)

const (
	sessionName = "_session"
)

type BeforeCreatedContext interface {
	DB() *pg.DB
	LoadFromEnv(opt interface{})
	InitConfig(key string, opt interface{})
	LoadConfig(key string, opt interface{}) error
	SaveConfig(key string, opt interface{}) error
	Add(def ...di.Def)
}

type CreatedContext interface {
	DB() *pg.DB
	GetSession(ctx context.Context) Session
	Get(name string) interface{}
	Logger() Logger
	Rest() *echo.Group
	HandleAction(name string, handler interface{})
	HandleEvent(name string, handler interface{})
	LoadConfig(key string, opt interface{}) error
	SaveConfig(key string, opt interface{}) error
}

type Session struct {
	UserId uuid.UUID
	Role   string
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

func (c *moduleContext) DB() *pg.DB {
	return c.app.db
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

func (c *moduleContext) InitConfig(key string, opt interface{}) {
	marshaller := conjson.NewMarshaler(opt, transform.ConventionalKeys())
	data, err := json.Marshal(marshaller)
	if err != nil {
		c.logger.Fatalf("Init config error, %s", err.Error())
	}

	cfg := Config{
		Key:  key,
		Data: data,
	}
	if _, err := c.app.db.Model(&cfg).OnConflict("DO NOTHING").SelectOrInsert(); err != nil {
		c.logger.Fatalf("Init config error, %s", err.Error())
	}
}

func (c *moduleContext) LoadConfig(key string, opt interface{}) error {
	t := reflect.TypeOf(opt)
	if t.Kind() != reflect.Ptr {
		return errors.New("opt must be pointer")
	}

	cfg := Config{}
	if err := c.app.db.Model(&cfg).Where("key = ?", key).Select(); err != nil {
		if err == pg.ErrNoRows {
			return errors.Errorf("config key %s not found", key)
		}
		return err
	}
	if err := json.Unmarshal(cfg.Data, conjson.NewUnmarshaler(opt, transform.ConventionalKeys())); err != nil {
		return err
	}
	return nil
}

func (c *moduleContext) SaveConfig(key string, opt interface{}) error {
	marshaller := conjson.NewMarshaler(opt, transform.ConventionalKeys())
	data, err := json.Marshal(marshaller)
	if err != nil {
		return err
	}

	cfg := Config{
		Key:  key,
		Data: data,
	}
	if _, err := c.app.db.Model(&cfg).Where("key = ?", key).Update(); err != nil {
		return err
	}
	return nil
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

func (c *moduleContext) Rest() *echo.Group {
	return c.app.externalEcho.Group("/" + c.module)
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
	c.logger.Infof("Event handler is added: %s", name)
	c.app.addEventHandler(name, NewHandler(handler))
}
