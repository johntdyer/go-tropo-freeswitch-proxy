package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/fukata/golang-stats-api-handler"
	_ "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/joho/godotenv/autoload"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/pmylund/go-cache"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/thoas/stats"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"strings"
	"time"
)

var (
	buildDate string
	//GoAuthProxy - Primary config
	GoAuthProxy *GoAuthProxyConfig
)

// init function does amazing things
func init() {

	expiredCachePurgeInterval, _ := strconv.Atoi(os.Getenv("EXPIRED_CACHE_PURGE_INTERVAL"))
	userCacheTTL, _ := strconv.Atoi(os.Getenv("APP_CACHE_TTL"))
	appCacheNegativeTTL, _ := strconv.Atoi(os.Getenv("APP_CACHE_NEGATIVE_TTL"))

	// Cache durations
	CacheTTLDuration := time.Duration(rand.Int31n(int32(userCacheTTL)))
	CachePurgeDuration := time.Duration(rand.Int31n(int32(expiredCachePurgeInterval)))

	// Do we validate the domain defined in ENV['CONNECT_DOMAIN'] against the request.
	validateDomain, _ := strconv.ParseBool(os.Getenv("VALIDATE_CONNECT_DOMAIN"))
	if validateDomain {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	GoAuthProxy = &GoAuthProxyConfig{
		LogLevel:                   os.Getenv("LOG_LEVEL"),
		ListenPort:                 os.Getenv("LISTEN_PORT"),
		PapiUser:                   os.Getenv("TROPO_API_USER"),
		PapiPass:                   os.Getenv("TROPO_API_PASS"),
		PapiURL:                    os.Getenv("TROPO_API_URL"),
		ConnectDomain:              os.Getenv("CONNECT_DOMAIN"),
		BasicAuthUser:              os.Getenv("API_AUTH_USER"),
		BasicAuthPass:              os.Getenv("API_AUTH_PASS"),
		FreeSwitchUserCacheTimeout: os.Getenv("FREESWITCH_CACHE_TIMEOUT"),
		DefaultTollPlan:            os.Getenv("DEFAULT_TOLL_PLAN"),
		ValidateDomain:             validateDomain,
		CacheTTL:                   userCacheTTL,
		CacheNegativeTTL:           appCacheNegativeTTL,
		ExpiredCachePurgeInterval:  expiredCachePurgeInterval,
		Version:                    Version,
		BuildDate:                  buildDate,
		AppName:                    "tropo-auth",
	}

	// Create a cache with a default expiration time of T minutes, and which  purges expired items every N seconds
	GoAuthProxy.AuthProxyCache = cache.New(CacheTTLDuration*time.Second, CachePurgeDuration*time.Second)

	log.SetFormatter(&log.TextFormatter{})
	level, err := log.ParseLevel(GoAuthProxy.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"buildDate":                  GoAuthProxy.BuildDate,
		"Version":                    GoAuthProxy.Version,
		"logLevel":                   GoAuthProxy.LogLevel,
		"ProxyCacheTTL":              GoAuthProxy.CacheTTL,
		"ProxyCacheNegativeTTL":      GoAuthProxy.CacheNegativeTTL,
		"ExpiredCachePurgeInterval":  GoAuthProxy.ExpiredCachePurgeInterval,
		"freeSwitchUserCacheTimeout": GoAuthProxy.FreeSwitchUserCacheTimeout,
		"listenPort":                 GoAuthProxy.ListenPort,
		"PapiURL":                    GoAuthProxy.PapiURL,
		"defaultTollPlan":            GoAuthProxy.DefaultTollPlan,
		"connectDomain":              GoAuthProxy.ConnectDomain,
		"validateDomain":             GoAuthProxy.ValidateDomain,
		"PapiUser":                   GoAuthProxy.PapiUser,
		"PapiPass":                   "xxxxxx",
	}).Info("Starting " + GoAuthProxy.AppName)

	log.SetLevel(level)
}

// VersionHandler handles incoming version / health requests
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

	ad.APIConnectivty = getProvisioningStatus(GoAuthProxy.PapiURL)
	body, _ := json.Marshal(ad)

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(body))
}

