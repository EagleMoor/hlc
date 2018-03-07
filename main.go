package main

import (
	"flag"
	"log"

	"hlc/models"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// flags
var (
	httpFlag     string
	dataPathFlag string
)

func init() {
	flag.StringVar(&httpFlag, "http", ":80", "HTTP address binding")
	flag.StringVar(&dataPathFlag, "data-path", "/tmp/hlc/data", "Path to data ZIP-file")
}

func main() {
	flag.Parse()

	storage := models.NewStorage()

	err := importDataFromPath(storage, dataPathFlag)
	if err != nil {
		log.Fatal(err)
	}

	server := &fasthttp.Server{
		Handler: makeRouter(storage),
		DisableHeaderNamesNormalizing: true,
	}

	log.Printf("HTTP server's starting on %s", httpFlag)
	err = server.ListenAndServe(httpFlag)
	if err != nil {
		log.Fatal(err)
	}
}

func makeRouter(storage *models.Storage) fasthttp.RequestHandler {
	router := routing.New()

	router.Get("/users/<id>", getUserByIDHandler(storage))
	router.Get("/locations/<id>", getLocationByIDHandler(storage))
	router.Get("/visits/<id>", getVisitByIDHandler(storage))

	router.Post("/users/new", createEntityHandler(storage, models.UserEntityType))
	router.Post("/locations/new", createEntityHandler(storage, models.LocationEntityType))
	router.Post("/visits/new", createEntityHandler(storage, models.VisitsEntityType))

	router.Post("/users/<id>", updateEntityHandler(storage, models.UserEntityType))
	router.Post("/locations/<id>", updateEntityHandler(storage, models.LocationEntityType))
	router.Post("/visits/<id>", updateEntityHandler(storage, models.VisitsEntityType))

	router.Get("/users/<id>/visits", getVisitsByUserIDHandler(storage))
	router.Get("/locations/<id>/avg", getLocationAvgMarkHandler(storage))

	return router.HandleRequest
}
