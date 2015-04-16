package main

import (
	"encoding/xml"
	"strconv"
)

// Render empty XML
// <?xml version="1.0" encoding="UTF-8" standalone="no"?>
// <document type="freeswitch/xml"/>

func RenderEmpty() string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"no\"?>\n<document type=\"freeswitch/xml\"/>"

}

// Render user not found XML
// <?xml version="1.0" encoding="UTF-8" standalone="no"?>
// <document type="freeswitch/xml">
//     <section name="result">
//         <result status="not found" />
//     </section>
// </document>
func RenderNotFound() string {
	b := NotFound{
		Type: "freeswitch/xml",
		Section: SectionResult{
			Name: "result",
			Result: Result{
				Status: "not found",
			},
		},
	}

	x, err := xml.MarshalIndent(b, "", "  ")
	if err != nil {
		panic(err)
	}
	return xml.Header + string(x)
}

/*
<?xml version="1.0" encoding="UTF-8"?>
<document type="freeswitch/xml">
  <section name="directory">
    <domain name="connect.tropo.com">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"></param>
      </params>
      <variables>
        <variable name="toll_allow" value="domestic,local"/>
      </variables>
      <groups>
        <group name="default">
          <users>
            <user id="+14158510070" cacheable="300000">
              <params>
                <param name="password" value="KwDuV8LwOfoQx0EpctLFgz4O2kER"></param>
              </params>
            </user>
          </users>
        </group>
      </groups>
    </domain>
  </section>
</document>
*/

func RenderUserDirectory(address string, secret string, domain string, tollPlan string, allowDirectSipOut string) string {

	user := &User{}
	user.Id = address

	//convert config value to milliseconds
	ms, _ := strconv.Atoi(GoAuthProxy.FreeSwitchUserCacheTimeout)

	user.Cachable = strconv.Itoa(ms * 1000) //GoAuthProxy.FreeSwitchUserCacheTimeout

	user.Params = append(user.Params, Param{
		Name:  "password",
		Value: secret,
	})

	//  Add toll plan variables to Directory
	//   <variables>
	//      <variable name="toll_allow" value="domestic"></variable>
	//    </variables>
	user.Variables = append(user.Variables, Variable{
		Name:  "toll_allow",
		Value: tollPlan,
	})

	//  Based on configuration in PAPI we guard wheather or not you can place direct outbound calls
	//   <variables>
	//      <variable name="allow_direct_sip_out" value="true"></variable>
	//    </variables>
	user.Variables = append(user.Variables, Variable{
		Name:  "allow_direct_sip_out",
		Value: allowDirectSipOut,
	})

	group := &Group{Name: "default"}

	group.Users = append(group.Users, user)
	d := Document{
		Type: "freeswitch/xml",
		Section: Section{
			Name: "directory",
			Domain: Domain{
				Name: domain,
			},
		},
	}

	d.Section.Domain.Params = append(d.Section.Domain.Params, Param{
		Name:  "dial-string",
		Value: "{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}",
	})
	d.Section.Domain.Groups = append(d.Section.Domain.Groups, group)

	x, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		panic(err)
	}
	return xml.Header + string(x)
}
