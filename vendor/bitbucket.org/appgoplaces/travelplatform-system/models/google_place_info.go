package models

import (
	"encoding/json"

	locationpb "bitbucket.org/appgoplaces/service-protos/location"
)

type AddPlaceInfo struct {
	tableName struct{} `sql:"google_maps.place_info"`
	// FormattedAddress is the human-readable address of this place
	FormattedAddress string `sql:"formatted_address"`
	// Geometry contains geometry information about the result, generally including the
	// location (geocode) of the place and (optionally) the viewport identifying its
	// general area of coverage.
	AddressGeometry string `sql:"address_geometry"`
	// Name contains the human-readable name for the returned result. For establishment
	// results, this is usually the business name.
	Name string `sql:"name"`
	// Icon contains the URL of a recommended icon which may be displayed to the user
	// when indicating this result.
	Icon string `sql:"icon"`
	// PlaceID is a textual identifier that uniquely identifies a place.
	PlaceID string `sql:"place_api_id"`
	// Scope indicates the scope of the PlaceID.
	Scope string `sql:"scope"`
	// Rating contains the place's rating, from 1.0 to 5.0, based on aggregated user
	// reviews.
	Rating float32 `sql:"rating"`
	// Types contains an array of feature types describing the given result.
	Types []string `sql:"types" pg:",array"`
	// OpeningHours may contain whether the place is open now or not.
	OpeningHours string `sql:"opening_hours"`
	// Photos is an array of photo objects, each containing a reference to an image.
	Photos string `sql:"photos"`
	// AltIDs — An array of zero, one or more alternative place IDs for the place, with
	// a scope related to each alternative ID.
	AltIDs string `sql:"alt_ids"`
	// PriceLevel is the price level of the place, on a scale of 0 to 4.
	PriceLevel int32 `sql:"price_level"`
	// Vicinity contains a feature name of a nearby location.
	Vicinity string `sql:"vicinity"`
	// PermanentlyClosed is a boolean flag indicating whether the place has permanently
	// shut down.
	PermanentlyClosed bool `sql:"permanently_closed"`
}

type PlaceCityCountry struct {
	tableName          struct{} `sql:"google_maps.place_city_country"`
	PlaceCityCountryID int64    `sql:"place_city_country_id, pk"`
	PlaceAPIID         string   `sql:"place_api_id"`
	CityID             int64    `sql:"city_id"`
	CountryID          int64    `sql:"country_id"`
}

type PlaceCityCountryQuery struct {
	PlaceCityCountryID int64  `sql:"place_city_country_id, pk" json:"place_city_country_id"`
	PlaceAPIID         string `sql:"place_api_id" json:"place_api_id"`
	CityID             int64  `sql:"city_id" json:"city_id"`
	CountryID          int64  `sql:"country_id" json:"country_id"`
}

// Takes list of proto locationpb.PlacesSearchResult
func (db *Db) ProcessInsertPlaces(results []*locationpb.PlacesSearchResult, cityID int64, countryID int64) []error {
	var errors []error
	// Check if place info exist before inserting
	for _, result := range results {
		geometry, geoErr := json.Marshal(result.AddressGeometry)
		if geoErr != nil {
			errors = append(errors, geoErr)
		}
		openingHours, openingHoursErr := json.Marshal(result.OpeningHours)
		if openingHoursErr != nil {
			errors = append(errors, openingHoursErr)
		}
		photos, photosErr := json.Marshal(result.Photos)
		if photosErr != nil {
			errors = append(errors, photosErr)
		}
		altIds, altIdsErr := json.Marshal(result.AltIDs)
		if altIdsErr != nil {
			errors = append(errors, altIdsErr)
		}
		if len(errors) > 0 {
			return errors
		}
		placeInfo := AddPlaceInfo{
			FormattedAddress:  result.FormattedAddress,
			AddressGeometry:   string(geometry),
			Name:              result.Name,
			Icon:              result.Icon,
			PlaceID:           result.PlaceID,
			Scope:             result.Scope,
			Rating:            result.Rating,
			Types:             result.Types,
			OpeningHours:      string(openingHours),
			Photos:            string(photos),
			AltIDs:            string(altIds),
			PriceLevel:        result.PriceLevel,
			Vicinity:          result.Vicinity,
			PermanentlyClosed: result.PermanentlyClosed,
		}

		placeCityCountry := PlaceCityCountry{
			PlaceAPIID: result.PlaceID,
			CityID:     cityID,
			CountryID:  countryID,
		}

		_, err := db.Client.Model(&placeInfo).
			Where("place_api_id = ?", result.PlaceID).
			SelectOrInsert()
		if err != nil {
			errors = append(errors, err)
		}

		err = db.Client.Insert(&placeCityCountry)

		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
