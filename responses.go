package main

import (
	"encoding/xml"
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

func RenderUserDirectorySimple(address string, secret string, domain string) string {

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
	return xml.Header + string(x)

}

func RenderUserDirectory(address string, secret string, domain string) string {

	user := &User{}
	user.Id = address

	user.Params = append(user.Params, Param{
		Name:  "password",
		Value: secret,
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
