package models

import (
	locationpb "bitbucket.org/appgoplaces/service-protos/location"
)

type City struct {
	tableName           struct{} `sql:"google_maps.city"`
	CityID              int64    `sql:"city_id, pk"`
	GeoNameID           int64    `sql:"geoname_id"`
	LocaleCode          string   `sql:"locale_code"`
	ContinentCode       string   `sql:"continent_code"`
	ContinentName       string   `sql:"continent_name"`
	CountryISOCode      string   `sql:"country_iso_code"`
	CountryName         string   `sql:"country_name"`
	SubDivision1ISOCode string   `sql:"subdivision_1_iso_code"`
	SubDivision1Name    string   `sql:"subdivision_1_name"`
	SubDivision2ISOCode string   `sql:"subdivision_2_iso_code"`
	SubDivision2Name    string   `sql:"subdivision_2_name"`
	CityName            string   `sql:"city_name"`
	MetroCode           string   `sql:"metro_code"`
	TimeZone            string   `sql:"time_zone"`
	IsInEuropeanUnion   int64    `sql:"is_in_european_union"`
}

type CityQuery struct {
	CityID              int64  `sql:"city_id, pk" json:"city_id"`
	GeoNameID           int64  `sql:"geoname_id" json:"geoname_id"`
	LocaleCode          string `sql:"locale_code" json:"locale_code"`
	ContinentCode       string `sql:"continent_code" json:"continent_code"`
	ContinentName       string `sql:"continent_name" json:"continent_name"`
	CountryISOCode      string `sql:"country_iso_code" json:"country_iso_code"`
	CountryName         string `sql:"country_name" json:"country_name"`
	SubDivision1ISOCode string `sql:"subdivision_1_iso_code" json:"subdivision_1_iso_code"`
	SubDivision1Name    string `sql:"subdivision_1_name" json:"subdivision_1_name"`
	SubDivision2ISOCode string `sql:"subdivision_2_iso_code" json:"subdivision_2_iso_code"`
	SubDivision2Name    string `sql:"subdivision_2_name" json:"subdivision_2_name"`
	CityName            string `sql:"city_name" json:"city_name"`
	MetroCode           string `sql:"metro_code" json:"metro_code"`
	TimeZone            string `sql:"time_zone" json:"time_zone"`
	IsInEuropeanUnion   int64  `sql:"is_in_european_union" json:"is_in_european_union"`
}

// GetCityData returns list of cities associated to a specific ISO2
func (db *Db) GetCityData(countryISO string) ([]CityQuery, error) {
	var cityData []CityQuery
	_, err := db.Client.Query(&cityData, `SELECT * FROM "google_maps".city WHERE country_iso_code LIKE ?`, countryISO)
	return cityData, err
}

// SearchCityCountry Returns back City, Country
func (db *Db) SearchCityCountry(query string) ([]*locationpb.SearchCityCountry, error) {
	var cities []*locationpb.SearchCityCountry
	_, err := db.Client.Query(&cities, `SELECT search_result, city_id FROM "public".search_city_country(?)`, query)
	return cities, err
}

// SearchCityProvince Returns back City, Province
func (db *Db) SearchCityProvince(query string) ([]*locationpb.SearchCityCountry, error) {
	var cities []*locationpb.SearchCityCountry
	_, err := db.Client.Query(&cities, `SELECT search_result, city_id FROM "public".search_city_province(?)`, query)
	return cities, err
}
