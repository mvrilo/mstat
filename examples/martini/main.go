package main

import (
	"github.com/codegangsta/martini"
	"github.com/mvrilo/mstat"
)

func main() {
	m := martini.Classic()

	//if martini.Env != "production" {
	m.Use(mstat.New().ServeHTTP)
	//}

	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}
