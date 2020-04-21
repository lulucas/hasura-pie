package pie

import (
	"context"
	"encoding/json"
	"fmt"
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
	InitConfig(cfg interface{})
	LoadConfig(cfg interface{}) error
	SaveConfig(cfg interface{}) error
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
	LoadConfig(cfg interface{}) error
	SaveConfig(cfg interface{}) error
	HandleCron(name, spec string, cmd func())
}

type Session struct {
	UserId *uuid.UUID
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

func (c *moduleContext) InitConfig(cfg interface{}) {
	t := reflect.TypeOf(cfg)
	if t.Kind() != reflect.Ptr {
		c.logger.Fatal("cfg must be pointer")
	}

	marshaller := conjson.NewMarshaler(cfg, transform.ConventionalKeys())
	data, err := json.Marshal(marshaller)
	if err != nil {
		c.logger.Fatalf("Init config error, %s", err.Error())
	}

	key := fmt.Sprintf("%s.%s", c.module, t.Elem().Name())

	config := Config{
		Key:  key,
		Data: data,
	}
	if _, err := c.app.db.Model(&config).Where("key = ?", key).Limit(1).OnConflict("DO NOTHING").SelectOrInsert(); err != nil {
		c.logger.Fatalf("Init config error, %s", err.Error())
	}
}

func (c *moduleContext) LoadConfig(cfg interface{}) error {
	t := reflect.TypeOf(cfg)
	if t.Kind() != reflect.Ptr {
		return errors.New("cfg must be pointer")
	}

	key := fmt.Sprintf("%s.%s", c.module, t.Elem().Name())

	config := Config{}
	if err := c.app.db.Model(&config).Where("key = ?", key).Select(); err != nil {
		if err == pg.ErrNoRows {
			return errors.Errorf("config key %s not found", key)
		}
		return err
	}
	if err := json.Unmarshal(config.Data, conjson.NewUnmarshaler(cfg, transform.ConventionalKeys())); err != nil {
		return err
	}
	return nil
}

func (c *moduleContext) SaveConfig(cfg interface{}) error {
	t := reflect.TypeOf(cfg)
	if t.Kind() != reflect.Ptr {
		return errors.New("cfg must be pointer")
	}

	marshaller := conjson.NewMarshaler(cfg, transform.ConventionalKeys())
	data, err := json.Marshal(marshaller)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s.%s", c.module, t.Elem().Name())

	config := Config{
		Key:  key,
		Data: data,
	}
	if _, err := c.app.db.Model(&config).Where("key = ?", key).Update(); err != nil {
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

// Scheduler timer
func (c *moduleContext) HandleCron(name, spec string, cmd func()) {
	c.logger.Infof("Cron func is added: %s on %s", name, spec)
	if _, err := c.app.cron.AddFunc(spec, cmd); err != nil {
		c.logger.Fatalf("Add cron func error, %s", err.Error())
	}
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
