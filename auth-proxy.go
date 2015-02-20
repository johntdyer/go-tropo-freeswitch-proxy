package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	_ "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	AppName          = "tropo-auth"
	buildDate        string
	applicationData  *AppVersion
	configPropertyId = os.Getenv("ADDRESS_CONFIG_PROPERTY_ID")
	PapiUser         = os.Getenv("TROPO_API_USER")
	PapiPass         = os.Getenv("TROPO_API_PASS")
	PapiUrl          = os.Getenv("TROPO_API_URL")
	BasicAuthUser    = os.Getenv("API_AUTH_USER")
	BasicAuthPass    = os.Getenv("API_AUTH_PASS")
	listenPort       = os.Getenv("LISTEN_PORT")
	cacheTime        = os.Getenv("USER_CACHE_VALUE")
	ConnectDomain    = "connect.tropo.com"
)

// versionRequestHandler handles incoming version / health requests
func VersionHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.WithFields(log.Fields{
		"method": req.Method,
		"url":    req.RequestURI,
	}).Debug("VersionHandler")
	w.Header().Add("Content-Type", "application/json")
	body, _ := json.Marshal(applicationData)
	fmt.Fprintf(w, string(body))

}

func AuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Debugf("AuthHandler: %s - %s", req.Method, req.RequestURI)

	log.WithFields(log.Fields{
		"method": req.Method,
		"url":    req.RequestURI,
	}).Debug("AuthHandler")

	auth := &AuthRequest{}

	if req.Method == "POST" {
		auth.Action = req.FormValue("action")
		auth.Username = req.FormValue("user")
		auth.Domain = req.FormValue("domain")
		auth.SipAuthUsername = req.FormValue("sip_auth_username")
	} else {
		q := req.URL.Query()
		auth.Action = q.Get("action")
		auth.Domain = q.Get("domain")
		auth.Username = q.Get("user")
		auth.SipAuthUsername = q.Get("sip_auth_username")

	}

	if auth.Action == "sip_auth" {

		w.Header().Set("Content-Type", "text/xml")

		addressData := GetAddressAuthData(auth.SipAuthUsername)
		if addressData.Value == "" {

			log.WithFields(log.Fields{
				"number":            auth.Username,
				"domain":            auth.Domain,
				"sip_auth_username": auth.SipAuthUsername,
			}).Debug("Address not found")

			fmt.Fprintf(w, RenderNotFound())
		} else {
			log.WithFields(log.Fields{
				"domain":            auth.Domain,
				"action":            auth.Action,
				"username":          auth.Username,
				"sip_auth_username": auth.SipAuthUsername,
			}).Debug("User found")
			fmt.Fprintf(w, RenderUserDirectory(auth.SipAuthUsername, addressData.Value, auth.Domain))
		}

	} else {
		log.WithFields(log.Fields{
			"domain": auth.Domain,
			"action": auth.Action,
		}).Debug("Unsupported action")

		fmt.Fprintf(w, RenderEmpty())
	}
}

func init() {
	applicationData = &AppVersion{
		Name:      AppName,
		Version:   Version,
		BuildDate: buildDate,
	}

	log.SetFormatter(&log.TextFormatter{})
	level, err := log.ParseLevel("debug")
	if err != nil {
		log.Fatal(err)
	}

	applicationLogLevel := level.String()

	log.WithFields(log.Fields{
		"buildDate":        applicationData.BuildDate,
		"Version":          applicationData.Version,
		"logLevel":         applicationLogLevel,
		"cacheTime":        cacheTime,
		"listenPort":       listenPort,
		"PapiUrl":          PapiUrl,
		"PapiUser":         PapiUser,
		"PapiPass":         "xxxxxx",
		"configPropertyId": configPropertyId,
	}).Info("Starting " + applicationData.Name)

	log.SetLevel(level)
}

func main() {
	user := []byte(BasicAuthUser)
	pass := []byte(BasicAuthPass)

	router := httprouter.New()

	router.GET("/connect-auth", BasicAuth(AuthHandler, user, pass))
	router.POST("/connect-auth", BasicAuth(AuthHandler, user, pass))
	router.GET("/", VersionHandler)
	router.GET("/version", VersionHandler)

	http.ListenAndServe(":"+listenPort, router)

}
