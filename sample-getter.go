package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type SampleGetter struct {
	user     string
	password string
	address  string
	database string
}

func NewSampleGetter(user, password, address, database string) SampleGetter {
	sg := SampleGetter{
		user:     user,
		password: password,
		address:  address,
		database: database,
	}
	return sg
}

func (sg SampleGetter) String() string {
	return fmt.Sprintf("SampleGetter{address: %s, database: %s}", sg.address, sg.database)
}

func (sg SampleGetter) Query(querystring string) (*sql.Rows, error) {
	con, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", sg.user, sg.password, sg.address, sg.database))
	if err != nil {
		return nil, err
	}
	defer con.Close()

	rows, err := con.Query(querystring)

	return rows, err
}

func (sg SampleGetter) GetCitationByNumber(number uint64) ([]Citation, error) {
	rows, err := sg.Query(fmt.Sprintf("SELECT * FROM citations WHERE citation_number = '%d';", number))

	if err != nil {
		log.Printf("There was an error getting a citation from Sample-Getter: %s", err)
		return []Citation{}, err
	}

	//build Citations
	cits := []Citation{}

	for rows.Next() {
		var cit Citation = NewCitation()
		var scrap int
		if err := rows.Scan(&scrap, &cit.Id, &cit.CitationNumber, &cit.CitationDate, &cit.FirstName, &cit.LastName, &cit.DOB, &cit.DefendantAddress, &cit.DefendantCity, &cit.DefendantState, &cit.DriversLicenseNumber, &cit.CourtDate, &cit.CourtLocation, &cit.CourtAddress); err != nil {
			return []Citation{}, err
		}

		cits = append(cits, cit)
	}

	return cits, nil
}

func (sg SampleGetter) GetCitationsByUser(lastname, licensenumber, dob string) ([]Citation, error) {
	rows, err := sg.Query(fmt.Sprintf("SELECT * FROM citations WHERE last_name = '%s' and drivers_license_number='%s' and STR_TO_DATE(date_of_birth, '%s') = STR_TO_DATE('%s', '%s');", lastname, licensenumber, "%m/%d/%Y", dob, "%m/%d/%Y"))

	//log.Printf(fmt.Sprintf("SELECT * FROM citations WHERE last_name = '%s';", lastname))

	if err != nil {
		log.Printf("There was an error getting a citation from Sample-Getter: %s", err)
		return []Citation{}, err
	}

	//build Citations
	cits := []Citation{}
	for rows.Next() {
		var cit Citation = NewCitation()
		var scrap int
		if err := rows.Scan(&scrap, &cit.Id, &cit.CitationNumber, &cit.CitationDate, &cit.FirstName, &cit.LastName, &cit.DOB, &cit.DefendantAddress, &cit.DefendantCity, &cit.DefendantState, &cit.DriversLicenseNumber, &cit.CourtDate, &cit.CourtLocation, &cit.CourtAddress); err != nil {
			return []Citation{}, err
		}

		violas, err := sg.GetViolationsForCitation(cit)
		if err != nil {
			log.Printf("There was an error getting violations for citation %s: %s", cit.CitationNumber, err)
			return []Citation{}, err
		}

		cit.Violations = violas

		cits = append(cits, cit)
	}

	return cits, nil
}

func (sg SampleGetter) GetViolationsForCitation(citation Citation) (Violations, error) {
	rows, err := sg.Query(fmt.Sprintf("SELECT * FROM violations WHERE citation_number = '%d';", citation.CitationNumber))

	//log.Printf(fmt.Sprintf("SELECT * FROM citations WHERE last_name = '%s';", lastname))

	if err != nil {
		log.Printf("There was an error getting a violation from Sample-Getter: %s", err)
		return Violations{}, err
	}

	//build Citations
	violas := Violations{}

	for rows.Next() {
		var viol Violation = NewViolation()
		var scrap int
		if err := rows.Scan(&scrap, &viol.Id, &viol.Citation_number, &viol.Violation_number, &viol.Violation_description, &viol.Warrant_status, &viol.Warrant_number, &viol.Status, &viol.Status_date, &viol.Fine_amount, &viol.Court_cost); err != nil {
			return Violations{}, err
		}

		log.Printf("Found viola successfully: %s", viol.Violation_description)
		violas = append(violas, viol)
	}

	return violas, nil

}
