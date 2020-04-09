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
	if err := bc.InitConfig(hello{
		World: "good",
	}); err != nil {
		log.Fatal(err)
	}
}

func (m *example) Created(cc pie.CreatedContext) {
	opt := &hello{}
	if err := cc.LoadConfig(opt); err != nil {
		log.Fatalln(err)
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
