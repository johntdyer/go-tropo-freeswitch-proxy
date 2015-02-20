# Freeswitch auth proxy

# Usage


### Development

    godep go run auth-proxy.go http_helpers.go structs.go responses.go


#### Examples


###### POST

```log
$ curl -v -X POST -u tropo:tropo "http://127.0.0.1:9082/connect-auth" -d "hostname=fs1&section=directory&tag_name=domain&key_name=name&key_value=connect.tropo.com&Event-Name=REQUEST_PARAMS&Core-UUID=0efa6954-cb1f-4b7c-b280-ba3434907a74&FreeSWITCH-Hostname=fs1&FreeSWITCH-Switchname=fs1&FreeSWITCH-IPv4=198.11.254.113&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2015-02-20%2014%3A08%3A45&Event-Date-GMT=Fri,%2020%20Feb%202015%2020%3A08%3A45%20GMT&Event-Date-Timestamp=1424462925448233&Event-Calling-File=sofia_reg.c&Event-Calling-Function=sofia_reg_parse_auth&Event-Calling-Line-Number=2741&Event-Sequence=518&action=sip_auth&sip_profile=internal&sip_user_agent=Blink%20Pro%204.1.0%20(MacOSX)&sip_auth_username=%2B14158510070&sip_auth_realm=connect.tropo.com&sip_auth_nonce=7cff7d80-afbb-4e0c-9455-c0d9d6baa3f7&sip_auth_uri=sip%3Aconnect.tropo.com&sip_contact_user=65341890&sip_contact_host=24.126.169.35&sip_to_user=%2B14158510070&sip_to_host=connect.tropo.com&sip_via_protocol=udp&sip_from_user=%2B14158510070&sip_from_host=connect.tropo.com&sip_call_id=XIiqYI.uxMNatGF6m76zb6HEFe0skr46&sip_request_host=connect.tropo.com&sip_auth_qop=auth&sip_auth_cnonce=mSfOYudMXnSRJsLwnrKiWJvhrChlrV00&sip_auth_nc=00000001&sip_auth_response=bcf3664ef7e5a84578967dde61d72c26&sip_auth_method=REGISTER&client_port=56510&key=id&user=%2B14158510070&domain=connect.tropo.com&ip=24.126.169.35"
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 9082 (#0)
* Server auth using Basic with user 'tropo'
> POST /connect-auth HTTP/1.1
> Authorization: Basic dHJvcG86dHJvcG8=
> User-Agent: curl/7.37.1
> Host: 127.0.0.1:9082
> Accept: */*
> Content-Length: 1296
> Content-Type: application/x-www-form-urlencoded
> Expect: 100-continue
>
< HTTP/1.1 100 Continue
< HTTP/1.1 200 OK
< Content-Type: text/xml
< Date: Fri, 20 Feb 2015 21:20:39 GMT
< Content-Length: 678
<
<?xml version="1.0" encoding="UTF-8"?>
<document type="freeswitch/xml">
  <section name="directory">
    <domain name="connect.tropo.com">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"></param>
      </params>
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
* Connection #0 to host 127.0.0.1 left intact
</document>%
# jdyer at retina.dyer in ~/.gvm/pkgsets/go1.4/global/src/bitbucket.org/voxeolabs/go-freeswitch-auth-proxy on git:master x [16:20:39]
$
```

###### GET

