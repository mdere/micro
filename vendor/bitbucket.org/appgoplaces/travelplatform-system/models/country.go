package models

type Country struct {
	tableName    struct{} `sql:"discovr.country"`
	CountryID    int64    `sql:"country_id, pk"`
	Name         string   `sql:"name"`
	ISO          string   `sql:"iso"`
	ISO3         string   `sql:"iso3"`
	ISONum       string   `sql:"isonum"`
	FIPS         string   `sql:"fips"`
	Captital     string   `sql:"capital"`
	RegionID     string   `sql:"region_id"`
	TLD          string   `sql:"tld"`
	CurrencyCode string   `sql:"currency_code"`
	CurrencyName string   `sql:"currency_name"`
	PostalFormat string   `sql:"postal_format"`
	Language     string   `sql:"language"`
	GeoNameID    string   `sql:"geonameid"`
	PhoneCode    string   `sql:"phonecode"`
	ImgSrc       string   `sql:"img_src"`
}

type CountryQuery struct {
	tableName    struct{} `sql:"discovr.country"`
	CountryID    int64    `sql:"country_id, pk" json:"country_id"`
	Name         string   `sql:"name" json:"name"`
	ISO          string   `sql:"iso" json:"iso"`
	ISO3         string   `sql:"iso3" json:"iso3"`
	ISONum       string   `sql:"isonum" json:"isonum"`
	FIPS         string   `sql:"fips" json:"fips"`
	Captital     string   `sql:"capital" json:"capital"`
	RegionID     string   `sql:"region_id" json:"region_id"`
	TLD          string   `sql:"tld" json:"tld"`
	CurrencyCode string   `sql:"currency_code" json:"currency_code"`
	CurrencyName string   `sql:"currency_name" json:"currency_name"`
	PostalFormat string   `sql:"postal_format" json:"postal_format"`
	Language     string   `sql:"language" json:"language"`
	GeoNameID    string   `sql:"geonameid" json:"geonameid"`
	PhoneCode    string   `sql:"phonecode" json:"phonecode"`
	ImgSrc       string   `sql:"img_src" json:"img_src"`
}

func (db *Db) GetCountryData() ([]CountryQuery, error) {
	var countryData []CountryQuery
	_, err := db.Client.Query(&countryData, `SELECT * FROM "discovr".country`)
	return countryData, err
}

func (db *Db) GetCountryDataByName(countryName string) (Country, error) {
	var countryData Country
	_, err := db.Client.Query(&countryData, `SELECT * FROM "discovr".country WHERE name LIKE ?`, countryName)
	return countryData, err
}

func (db *Db) GetCountryDataByRegionID(regionID string) (Country, error) {
	var countryData Country
	_, err := db.Client.Query(&countryData, `SELECT * FROM "discovr".country WHERE region_id = ?`, regionID)
	return countryData, err
}

func (db *Db) GetCountryDataByISO(isoValue string) (Country, error) {
	var countryData Country
	_, err := db.Client.Query(&countryData, `SELECT * FROM "discovr".country WHERE iso LIKE ?`, isoValue)
	return countryData, err
}
