package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetAddressAuthData(number string) *Auth {
	client := &http.Client{}
	resp_json := &Auth{}
	req, err := http.NewRequest("GET", PapiUrl+"/addresses/number/"+number+"/config/"+configPropertyId, nil)
	req.SetBasicAuth(PapiUser, PapiPass)
	api_resp, err := client.Do(req)
	if err != nil {
		log.Error("PAPI Error : %s", err)
		return resp_json
	}

	fields := log.Fields{
		"responseCode": api_resp.StatusCode,
		"number":       number,
	}

	if api_resp.StatusCode == 200 {

		body, _ := ioutil.ReadAll(api_resp.Body)
		resp_json.Address = number
		log.WithFields(fields).Debug("PAPI Request successful")

		err := json.Unmarshal(body, &resp_json)
		if err != nil {
			log.Error("PAPI Error : %s", err)
			return resp_json
		}

	} else {

		log.WithFields(fields).Warn("PAPI request returned non-2xx")

	}
	return resp_json
}
