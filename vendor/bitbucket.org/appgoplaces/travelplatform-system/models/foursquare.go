package models

import (
	"net/http"

	"encoding/json"

	"github.com/go-log/log"
	"github.com/peppage/foursquarego"
)

func FourSquare(queryInput string) (string, string) {

	httpClient := http.DefaultClient
	var clientId = "MM1OHBYY2QGQPCVTYM2GP1PEQK0PKOWUOM11X0XALJEVIAWG"
	var clientSecret = "XAKTBN4ASNO4O2ZHLUS3DDH4ASVBZOWBEJS33ULMVRRUM0BS"
	// When creating the client you can specify either clientSecret or the accesstoken
	client := foursquarego.NewClient(httpClient, "foursquare", clientId, clientSecret, "")

	params := foursquarego.VenueSuggestParams{
		Query: queryInput,
	}

	// Search Venues
	venues, resp, err := client.Venues.SuggestCompletion(&params)

	if err != nil {
		log.Log(err)
	}

	venuesJSON, error := json.Marshal(venues)

	if error != nil {
		log.Log(err)
	}

	return string(venuesJSON), string(resp.Status)
}
