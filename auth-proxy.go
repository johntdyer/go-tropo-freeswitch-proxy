package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	_ "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	AppName          = "tropo-auth"
	buildDate        string
	applicationData  *AppData
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

// init function does amazing things
func init() {
	applicationData = &AppData{
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

// versionRequestHandler handles incoming version / health requests
func VersionHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.WithFields(log.Fields{
		"method":     req.Method,
		"url":        req.RequestURI,
		"version":    applicationData.Version,
		"build_date": applicationData.BuildDate,
	}).Debug("VersionHandler")

	applicationData.ApiConnectivty = getProvisioningStatus(PapiUrl)
	body, _ := json.Marshal(applicationData)

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(body))
}

func AuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	authResponse := &AuthHandlerResponse{}

	freeswitch := parseFreeswitchRequest(req)

	authResponse.Fields = log.Fields{
		"remote_ip":         freeswitch.Ip,
		"sip_user_agent":    freeswitch.UserAgent,
		"domain":            freeswitch.Domain,
		"url":               req.RequestURI,
		"method":            req.Method,
		"number":            freeswitch.Username,
		"action":            freeswitch.Action,
		"sip_auth_username": freeswitch.SipAuthUsername,
	}

	// Generally we only care about sip_auth action
	if freeswitch.Action == "sip_auth" {

		// Make sure the address starts w/ a plus
		if strings.HasPrefix(freeswitch.SipAuthUsername, "+") {

			addressData := GetAddressAuthData(freeswitch.SipAuthUsername)

			// If we get no auth data back then the number is not found
			if addressData.Value == "" {
				authResponse.Message = "Address not found"
				authResponse.XmlResponse = RenderNotFound()
			} else {
				authResponse.Message = "User found"
				authResponse.XmlResponse = RenderUserDirectory(freeswitch.SipAuthUsername, addressData.Value, freeswitch.Domain)
			}

		} else {

			authResponse.Message = "Not e164 encoded"
			authResponse.XmlResponse = RenderNotFound()
		}

	} else {

		authResponse.Message = "Unsupported action"
		authResponse.XmlResponse = RenderEmpty()

	}

	log.WithFields(authResponse.Fields).Info(authResponse.Message)
	w.Header().Set("X-Tropo-Reason", authResponse.Message)
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, authResponse.XmlResponse)
}

func main() {
	user := []byte(BasicAuthUser)
	pass := []byte(BasicAuthPass)

	router := httprouter.New()

	router.GET("/connect-auth", BasicAuth(AuthHandler, user, pass))
	router.POST("/connect-auth", BasicAuth(AuthHandler, user, pass))
	router.GET("/", VersionHandler)
	router.GET("/version", VersionHandler)
	router.GET("/health", VersionHandler)

	http.ListenAndServe(":"+listenPort, router)

}
