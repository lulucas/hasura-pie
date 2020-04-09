package main

import (
	"fmt"
	pie "github.com/lulucas/hasura-pie"
	"log"
)

type example struct {
}

type hello struct {
	World string
}

func (m *example) BeforeCreated(bc pie.BeforeCreatedContext) {
	bc.InitConfig("hello", hello{
		World: "good",
	})
}

func (m *example) Created(cc pie.CreatedContext) {
	opt := &hello{}
	if err := cc.LoadConfig("hello", opt); err != nil {
		log.Fatal(err)
	}
	fmt.Println(opt.World)
}

func main() {
	app := pie.NewApp()
	app.AddModule(
		&example{},
	)
	app.Start()
}
