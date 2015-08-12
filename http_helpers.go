package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/julienschmidt/httprouter"
	"bytes"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//HostSwitch -  We need an object that implements the http.Handler interface.
// Therefore we need a type for which we implement the ServeHTTP method.
// We just use a map here, in which we map host names (with port) to http.Handlers
type HostSwitch map[string]http.Handler

// Implement the ServerHTTP method on our new type
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		// Handle host names for wich no handler is registered
		http.Error(w, "Forbidden", 403) // Or Redirect?

	}
}

// Parse request ( GET or POST ) and get freeswitch data
func parseFreeswitchRequest(req *http.Request) (a *FreeswitchRequest) {
	a = &FreeswitchRequest{}

	if req.Method == "POST" {
		a.Action = req.FormValue("action")
		a.IP = req.FormValue("ip")
		a.UserAgent = req.FormValue("sip_user_agent")
		a.Username = req.FormValue("user")
		a.Domain = req.FormValue("domain")
		a.SipAuthUsername = req.FormValue("sip_auth_username")
	} else {
		q := req.URL.Query()
		a.Action = q.Get("action")
		a.IP = q.Get("ip")
		a.UserAgent = q.Get("sip_user_agent")
		a.Domain = q.Get("domain")
		a.Username = q.Get("user")
		a.SipAuthUsername = q.Get("sip_auth_username")
	}
	return a
}

// getProvisioningStatus checks connectivity with PAPI
func getProvisioningStatus(papiURL string) bool {
	result := false
	u, err := url.Parse(papiURL)
	if err != nil {
		log.Error(err)
	}

	site := &Site{u.Host + ":80"}

	t, _ := site.Status()

	// status 2 is pass
	if t == 2 {
		result = true
	}

	fields := log.Fields{
		"status_code": t,
		"status_msg":  strconv.FormatBool(result),
	}

	if err != nil {
		fields["error"] = err
	}

	log.WithFields(fields).Debug("getProvisioningStatus()")

	return result
}

// BasicAuth - Middleware for basic auth
func BasicAuth(h httprouter.Handle, user, pass []byte) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		const basicAuthPrefix string = "Basic "

		// Get the Basic Authentication credentials
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, basicAuthPrefix) {
			// Check credentials
			payload, err := base64.StdEncoding.DecodeString(auth[len(basicAuthPrefix):])
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				if len(pair) == 2 && bytes.Equal(pair[0], user) && bytes.Equal(pair[1], pass) {
					// Delegate request to the given handle
					h(w, r, ps)
					return
				}
			}
		}

		log.Warn("Auth error: " + r.RequestURI + "")
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
