package trec

import (
	"os"

	"github.com/codegangsta/martini-contrib/binding"
	"github.com/dave-malone/email"
	"github.com/dave-malone/trec/user"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xchapter7x/lo"
)

// NewServer configures and returns a Server.
func NewServer() *martini.ClassicMartini {
	m := martini.Classic()
	initRoutes(m)
	initMappings(m)
	m.Use(render.Renderer())

	return m
}

func initMappings(m *martini.ClassicMartini) {
	profile := os.Getenv("PROFILE")

	if profile == "mysql" {
		db, err := newDbConn()
		if err != nil {
			user.NewRepository = user.NewMysqlRepositoryFactory(db)
		}
	} else {
		lo.G.Info("Using in-memory repositories")
		user.NewRepository = user.NewInMemoryRepository
	}

	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if awsEndpoint != "" && awsAccessKeyID != "" && awsSecretAccessKey != "" {
		email.NewSender = email.NewAmazonSESSender(awsEndpoint, awsAccessKeyID, awsSecretAccessKey)
	} else {
		email.NewSender = email.NewNoopSender
	}

	m.Map(email.NewSender())
	m.Map(newUserRepository())
}

func initRoutes(m *martini.ClassicMartini) {
	m.Get("/", func() string {
		return "trec api home; nothing to see here"
	})

	m.Group("/user", func(r martini.Router) {
		r.Get("/info", func() string {
			return "An API that allows you to work with Users"
		})
		r.Get("/", getUsersHandler)
		r.Get("/:id", getUserHandler)
		r.Post("/", binding.Json(User{}), createUserHandler)
	})

	m.Group("/company", func(r martini.Router) {
		r.Get("/info", func() string {
			return "An API that allows you to work with Companies"
		})
	})
}