package models

import (
	"fmt"

	user "bitbucket.org/appgoplaces/service-protos/user"
)

type User struct {
	Id              int64  `sql:"id"`
	Email           string `sql:"email"`
	Firstname       string `sql:"firstname"`
	Lastname        string `sql:"lastname"`
	Password        string `sql:"password"`
	VerifyToken     string `sql:"verify_token"`
	VerifyCode      int32  `sql:"verify_code"`
	LastSignedin    int64  `sql:"last_signedin"`
	Verified        bool   `sql:"verified"`
	Mobile          int32  `sql:"mobile"`
	JwtToken        string `sql:"jwt_token"`
	PreferenceSetup bool   `sql:"preference_setup"`
}

// type UserCreate struct {
// 	Email       string `sql:"email"`
// 	Firstname   string `sql:"firstname"`
// 	Lastname    string `sql:"lastname"`
// 	Password    string `sql:"password"`
// 	VerifyToken string `sql:"verify_token"`
// 	VerifyCode  int32  `sql:"verify_code"`
// 	Mobile      int32  `sql:"mobile"`
// }

type UserIdentification struct {
	Id    int64
	Email string
}

type UserInfo struct {
	tableName       struct{} `sql:"users"`
	ID              string   `sql:"id"`
	Email           string   `sql:"email"`
	Firstname       string   `sql:"firstname"`
	Lastname        string   `sql:"lastname"`
	Verified        bool     `sql:"verified"`
	PreferenceSetup bool     `sql:"preference_setup"`
}

type UserPreference struct {
	tableName  struct{} `sql:"discovr.user_preference"`
	Countries  []int64  `pg:",array" sql:"countries" json:"countries"`
	Activities []int64  `pg:",array" sql:"activities" json:"activities"`
}

type UserPreferenceCreate struct {
	tableName  struct{} `sql:"discovr.user_preference"`
	UserId     int64    `sql:"user_id"`
	Countries  []int64  `sql:"countries" json:"countries"`
	Activities []int64  `sql:"activities" json:"activities"`
}

type AllRegionPreferences struct {
	RegionId int64  `sql:"region_id"`
	Name     string `sql:"region_name"`
	ImgSrc   string `sql:"img_src"`
}

type AllCountryPreferences struct {
	CountryId int64  `sql:"country_id"`
	RegionId  int64  `sql:"region_id"`
	ImgSrc    string `sql:"img_src"`
	Name      string `sql:"name"`
}

type AllRegionCountriesPreferences struct {
	Region    AllRegionPreferences
	Countries []AllCountryPreferences
}

type AllActivityPreferences struct {
	ActivityId int64  `sql:"activity_id"`
	ImgSrc     string `sql:"img_src"`
	Name       string `sql:"name"`
}

type UpdatePreferenceResponse struct {
	Countries  []*int64
	Activities []*int64
}

