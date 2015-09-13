package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func main() {
	configFile := flag.String("config", "config.json", "-config <configfile>: use to override the default config file (config.json)")
	flag.Parse()
	config := ParseConfig(*configFile)
	if config.DBUser == "" {
		log.Fatalf("Could not read config file! exiting....")
	}

	cm := &CitationManager{}
	cm.AddSource(NewSampleGetter(config.DBUser, config.DBPassword, config.DBAddress, config.Database))

	ws := new(restful.WebService)
	ws.Path("/api").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("citation/{number}").To(cm.findCitation).
		Doc("Gets a citation by number").
		Param(ws.PathParameter("number", "identifier of the citation").DataType("uint64")).
		Writes(CitationResponse{}))

	ws.Route(ws.POST("citations/").To(cm.findAllCitationsForUser).
		Doc("Get all citations for a user given the correct params.").
		Param(ws.BodyParameter("last_name", "Last name of the person to get citations for.")).
		Param(ws.BodyParameter("dob", "The date of birth of the person to get citations for.")).
		Param(ws.BodyParameter("license_number", "License number of the person to get citations for.")).
		Writes(CitationResponse{}))

	restful.Add(ws)
	log.Printf("Starting server on :%d", config.ServerPort)
	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.ServerPort), nil)
	if err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
