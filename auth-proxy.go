package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	_ "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/pmylund/go-cache"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	buildDate   string
	GoAuthProxy *GoAuthProxyConfig
)

// init function does amazing things
func init() {
	cacheTimeout, _ := strconv.Atoi(os.Getenv("PAPI_CACHE"))
	GoAuthProxy = &GoAuthProxyConfig{
		AppName:         "tropo-auth",
		LogLevel:        os.Getenv("LOG_LEVEL"),
		AppCacheTimeout: cacheTimeout,
		AuthProxyCache:  cache.New(2*time.Minute, 30*time.Second),
		ListenPort:      os.Getenv("LISTEN_PORT"),
		PropertyId:      os.Getenv("ADDRESS_CONFIG_PROPERTY_ID"),
		CacheTime:       os.Getenv("USER_CACHE_VALUE"),
		PapiUser:        os.Getenv("TROPO_API_USER"),
		PapiPass:        os.Getenv("TROPO_API_PASS"),
		PapiUrl:         os.Getenv("TROPO_API_URL"),
		BasicAuthUser:   os.Getenv("API_AUTH_USER"),
		BasicAuthPass:   os.Getenv("API_AUTH_PASS"),
		Version:         Version,
		BuildDate:       buildDate,
	}

	log.SetFormatter(&log.TextFormatter{})
	level, err := log.ParseLevel(GoAuthProxy.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"buildDate":        GoAuthProxy.BuildDate,
		"Version":          GoAuthProxy.Version,
		"logLevel":         GoAuthProxy.LogLevel,
		"cacheTime":        GoAuthProxy.CacheTime,
		"listenPort":       GoAuthProxy.ListenPort,
		"PapiUrl":          GoAuthProxy.PapiUrl,
		"configPropertyId": GoAuthProxy.PropertyId,
		"PapiUser":         GoAuthProxy.PapiUser,
		"PapiPass":         "xxxxxx",
	}).Info("Starting " + GoAuthProxy.AppName)

	log.SetLevel(level)
}

// versionRequestHandler handles incoming version / health requests
func VersionHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ad := AppData{
		Name:      GoAuthProxy.AppName,
		BuildDate: GoAuthProxy.BuildDate,
		Version:   GoAuthProxy.Version,
	}

	log.WithFields(log.Fields{
		"method":     req.Method,
		"url":        req.RequestURI,
		"version":    ad.Version,
		"build_date": ad.BuildDate,
	}).Debug("VersionHandler")

	ad.ApiConnectivty = getProvisioningStatus(GoAuthProxy.PapiUrl)
	body, _ := json.Marshal(ad)

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(body))
}

func DirectoryAuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	authResponse := &DirectoryAuthResponse{}

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

func UserAuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := &UserAuthHandlerResponse{}
	res.Fields = log.Fields{
		"url":     req.RequestURI,
		"method":  req.Method,
		"address": res.Address,
	}
	res.Address = ps.ByName("address")
	result, found := GoAuthProxy.AuthProxyCache.Get(res.Address)

	// Address found in cache
	if found {
		if result == true {
			res.Message = "Cached address found"
			res.Header = http.StatusNonAuthoritativeInfo
		} else {
			res.Message = "Cached address not found"
			res.Header = http.StatusNotFound
		}

	} else {

		if strings.HasPrefix(res.Address, "+") {

			addressData := GetAddressAuthData(res.Address)

			if addressData.Value == "" {
				GoAuthProxy.AuthProxyCache.Set(res.Address, false, 1*time.Minute)

				res.Message = "Address not found"
				res.Header = http.StatusNotFound
			} else {
				GoAuthProxy.AuthProxyCache.Set(res.Address, true, cache.DefaultExpiration)
				res.Message = "Address found"
				res.Header = http.StatusNonAuthoritativeInfo
			}

		} else {
			res.Message = "Missing plus"
			res.Header = http.StatusBadRequest
		}
	}

	log.WithFields(res.Fields).Debug("UserAuthHandler()")

	w.Header().Set("X-Tropo-Lookup-Result", res.Message)
	w.WriteHeader(res.Header)

}
func CacheHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	log.WithFields(log.Fields{
		"url":    req.RequestURI,
		"method": req.Method,
	}).Debug("CacheHandler()")

	msg := ""

	if req.Method == "DELETE" {
		GoAuthProxy.AuthProxyCache.Flush()
		msg = "Flushed"

	} else {
		msg = fmt.Sprintf("Cache count: %s", strconv.Itoa(GoAuthProxy.AuthProxyCache.ItemCount()))

	}
	fmt.Fprintf(w, msg)
}

func main() {
	user := []byte(GoAuthProxy.BasicAuthUser)
	pass := []byte(GoAuthProxy.BasicAuthPass)

	router := httprouter.New()

	router.GET("/connect-auth", BasicAuth(DirectoryAuthHandler, user, pass))
	router.POST("/connect-auth", BasicAuth(DirectoryAuthHandler, user, pass))
	router.GET("/", VersionHandler)
	router.GET("/version", VersionHandler)
	router.GET("/health", VersionHandler)
	router.DELETE("/cache", CacheHandler)
	router.GET("/cache", CacheHandler)
	router.GET("/users/:address", UserAuthHandler)

	http.ListenAndServe(":"+GoAuthProxy.ListenPort, router)

}
