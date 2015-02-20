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
	ConnectDomain    = "connect.tropo.com"
)

// versionRequestHandler handles incoming version / health requests
func VersionHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	body, _ := json.Marshal(applicationData)
	fmt.Fprintf(w, string(body))

}

func AuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	auth := &AuthRequest{}
	q := req.URL.Query()

	if q.Get("action") == "sip_auth" {
		auth.Domain = q.Get("domain")
		auth.Username = q.Get("user")
		auth.SipAuthUsername = q.Get("sip_auth_username")
		log.Debugf("domain: %s, user: %s, sip_auth_username: %s", auth.Domain, auth.Username, auth.SipAuthUsername)
		w.Header().Set("Content-Type", "text/xml")

		addressData := GetAddressAuthData(auth.SipAuthUsername)
		if addressData.Value == "" {
			log.Debug("Returned user not found response")
			fmt.Fprintf(w, RenderNotFound())
		} else {
			fmt.Fprintf(w, RenderUserDirectory(auth.SipAuthUsername, addressData.Value, auth.Domain))
		}

	} else {
		log.Debug("Returned user not found response")
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
	log.Info("Starting " + applicationData.Name + " Version: " + applicationData.Version + " Built on: " + applicationData.BuildDate)
	log.Warn("Logging at " + applicationLogLevel + " level")

	log.SetLevel(level)
}

func main() {
	user := []byte(BasicAuthUser)
	pass := []byte(BasicAuthPass)

	router := httprouter.New()
	router.GET("/", VersionHandler)

	router.GET("/connect-auth", BasicAuth(AuthHandler, user, pass))
	router.GET("/version", VersionHandler)
	hs := make(HostSwitch)
	hs["localhost:9082"] = router

	http.ListenAndServe(":"+listenPort, hs)

}
