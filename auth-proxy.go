package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	_ "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"encoding/json"
	"encoding/xml"
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
	ConnectDomain    = "connect.tropo.com"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func AuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	q := req.URL.Query()

	if q.Get("action") == "sip_auth" {
		domain := q.Get("domain")
		username := q.Get("user")
		sip_auth_username := q.Get("sip_auth_username")
		log.Debugf("domain: %s, user: %s, sip_auth_username: %s", domain, username, sip_auth_username)
		fmt.Fprintf(w, "domain: %s, user: %s, sip_auth_username: %s\n", domain, username, sip_auth_username)
	} else {

		fmt.Fprintf(w, "wtf")
	}

	// addressData := getAuth(ps.ByName("name"))
	// fmt.Printf("%+v\n", addressData)
	// number := fmt.Sprintf("+%s", addressData.Address)
	// fmt.Fprintf(w, renderUser(number, addressData.Value, ConnectDomain))
	//
}

func renderUser(address string, secret string, domain string) string {
	user := &User{}
	user.Id = address

	user.Params = append(user.Params, Param{
		Name:  "password",
		Value: secret,
	})

	profile := Address{
		Name: domain,
		User: user,
	}
	x, _ := xml.MarshalIndent(profile, "", "  ")
	return string(x)
}

func getAuth(number string) *Auth {
	client := &http.Client{}

	req, err := http.NewRequest("GET", PapiUrl+"/addresses/number/+"+number+"/config/"+configPropertyId, nil)
	req.SetBasicAuth(PapiUser, PapiPass)
	api_resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	body, _ := ioutil.ReadAll(api_resp.Body)

	resp_json := &Auth{
		Address: number,
	}

	json.Unmarshal(body, &resp_json)
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
	hs["localhost:9082"] = router

	http.ListenAndServe("localhost:9082", hs)

}
