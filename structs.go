package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/pmylund/go-cache"
	"encoding/xml"
)

// GoAuthProxyConfig is the primary application struct
type GoAuthProxyConfig struct {
	AppName                    string
	LogLevel                   string
	Version                    string
	BuildDate                  string
	PapiUser                   string
	PapiPass                   string
	PapiURL                    string
	BasicAuthUser              string
	ConnectDomain              string
	ValidateDomain             bool
	BasicAuthPass              string
	FreeSwitchUserCacheTimeout string
	CacheNegativeTTL           int
	CacheTTL                   int
	ExpiredCachePurgeInterval  int
	ListenPort                 string
	DefaultTollPlan            string
	AuthProxyCache             *cache.Cache
	AppCacheTimeout            int
}

// ConfigData to be used by Proxy to look up address data by ID or key
type ConfigData struct {
	ID   map[int]string
	Name map[string]string
}

// PapiResponse - Provisioning API response
type PapiResponse struct {
	Address string `json:"address"`
	Configs []struct {
		Href        string `json:"href"`
		ID          string `json:"id"`
		Description string `json:"description"`
		Value       string `json:"value"`
		Name        string `json:"name"`
	}
}

// DirectoryAuthResponse stucted used the the Directory handler
type DirectoryAuthResponse struct {
	Message     string
	XMLResponse string
	Fields      log.Fields
}

// UserAuthHandlerResponse - Struct for handling respoes from UserAuthHandler
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
	APIConnectivty bool   `json:"api_connectivity"`
}

// Auth - User auth data
type Auth struct {
	Href        string `json:"href"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Address     string `json:"address"`
}

// FreeswitchRequest - request data
type FreeswitchRequest struct {
	Domain          string
	IP              string
	UserAgent       string
	Action          string
	Username        string
	SipAuthUsername string
}

// Variable - xml response data
type Variable struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// Address - address data
type Address struct {
	XMLName xml.Name `xml:"domain"`
	Name    string   `xml:"name,attr"`
	User    *User    `xml:"user"`
}

// Document - XML doc
type Document struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Section Section  `xml:"section"`
}

// Section - Section directory data
type Section struct {
	Name   string `xml:"name,attr"`
	Domain Domain `xml:"domain"`
}

// SectionResult - Result
type SectionResult struct {
	Name   string `xml:"name,attr"`
	Result Result `xml:"result"`
}

// Domain - Domain data
type Domain struct {
	Name      string     `xml:"name,attr"`
	Params    []Param    `xml:"params>param"`
	Variables []Variable `xml:"variables>variable"`
	Groups    []*Group   `xml:"groups>group"`
}

// Group - group data
type Group struct {
	Name  string  `xml:"name,attr"`
	Users []*User `xml:"users>user"`
}

// User - User
type User struct {
	ID        string     `xml:"id,attr"`
	Cachable  string     `xml:"cacheable,attr"`
	Params    []Param    `xml:"params>param"`
	Variables []Variable `xml:"variables>variable"`
}

// Param - directory params
type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// NotFound  - 404 user
type NotFound struct {
	XMLName xml.Name      `xml:"document"`
	Type    string        `xml:"type,attr"`
	Section SectionResult `xml:"result"`
}

// Result - user status
type Result struct {
	Status string `xml:"status,attr"`
}