// DirectoryAuthHandler - HTTP handler for directory requests
func DirectoryAuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	authResponse := &DirectoryAuthResponse{}

	freeswitch := parseFreeswitchRequest(req)

	authResponse.Fields = log.Fields{
		"remote_ip":         freeswitch.IP,
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

		// Validate domain and make sure its supported.  This is intended to help out w/ asshat's trying to register against FS
		if GoAuthProxy.ValidateDomain == true && freeswitch.Domain != GoAuthProxy.ConnectDomain {

			authResponse.Fields["valid_domain"] = GoAuthProxy.ConnectDomain
			log.WithFields(authResponse.Fields).Warn("Domain does not match. Expected " + GoAuthProxy.ConnectDomain)

			w.Header().Set("X-Tropo-Reason", "Domain does not match. Expected "+GoAuthProxy.ConnectDomain)
			w.Header().Set("Content-Type", "text/xml")

			fmt.Fprintf(w, RenderNotFound())
			return
		}

		// Make sure the address starts w/ a plus
		if strings.HasPrefix(freeswitch.SipAuthUsername, "+") {

			configData := GetAddressConfigData(freeswitch.SipAuthUsername)
			// If we get no auth data back then the number is not found
			if configData.Name["com.tropo.connect.address.secret"] == "" {
				authResponse.Message = "Address not found"
				authResponse.XMLResponse = RenderNotFound()
			} else {
				authResponse.Message = "User found"

				//If the tollplan config is set we'll use that otherwise we'll use the default
				tollPlan := ""
				if configData.Name["com.tropo.connect.tollAllow"] == "" {
					tollPlan = GoAuthProxy.DefaultTollPlan
				} else {
					tollPlan = configData.Name["com.tropo.connect.tollAllow"]
				}

				allowDirectSipOut := configData.Name["com.tropo.connect.sip_outbound_allow"]

				authResponse.XMLResponse = RenderUserDirectory(freeswitch.SipAuthUsername,
					configData.Name["com.tropo.connect.address.secret"],
					freeswitch.Domain,
					tollPlan,
					allowDirectSipOut,
				)

			}

		} else {

			authResponse.Message = "Not e164 encoded"
			authResponse.XMLResponse = RenderNotFound()

		}
	} else {

		authResponse.Message = "Unsupported action"
		authResponse.XMLResponse = RenderEmpty()

	}

	log.WithFields(authResponse.Fields).Info(authResponse.Message)
	w.Header().Set("X-Tropo-Reason", authResponse.Message)
	w.Header().Set("Content-Type", "text/xml")
	fmt.Fprintf(w, authResponse.XMLResponse)
}

// UserAuthHandler - HTTP handler for authenticating users
func UserAuthHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res := &UserAuthHandlerResponse{}
	res.Address = ps.ByName("address")

	res.Fields = log.Fields{
		"url":     req.RequestURI,
		"method":  req.Method,
		"address": res.Address,
		"cache":   nil,
	}
	result, found := GoAuthProxy.AuthProxyCache.Get(res.Address)

	// Address found in cache
	if found {
		if result == true {
			res.Fields["cache"] = "hit"
			res.Message = "Cached address found"
			res.Header = http.StatusNonAuthoritativeInfo
		} else {
			res.Fields["cache"] = "miss"
			res.Message = "Cached address not found"
			res.Header = http.StatusNotFound
		}

	} else {

		if strings.HasPrefix(res.Address, "+") {

			configData := GetAddressConfigData(res.Address)
			// If we get no auth data back then the number is not found
			if configData.Name["com.tropo.connect.address.secret"] == "" {

				GoAuthProxy.AuthProxyCache.Set(res.Address, false, time.Duration(rand.Int31n(int32(GoAuthProxy.CacheNegativeTTL)))*time.Second)
				res.Fields["cache"] = "set-negative"
				res.Message = "Address not found"
				res.Header = http.StatusNotFound
			} else {
				GoAuthProxy.AuthProxyCache.Set(res.Address, true, cache.DefaultExpiration)
				res.Fields["cache"] = "set-positive"
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

// CacheHandler - Middleware for caching
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

// StatsHandler - Allow httprouter to call  stats_api.Handler
func StatsHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	stats_api.Handler(w, req)
}

func main() {
	user := []byte(GoAuthProxy.BasicAuthUser)
	pass := []byte(GoAuthProxy.BasicAuthPass)

	//http stats
	stats := stats.New()

	router := httprouter.New()

	router.GET("/stats/go", StatsHandler)
	router.GET("/connect-auth", BasicAuth(DirectoryAuthHandler, user, pass))
	router.POST("/connect-auth", BasicAuth(DirectoryAuthHandler, user, pass))
	router.GET("/", VersionHandler)
	router.GET("/version", VersionHandler)
	router.GET("/health", VersionHandler)
	router.DELETE("/cache", CacheHandler)
	router.GET("/cache", CacheHandler)
	router.GET("/users/:address", UserAuthHandler)

	router.GET("/stats/http", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s, err := json.Marshal(stats.Data())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(s)
	})

	http.ListenAndServe(":"+GoAuthProxy.ListenPort, stats.Handler(router))

}
