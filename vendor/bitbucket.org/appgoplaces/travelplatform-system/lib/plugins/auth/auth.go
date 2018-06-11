package auth_plugin

// stats_auth enables basic auth on the /stats endpoint

import (
	"net/http"
	"strconv"

	"github.com/micro/cli"
	"github.com/micro/micro/plugin"

	"github.com/micro/go-log"
	// "fmt"
	"strings"
	"regexp"
	"encoding/json"
	"github.com/pkg/errors"
	jwt "github.com/dgrijalva/jwt-go"
	rdbms "bitbucket.org/appgoplaces/travelplatform-system/db/rdbms"
	models "bitbucket.org/appgoplaces/travelplatform-system/models"
)

var (
	db = rdbms.Connect("micro")
	secretKey = []byte("secret-key")
)

var exclude = []string{
	"/user/register",
	"/user/signin",
	"/user/requestpasswordreset",
	"/user/sendverifycode",
	"/user/verifycode",
	"/management/user/signin",
	"/management/user/register",
}

type auth struct {
}

func (a *auth) Flags() []cli.Flag {
	return []cli.Flag{}
}

func (a *auth) Commands() []cli.Command {
	return nil
}

func (a *auth) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Log("--------------------")
			for _, url := range exclude {
				if url == strings.ToLower(r.RequestURI) {
					h.ServeHTTP(w,r)
					return
				}
			}
			// GET Token and append to header
			bearerToken := getBearerToken(r.Header["Authorization"])
			if len(bearerToken) > 0 {
				userIdentification, err := parseToken(bearerToken)
				if err == nil {
					var user userMap
					var qErr error
					if strings.Contains(r.RequestURI, "management") {
						qErr = user.authenticateManagement(userIdentification)
					} else {
						qErr = user.authenticateUser(userIdentification)
					}
					if qErr != nil {
						// respond back 500
						log.Log("qErr")
						log.Log(qErr)
						return
					}
					if !user.Valid {
						// respond back 401
						log.Log("NOT AUTHORIZED")
						return 
					}
					data, jsonErr := json.Marshal(&user)
					if jsonErr != nil {
						log.Log("jsonErr")
						log.Log(jsonErr)
						// respond back 500
						return
					}
					r.Header.Set("user", string(data))
				}  else if vErr, ok := errors.Cause(err).(*jwt.ValidationError); ok  {
					if vErr.Errors&jwt.ValidationErrorExpired > 0 {
						// respond back 401 and expired
						log.Log("NOT AUTHORIZED AND EXPIRED")
						return
					}
				}
			}
			h.ServeHTTP(w,r)
			log.Log("-------------------")
			return
		})
	}
}

type userMap struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Valid bool   `json:"valid"`
}

type userIdentity struct {
	Id    int64
	Email string
}

func (u *userMap) authenticateManagement(userIdentify userIdentity) error {
	user := models.ManagementUser{}
	_, qErr := db.Query(&user,`SELECT * FROM management_user WHERE email = ? AND management_user_id = ?`, userIdentify.Email, userIdentify.Id)
	if qErr != nil {
		if qErr.Error() != "pg: no rows in result set" {
			return qErr
		}
	}
	valid := len(user.Email) > 0 && user.Id > 0 && user.Role == "Admin" 
	u.Id = user.Id
	u.Email = user.Email
	u.Valid = valid
	return nil
}

func (u *userMap) authenticateUser(userIdentify userIdentity) error {
	user := models.User{}
	_, qErr := db.Query(&user,`SELECT * FROM Users WHERE email = ? AND management_user_id = ?`, userIdentify.Email, userIdentify.Id)
	if qErr != nil {
		if qErr.Error() != "pg: no rows in result set" {
			return qErr
		}
	}
	valid := len(user.Email) > 0 && user.Id > 0 
	u.Id = user.Id
	u.Email = u.Email
	u.Valid = valid
	return nil
}

func parseToken(bearerToken string) (userIdentity, error) {
	token, err := jwt.ParseWithClaims(bearerToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return secretKey, nil
	})

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid && claims.Valid() == nil {
		// Need to break down Subject into UserID model
		subjects := strings.Split(claims.Subject, ",")
		id, err := strconv.ParseInt(subjects[1], 10, 64)
		return userIdentity{id, subjects[0]}, err
	} else {
		return userIdentity{}, errors.Wrap(err, "parsing jwt token string")
	}
}

func getBearerToken(authorization []string) string {
	if len(authorization) > 0 {
		value := authorization[0]
		if bearerExists, err := regexp.MatchString("Bearer", value); bearerExists && err == nil {
			bearerData := strings.Split(value, " ")
			log.Log(bearerData)
			if len(bearerData) > 1 {
				return bearerData[1]
			}
		}
	}
	return ""
}

func (a *auth) Init(ctx *cli.Context) error {
	return nil
}

func (a *auth) String() string {
	return "auth"
}

func New() plugin.Plugin {
	return &auth{}
}