func (db *Db) GetAllPreferences() (*user.GetAllPreferencesResponse, error) {
	// {
	// 	regions: [{
	// 		name,
	// 		id
	// 		}],
	// 	regionCountries: [{
	// 		region: {},
	// 		countries: []
	// 		}].
	// 	activities: []
	// }
	var allPreferences user.GetAllPreferencesResponse
	// Get Regions and assign
	regions, regionErr := db.GetAllRegionPreferences()
	if regionErr != nil {
		fmt.Println("1")
		fmt.Println(regionErr)
		panic(regionErr)
	}
	// Get Countries and assign
	countries, countryErr := db.GetAllCountryPreferences()
	if countryErr != nil {
		fmt.Println("2")
		fmt.Println(countryErr)
		panic(countryErr)
	}
	// Get Activities and assign
	activities, activityErr := db.GetAllActivityPreferences()
	if activityErr != nil {
		fmt.Println("3")
		fmt.Println(activityErr)
		panic(activityErr)
	}

	var pbregionCountries []*user.RegionCountryPreferences
	pbRegionPreferences := make([]*user.RegionPreferences, 0, len(regions))
	// Organize Countries to a region
	for _, region := range regions {
		var pbregionCountry user.RegionCountryPreferences
		var pbCountries []*user.CountryPreferences
		for _, country := range countries {
			if country.RegionId == region.RegionId {
				pbCountry := user.CountryPreferences(country)
				pbCountries = append(pbCountries, &pbCountry)
			}
		}
		pbRegion := user.RegionPreferences(region)
		pbregionCountry.Region = &pbRegion
		fmt.Println(pbCountries)
		pbregionCountry.Countries = pbCountries
		pbregionCountries = append(pbregionCountries, &pbregionCountry)
		pbRegionPreferences = append(pbRegionPreferences, &pbRegion)
	}

	pbActivityPreferences := make([]*user.ActivityPreferences, 0, len(activities))

	for _, activity := range activities {
		activitypb := user.ActivityPreferences(activity)
		pbActivityPreferences = append(pbActivityPreferences, &activitypb)
	}

	allPreferences.Regions = pbRegionPreferences
	allPreferences.RegionGroupings = pbregionCountries
	allPreferences.Activities = pbActivityPreferences

	return &allPreferences, nil
}

func (db *Db) GetAllRegionPreferences() ([]AllRegionPreferences, error) {
	var regions []AllRegionPreferences
	_, err := db.Client.Query(&regions, `SELECT region_id, region_name, img_src FROM "discovr".region ORDER BY region_name`)
	return regions, err
}

func (db *Db) GetAllCountryPreferences() ([]AllCountryPreferences, error) {
	var countries []AllCountryPreferences
	_, err := db.Client.Query(&countries, `SELECT country_id, name, region_id, img_src FROM "discovr".country ORDER BY region_id, name`)
	return countries, err
}

func (db *Db) GetAllActivityPreferences() ([]AllActivityPreferences, error) {
	var activities []AllActivityPreferences
	_, err := db.Client.Query(&activities, `SELECT activity_id, name, img_src FROM "discovr".activity`)
	return activities, err
}

func (db *Db) UpdatePreference(u *UserPreference, userId int64) error {
	_, err := db.Client.Model(u).Where("user_id = ?", userId).Update()
	return err
}

func (db *Db) CreatePreference(u *UserPreferenceCreate) error {
	return db.Client.Insert(u)
}

func (db *Db) CreateUser(u *User) (*User, error) {
	_, err := db.Client.Model(u).
		Column("id").
		Where("email = ?", u.Email).
		Returning("id").
		SelectOrInsert()
	return u, err
}

type UserVerification struct {
	TableName   struct{} `sql:"public.users"`
	Email       string   `sql:"email,pk"`
	Mobile      int32    `sql:"mobile"`
	VerifyCode  int32    `sql:"verify_code"`
	VerifyToken string   `sql:"verify_token"`
}

func (db *Db) UpdateUserVerification(uv *UserVerification) error {
	_, err := db.Client.Model(uv).Where("email = ?", uv.Email).Update()
	return err
}

func (db *Db) GetUserById(id int64) (User, error) {
	var user User
	_, err := db.Client.Query(&user, `SELECT * FROM Users WHERE id = ?`, id)
	return user, err
}

func (db *Db) GetUserInfoById(id int64) (*UserInfo, error) {
	user := new(UserInfo)
	err := db.Client.Model(user).Where("id = ?", id).Select()
	return user, err
}

func (db *Db) GetUserByEmail(email string) (*User, error) {
	var user User
	_, err := db.Client.Query(&user, `SELECT * FROM Users WHERE email = ?`, email)
	return &user, err
}

func (db *Db) GetUserInfoByEmail(email string) (*UserInfo, error) {
	var user UserInfo
	_, err := db.Client.Query(&user, `SELECT email, firstname, lastname, verified
        FROM Users WHERE email = ?`, email)
	return &user, err
}
