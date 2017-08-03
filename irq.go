package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/stvp/go-toml-config"

	"github.com/firba1/irq/model"
	"github.com/firba1/irq/view"
)

// command line flags
var (
	Port   = flag.Int("port", 3000, "port to run on")
	Config = flag.String("config", "irq.toml", "path for config for db and stuff")
)

// toml config values
var (
	dbType = config.String("db.type", "sqlite3")
	dbPath = config.String("db.path", "quotes.db")
)

func main() {
	flag.Parse()
	err := config.Parse(*Config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m := martini.Classic()

	// allow for custom assets
	m.Use(martini.Static("custom_assets"))

	m.Use(render.Renderer(render.Options{
		Layout: "base",
	}))

	m.Use(func(req *http.Request, c martini.Context) {
		c.Map(view.IsJson(strings.HasSuffix(req.URL.Path, "json")))
	})

	db, err := model.NewModel(*dbType, *dbPath)
	// TODO allow this to be configured
	db.SetMaxOpenConns(100)
	if err != nil {
		panic(err)
	}

	// middleware to inject DB connection into each request
	m.Use(func(r render.Render, c martini.Context, isJson view.IsJson) {
		c.Map(db)
	})

	m.Get("/", view.Index)
	m.Get(json("/latest"), view.Latest)
	m.Get(json("/all"), view.All)
	m.Get(json("/random"), view.Random)
	m.Get(json("/search"), view.Search)
	m.Get("/submit", view.Submit)
	m.Get(json("/top"), view.Top)
	m.Get(json("/quote/:id"), view.Quote)
	m.Get("/status", view.Status)

	m.NotFound(func(r render.Render, isJson view.IsJson) {
		view.RenderError(r, 404, isJson, "not found")
		return
	})

	m.Post("/submit", view.SubmitForm)
	m.Post("/vote/:id", view.Vote)
	m.Delete("/quote/:id", view.Delete)

	m.RunOnAddr(fmt.Sprintf(":%v", *Port))
}

func json(pattern string) string {
	return pattern + "((/|.)json|)"
}
