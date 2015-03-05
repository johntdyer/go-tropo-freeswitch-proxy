package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/pmylund/go-cache"
	"encoding/xml"
)

// GoAuthProxyConfig is the primary application struct
type GoAuthProxyConfig struct {
	AppName                    string
	PropertyId                 string
	LogLevel                   string
	Version                    string
	BuildDate                  string
	PapiUser                   string
	PapiPass                   string
	PapiUrl                    string
	BasicAuthUser              string
	BasicAuthPass              string
	FreeSwitchUserCacheTimeout string
	CacheNegativeTTL           int
	CacheTTL                   int
	ExpiredCachePurgeInterval  int
	ListenPort                 string
	ConnectDomain              string
	DefaultTollPlan            string
	AuthProxyCache             *cache.Cache
	AppCacheTimeout            int
}

// Config data to be used by Proxy to look up address data by ID or key
type ConfigData struct {
	Id   map[int]string
	Name map[string]string
}

// Provisioning API response
type PapiResponse struct {
	Address string `json:"address"`
	Configs []struct {
		Href        string `json:"href"`
		Id          string `json:"id"`
		Description string `json:"description"`
		Value       string `json:"value"`
		Name        string `json:"name"`
	}
}

// DirectoryAuthResponse stucted used the the Directory handler
type DirectoryAuthResponse struct {
	Message     string
	XmlResponse string
	Fields      log.Fields
}

// Struct for handling respoes from UserAuthHandler
type UserAuthHandlerResponse struct {
	Message string
	Header  int
	Address string
	Fields  log.Fields
}

// AppData contains information about the version of the application
type AppData struct {
	Name           string `json:"application"`
	Version        string `json:"version"`
	BuildDate      string `json:"build_date"`
	ApiConnectivty bool   `json:"api_connectivity"`
}

type Auth struct {
	Href        string `json:"href"`
	Id          string `json:"id"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Address     string `json:"address"`
}

type FreeswitchRequest struct {
	Domain          string
	Ip              string
	UserAgent       string
	Action          string
	Username        string
	SipAuthUsername string
}

type Variable struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Address struct {
	XMLName xml.Name `xml:"domain"`
	Name    string   `xml:"name,attr"`
	User    *User    `xml:"user"`
}

type Document struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Section Section  `xml:"section"`
}

type Section struct {
	Name   string `xml:"name,attr"`
	Domain Domain `xml:"domain"`
}

type SectionResult struct {
	Name   string `xml:"name,attr"`
	Result Result `xml:"result"`
}
type Domain struct {
	Name      string     `xml:"name,attr"`
	Params    []Param    `xml:"params>param"`
	Variables []Variable `xml:"variables>variable"`
	Groups    []*Group   `xml:"groups>group"`
}

type Group struct {
	Name  string  `xml:"name,attr"`
	Users []*User `xml:"users>user"`
}

type User struct {
	Id       string  `xml:"id,attr"`
	Cachable string  `xml:"cacheable,attr"`
	Params   []Param `xml:"params>param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type NotFound struct {
	XMLName xml.Name      `xml:"document"`
	Type    string        `xml:"type,attr"`
	Section SectionResult `xml:"result"`
}

type Result struct {
	Status string `xml:"status,attr"`
}
