package main

import (
	log "bitbucket.org/voxeolabs/go-freeswitch-auth-proxy/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// fetchProvisioningAddressConfig - Return struct from PAPI of address config data
func fetchProvisioningAddressConfig(address string) (*PapiResponse, error) {

	client := &http.Client{}
	respJSON := &PapiResponse{}

	req, err := http.NewRequest("GET", GoAuthProxy.PapiURL+"/addresses/number/"+address+"/config", nil)
	req.Close = true
	req.SetBasicAuth(GoAuthProxy.PapiUser, GoAuthProxy.PapiPass)
	resp, err := client.Do(req)

	// defer close
	defer resp.Body.Close()

	if err != nil {
		log.Error("PAPI Error : %s", err)
		return respJSON, err
	}

	fields := log.Fields{
		"responseCode": resp.StatusCode,
		"address":      address,
	}

	if resp.StatusCode == 200 {

		body, _ := ioutil.ReadAll(resp.Body)
		respJSON.Address = address

		log.WithFields(fields).Debug("PAPI Request successful")

		err := json.Unmarshal(body, &respJSON.Configs)
		if err != nil {
			log.Error("PAPI Error : %s", err)
			return respJSON, err
		}

	} else {

		log.WithFields(fields).Warn("PAPI request returned non-2xx")

	}
	return respJSON, nil

}

//setConfigDataStruct returns a stuct with maps of config data. One is indexed based on ID, the other is indexes on name
//  cd.ID[16] = "abc124"
//  cd.Name["com.tropo.connect.address.secret"] = "abc123"
func setConfigDataStruct(papi *PapiResponse) *ConfigData {
	cd := &ConfigData{
		ID:   make(map[int]string),
		Name: make(map[string]string),
	}

	// Add to map using Config.ID as key
	for _, config := range papi.Configs {
		idx, err := strconv.Atoi(config.ID)
		if err != nil {
			log.Error(err)
			return cd
		}
		cd.ID[idx] = config.Value
		cd.Name[config.Name] = config.Value
	}

	return cd

}

// GetAddressConfigData - Get address config data
func GetAddressConfigData(number string) *ConfigData {
	papiData, err := fetchProvisioningAddressConfig(number)
	if err != nil {
		log.WithFields(log.Fields{
			"method": "GetAddressConfigData",
		}).Error(err)
	}
	return setConfigDataStruct(papiData)
}