```log
$ curl -v -X POST -u tropo:tropo "http://127.0.0.1:9082/connect-auth?hostname=fs1&section=directory&tag_name=domain&key_name=name&key_value=connect.tropo.com&Event-Name=REQUEST_PARAMS&Core-UUID=0efa6954-cb1f-4b7c-b280-ba3434907a74&FreeSWITCH-Hostname=fs1&FreeSWITCH-Switchname=fs1&FreeSWITCH-IPv4=198.11.254.113&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2015-02-20%2014%3A08%3A45&Event-Date-GMT=Fri,%2020%20Feb%202015%2020%3A08%3A45%20GMT&Event-Date-Timestamp=1424462925448233&Event-Calling-File=sofia_reg.c&Event-Calling-Function=sofia_reg_parse_auth&Event-Calling-Line-Number=2741&Event-Sequence=518&action=sip_auth&sip_profile=internal&sip_user_agent=Blink%20Pro%204.1.0%20(MacOSX)&sip_auth_username=%2B14158510070&sip_auth_realm=connect.tropo.com&sip_auth_nonce=7cff7d80-afbb-4e0c-9455-c0d9d6baa3f7&sip_auth_uri=sip%3Aconnect.tropo.com&sip_contact_user=65341890&sip_contact_host=24.126.169.35&sip_to_user=%2B14158510070&sip_to_host=connect.tropo.com&sip_via_protocol=udp&sip_from_user=%2B14158510070&sip_from_host=connect.tropo.com&sip_call_id=XIiqYI.uxMNatGF6m76zb6HEFe0skr46&sip_request_host=connect.tropo.com&sip_auth_qop=auth&sip_auth_cnonce=mSfOYudMXnSRJsLwnrKiWJvhrChlrV00&sip_auth_nc=00000001&sip_auth_response=bcf3664ef7e5a84578967dde61d72c26&sip_auth_method=REGISTER&client_port=56510&key=id&user=%2B14158510070&domain=connect.tropo.com&ip=24.126.169.35"
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 9082 (#0)
* Server auth using Basic with user 'tropo'
> POST /connect-auth?hostname=fs1&section=directory&tag_name=domain&key_name=name&key_value=connect.tropo.com&Event-Name=REQUEST_PARAMS&Core-UUID=0efa6954-cb1f-4b7c-b280-ba3434907a74&FreeSWITCH-Hostname=fs1&FreeSWITCH-Switchname=fs1&FreeSWITCH-IPv4=198.11.254.113&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2015-02-20%2014%3A08%3A45&Event-Date-GMT=Fri,%2020%20Feb%202015%2020%3A08%3A45%20GMT&Event-Date-Timestamp=1424462925448233&Event-Calling-File=sofia_reg.c&Event-Calling-Function=sofia_reg_parse_auth&Event-Calling-Line-Number=2741&Event-Sequence=518&action=sip_auth&sip_profile=internal&sip_user_agent=Blink%20Pro%204.1.0%20(MacOSX)&sip_auth_username=%2B14158510070&sip_auth_realm=connect.tropo.com&sip_auth_nonce=7cff7d80-afbb-4e0c-9455-c0d9d6baa3f7&sip_auth_uri=sip%3Aconnect.tropo.com&sip_contact_user=65341890&sip_contact_host=24.126.169.35&sip_to_user=%2B14158510070&sip_to_host=connect.tropo.com&sip_via_protocol=udp&sip_from_user=%2B14158510070&sip_from_host=connect.tropo.com&sip_call_id=XIiqYI.uxMNatGF6m76zb6HEFe0skr46&sip_request_host=connect.tropo.com&sip_auth_qop=auth&sip_auth_cnonce=mSfOYudMXnSRJsLwnrKiWJvhrChlrV00&sip_auth_nc=00000001&sip_auth_response=bcf3664ef7e5a84578967dde61d72c26&sip_auth_method=REGISTER&client_port=56510&key=id&user=%2B14158510070&domain=connect.tropo.com&ip=24.126.169.35 HTTP/1.1
> Authorization: Basic dHJvcG86dHJvcG8=
> User-Agent: curl/7.37.1
> Host: 127.0.0.1:9082
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: text/xml
< Date: Fri, 20 Feb 2015 21:21:10 GMT
< Content-Length: 678
<
<?xml version="1.0" encoding="UTF-8"?>
<document type="freeswitch/xml">
  <section name="directory">
    <domain name="connect.tropo.com">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"></param>
      </params>
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
* Connection #0 to host 127.0.0.1 left intact
</document>%
# jdyer at retina.dyer in ~/.gvm/pkgsets/go1.4/global/src/bitbucket.org/voxeolabs/go-freeswitch-auth-proxy on git:master x [16:21:10]
$
```
