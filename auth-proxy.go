package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	_ "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"encoding/json"
	// "encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	configPropertyId = os.Getenv("ADDRESS_CONFIG_PROPERTY_ID")
	PapiUser         = os.Getenv("TROPO_API_USER")
	PapiPass         = os.Getenv("TROPO_API_PASS")
	PapiUrl          = os.Getenv("TROPO_API_URL")
	BasicAuthUser    = os.Getenv("API_AUTH_USER")
	BasicAuthPass    = os.Getenv("API_AUTH_PASS")
	listenPort       = os.Getenv("LISTEN_PORT")
	ConnectDomain    = "connect.tropo.com"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
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

		addressData := getAuth(auth.SipAuthUsername)
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

	//
}

func getAuth(number string) *Auth {
	client := &http.Client{}
	resp_json := &Auth{}
	req, err := http.NewRequest("GET", PapiUrl+"/addresses/number/"+number+"/config/"+configPropertyId, nil)
	req.SetBasicAuth(PapiUser, PapiPass)
	api_resp, err := client.Do(req)
	if err != nil {
		log.Error("Error : %s", err)
		return resp_json
	}

	if api_resp.StatusCode == 200 {

		body, _ := ioutil.ReadAll(api_resp.Body)

		resp_json.Address = number

		json.Unmarshal(body, &resp_json)

	} else {
		log.Warn("user not found")
	}
	return resp_json
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	level, err := log.ParseLevel("debug")
	if err != nil {
		log.Fatal(err)
	}
	applicationLogLevel := level.String()
	log.Warn("Logging at " + applicationLogLevel + " level")
	log.SetLevel(level)
}

func main() {
	user := []byte(BasicAuthUser)
	pass := []byte(BasicAuthPass)

	router := httprouter.New()
	router.GET("/", Index)

	router.GET("/connect-auth", BasicAuth(AuthHandler, user, pass))

	hs := make(HostSwitch)
	hs[":9082"] = router

	http.ListenAndServe(":"+listenPort, hs)

}
