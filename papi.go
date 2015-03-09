package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// getAddressConfig - Return struct from PAPI of address config data
func fetchProvisioningAddressConfig(address string) (*PapiResponse, error) {

	client := &http.Client{}
	resp_json := &PapiResponse{}

	req, err := http.NewRequest("GET", GoAuthProxy.PapiUrl+"/addresses/number/"+address+"/config", nil)
	req.SetBasicAuth(GoAuthProxy.PapiUser, GoAuthProxy.PapiPass)
	api_resp, err := client.Do(req)
	if err != nil {
		log.Error("PAPI Error : %s", err)
		return resp_json, err
	}

	fields := log.Fields{
		"responseCode": api_resp.StatusCode,
		"address":      address,
	}

	if api_resp.StatusCode == 200 {

		body, _ := ioutil.ReadAll(api_resp.Body)
		resp_json.Address = address

		log.WithFields(fields).Debug("PAPI Request successful")

		err := json.Unmarshal(body, &resp_json.Configs)
		if err != nil {
			log.Error("PAPI Error : %s", err)
			return resp_json, err
		}

	} else {

		log.WithFields(fields).Warn("PAPI request returned non-2xx")

	}
	return resp_json, nil

}

//setConfigDataStruct returns a stuct with maps of config data. One is indexed based on ID, the other is indexes on name
//  cd.Id[16] = "abc124"
//  cd.Name["com.tropo.connect.address.secret"] = "abc123"
func setConfigDataStruct(papi *PapiResponse) *ConfigData {
	cd := &ConfigData{
		Id:   make(map[int]string),
		Name: make(map[string]string),
	}

	// Add to map using Config.Id as key
	for _, config := range papi.Configs {
		idx, err := strconv.Atoi(config.Id)
		if err != nil {
			log.Error(err)
			return cd
		}
		cd.Id[idx] = config.Value
		cd.Name[config.Name] = config.Value
	}

	return cd

}

func GetAddressConfigData(number string) *ConfigData {
	papiData, err := fetchProvisioningAddressConfig(number)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "GetAddressConfigData",
		}).Error(err)
	}
	return setConfigDataStruct(papiData)
}
