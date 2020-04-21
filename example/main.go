package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	pie "github.com/lulucas/hasura-pie"
	"log"
	"net/http"
)

type example struct {
}

type hello struct {
	World string
}

func (m *example) BeforeCreated(bc pie.BeforeCreatedContext) {
	bc.InitConfig(&hello{
		World: "good",
	})
}

func (m *example) Created(cc pie.CreatedContext) {
	opt := &hello{}
	if err := cc.LoadConfig(opt); err != nil {
		log.Fatal(err)
	}
	fmt.Println(opt.World)

	cc.Rest().POST("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "world")
	})
}

func main() {
	app := pie.NewApp()
	app.AddModule(
		&example{},
	)
	app.Start()
}
