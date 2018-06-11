package conf

import (
	"fmt"

	config "github.com/micro/go-config"
	envvar "github.com/micro/go-config/source/envvar"
	"github.com/micro/go-log"
)

type smtpOpts struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// for a struct scan
type rdbmsOpts struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

type apiOpts struct {
	Url              string `json:"url"`
	Port             string `json:"port"`
	FoursquareSecret string `json:"foursquare"`
}

type domainOpts struct {
	Url  string `json:"url"`
	Port string `json:"port"`
}

type EnvOpts struct {
	Api    apiOpts
	Smtp   smtpOpts
	Domain domainOpts
	Rdbms  rdbmsOpts
}

func GetConf() *EnvOpts {
	envvarsrc := envvar.NewSource(
		envvar.WithStrippedPrefix("APP"),
		envvar.WithPrefix("DB"),
	)
	// Create new config
	conf := config.NewConfig()
	// Load file source
	conf.Load(envvarsrc)

	var envVar EnvOpts
	var smtpVar smtpOpts
	var apiVar apiOpts
	var domainVar domainOpts
	var rdbmsVar rdbmsOpts

	conf.Get("smtp").Scan(&smtpVar)
	conf.Get("api").Scan(&apiVar)
	conf.Get("domain").Scan(&domainVar)
	conf.Get("db", "rdbms").Scan(&rdbmsVar)

	envVar.Smtp = smtpVar
	envVar.Api = apiVar
	envVar.Domain = domainVar
	envVar.Rdbms = rdbmsVar

	log.Log(envVar)

	return &envVar
}

func GetApiHostname(env *EnvOpts) string {
	api := env.Api
	if api.Port != "" {
		return fmt.Sprintf("%s:%s", api.Url, api.Port)
	}
	return api.Url
}

func GetDomainHostname(env *EnvOpts) string {
	domain := env.Domain
	if domain.Port != "" {
		return fmt.Sprintf("%s:%s", domain.Url, domain.Port)
	}
	return domain.Url
}
