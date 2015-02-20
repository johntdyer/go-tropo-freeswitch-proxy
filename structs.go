package main

import (
	"encoding/xml"
)

// AppVersion contains information about the version of the application
type AppVersion struct {
	Name      string `json:"application"`
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
}

type Auth struct {
	Href        string `json:"href"`
	Id          string `json:"id"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Address     string `json:"address"`
}

type AuthRequest struct {
	Domain          string
	Username        string
	SipAuthUsername string
}

type Address struct {
	XMLName xml.Name `xml:"domain"`
	Name    string   `xml:"name,attr"`
	User    *User    `xml:"user"`
}

type Document struct {
	XMLName xml.Name `xml:"document"`
	Type    string   `xml:"type,attr"`
	Section Section  `xml:"section"`
}

type Section struct {
	Name   string `xml:"name,attr"`
	Domain Domain `xml:"domain"`
}

type SectionResult struct {
	Name   string `xml:"name,attr"`
	Result Result `xml:"result"`
}
type Domain struct {
	Name   string   `xml:"name,attr"`
	Params []Param  `xml:"params>param"`
	Groups []*Group `xml:"groups>group"`
}

type Group struct {
	Name  string  `xml:"name,attr"`
	Users []*User `xml:"users>user"`
}

type User struct {
	Id       string  `xml:"id,attr"`
	Cachable bool    `xml:"cacheable,attr"`
	Params   []Param `xml:"params>param"`
}

type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type NotFound struct {
	XMLName xml.Name      `xml:"document"`
	Type    string        `xml:"type,attr"`
	Section SectionResult `xml:"result"`
}

type Result struct {
	Status string `xml:"status,attr"`
}
