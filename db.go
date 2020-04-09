package pie

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/vrischmann/envconfig"
)

type dbOption struct {
	Host     string `envconfig:"default=127.0.0.1"`
	Port     int    `envconfig:"default=5432"`
	User     string `envconfig:"default=postgres"`
	Password string `envconfig:"optional"`
	Database string `envconfig:"default=postgres"`
}

func newDB(logger Logger) *pg.DB {
	opt := &dbOption{}
	if err := envconfig.InitWithPrefix(opt, "DB"); err != nil {
		logger.Fatalf("Load env error, %s", err.Error())
	}
	return pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		User:     opt.User,
		Password: opt.Password,
		Database: opt.Database,
	})
}
