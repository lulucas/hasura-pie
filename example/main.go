package main

import pie "github.com/lulucas/hasura-pie"

func main() {
	app := pie.NewApp()
	app.AddModule()
	app.Start()
}
