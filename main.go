package main

import (
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful"
)

func main() {
	cm := &CitationManager{}
	cm.AddSource(NewSampleGetter("root", "my-secret-pw", "tcp(192.168.1.116:3306)", "civicCitations"))

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
	fmt.Println("Starting server on :6969")
	http.ListenAndServe("0.0.0.0:6969", nil)
}
