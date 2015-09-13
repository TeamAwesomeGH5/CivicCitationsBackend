package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/emicklei/go-restful"
)

var NoCitationFoundText string = "No citations were found matching the provided criteria."
var NoCitationFound = errors.New(NoCitationFoundText)

type Citation struct {
	Id                   int    `json:"id"`
	CitationNumber       uint64 `json:"citation_number"`
	CitationDate         string `json:"citation_date"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	DOB                  string `json:"date_of_birth"`
	DefendantAddress     string `json:"defendant_address"`
	DefendantCity        string `json:"defendant_city"`
	DefendantState       string `json:"defendant_state"`
	DriversLicenseNumber string `json:"drivers_license_number"`
	CourtDate            string `json:"court_date"`
	CourtLocation        string `json:"court_location"`
	CourtAddress         string `json:"court_address"`
}

func NewCitation() Citation {
	return Citation{
		Id:                   0,
		CitationNumber:       0,
		CitationDate:         "",
		FirstName:            "",
		LastName:             "",
		DOB:                  "",
		DefendantAddress:     "",
		DefendantCity:        "",
		DefendantState:       "",
		DriversLicenseNumber: "",
		CourtDate:            "",
		CourtLocation:        "",
		CourtAddress:         "",
	}
}

type CitationResponse struct {
	Citations []Citation
	Valid     bool
	Message   string
}

type CitationManager struct {
	Sources []Getter
}

func (cm *CitationManager) AddSource(g Getter) {
	cm.Sources = append(cm.Sources, g)
}

func (cm *CitationManager) findCitation(request *restful.Request, response *restful.Response) {
	citationNumber := request.PathParameter("number")
	number, err := GetCitationNumber(citationNumber)
	if err != nil {
		response.WriteEntity(CitationResponse{Message: err.Error()})
		return
	}
	citations := []Citation{}
	for _, getter := range cm.Sources {
		citation, err := getter.GetCitationByNumber(number)
		if err != nil && err != NoCitationFound {
			log.Printf("There was an error getting citations from %s: %s", getter.String(), err)
		}

		citations = append(citations, citation...)
	}

	if len(citations) < 1 {
		response.WriteEntity(CitationResponse{Citations: citations, Valid: false, Message: NoCitationFoundText})
		return
	}
	response.WriteEntity(CitationResponse{Citations: citations, Valid: true, Message: ""})
}

type Params struct {
	LicenseNumber string `json:"license_number"`
	LastName      string `json:"last_name"`
	Dob           string `json:"dob"`
}

func (cm *CitationManager) findAllCitationsForUser(request *restful.Request, response *restful.Response) {
	body, err := ioutil.ReadAll(request.Request.Body)
	defer request.Request.Body.Close()
	if err != nil {
		log.Printf("Unable to print body %s", err)

	}
	var params Params
	err = json.Unmarshal(body, &params)
	if err != nil {
		log.Printf("Error reading request body\n %s\n %s", string(body), err)
		response.WriteEntity(CitationResponse{Message: fmt.Sprintf("Could not read request body %s", string(body))})
	}

	// lastName, err := request.BodyParameter("last_name")
	// if err != nil {
	// 	log.Printf("lastName is invalid: %s", lastName)
	// }
	// log.Printf("lastName is %s", params.LastName)
	// licenseNumber, err := request.BodyParameter("license_number")
	// if err != nil {
	// 	log.Printf("lastName is invalid: %s", licenseNumber)
	// }
	// log.Printf("license_number is %s", params.LicenseNumber)
	// dob, err := request.BodyParameter("dob")
	// if err != nil {
	// 	log.Printf("lastName is invalid: %s", dob)
	// }
	// log.Printf("dob is %s", params.Dob)

	citations := []Citation{}
	for _, getter := range cm.Sources {
		citation, err := getter.GetCitationsByUser(params.LastName, params.LicenseNumber, params.Dob)
		if err != nil && err != NoCitationFound {
			log.Printf("There was an error getting citations from %s: %s", getter.String(), err)
		}

		citations = append(citations, citation...)
	}

	if len(citations) < 1 {
		response.WriteEntity(CitationResponse{Citations: citations, Valid: false, Message: NoCitationFoundText})
		return
	}
	response.WriteEntity(CitationResponse{Citations: citations, Valid: true, Message: ""})
}

func GetCitationNumber(citationNumber string) (uint64, error) {
	parsedCitationNumber, err := strconv.ParseUint(citationNumber, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Citation number invalid - must be numeric")
	}
	if parsedCitationNumber == 0 {
		return 0, fmt.Errorf("Citation number invalid - blank string")
	}
	return parsedCitationNumber, nil
}
