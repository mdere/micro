package models

type PlaceType struct {
	tableName   struct{} `sql:"google_maps.city"`
	PlaceTypeID string   `sql:"place_type_id, pk"`
	Name        string   `sql:"name"`
}

type PlaceTypeQuery struct {
	PlaceTypeID string `sql:"place_type_id, pk" json:"place_type_id"`
	Name        string `sql:"name" json:"name"`
}

func (db *Db) GetPlaceTypes() ([]PlaceTypeQuery, error) {
	var placeTypeData []PlaceTypeQuery
	_, err := db.Client.Query(&placeTypeData, `SELECT * FROM "google_maps".place_type`)
	return placeTypeData, err
}
