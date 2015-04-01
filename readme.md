# Freeswitch auth proxy

# Usage

* Configure .env file

`cp .env-example .env`

* Build

`./make.sh`

* Run

`./builds/tropo-auth-0.0.1.osx`

### Development

```shell
go get github.com/tools/godep
godep go run auth-proxy.go http_helpers.go check_api.go structs.go responses.go version.go papi.go
```


#### Installation / Config


To use this plugin you must made a few config changes to Freeswitch.

* Install mod-curl

`apt-get install freeswitch-mod-xml-curl`

* Enable `mod_xml_curl` in `/etc/freeswitch/autoload_configs/modules.conf.xml`. This is done by simply uncommenting it

* Suggest you disable mod_voicemail `<!--load module="mod_voicemail"/-->` inside `/etc/freeswitch/autoload_configs/modules.conf.xml`

* Configure Freeswitch to use curl for directory lookups `autoload_configs/modules.conf.xml`

```xml
<configuration description="cURL XML Gateway" name="xml_curl.conf">
  <bindings>
    <binding name="directory">
      <param bindings="directory" name="gateway-url" value="http://127.0.0.1:9082/connect-auth"/>
      <param name="gateway-credentials" value="user:pass"/>
      <param name="auth-scheme" value="basic"/>
    </binding>
  </bindings>
</configuration>
```

##### Application Config

The application gets it config from a .env file.  The config items are as follows


* `TROPO_API_USER`               - Username for user with read access to address config data
* `TROPO_API_PASS`               - Password for user with read access to address config data
* `TROPO_API_URL`                - Provisioning API URL
* `API_AUTH_USER`                - Basic auth user for making requests against proxy
* `API_AUTH_PASS`                - Basic auth pass for making requests against proxy
* `LISTEN_PORT`                  - Listen port
* `LOG_LEVEL`                    - Log level
* `DEFAULT_TOLL_PLAN`            - Plan to use if there is none found in PAPI
* `APP_CACHE_TTL`                - How long to cache application lookups from PAPI
* `APP_CACHE_NEGATIVE_TTL`       - How long to cache negative lookups
* `FREESWITCH_CACHE_TIMEOUT`     - How long freeswich should cache directoy lookups
* `EXPIRED_CACHE_PURGE_INTERVAL` - How often expired items are purged from cache
* `CONNECT_DOMAIN`               - Connect domain
* `VALIDATE_CONNECT_DOMAIN`      - validate requests against `CONNECT_DOMAIN`

##### Upstart

Using upstart is a simple way to start the service.  The following was tested on Ubuntu 14.04

```shell
sudo apt-get install supervisor
sudo addgroup --system supervisor
sudo adduser tropo supervisor


cat <<EOF > /etc/supervisor/supervisord.conf
[unix_http_server]
file=/var/run/supervisor.sock
chmod=0770
chown=root:supervisor

[supervisord]
logfile=/var/log/supervisor/supervisord.log
pidfile=/var/run/supervisord.pid
childlogdir=/var/log/supervisor

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///var/run/supervisor.sock

[include]
files = /etc/supervisor/conf.d/*.conf
EOF

cat <<EOF > /etc/supervisord/conf.d/tropo-auth-proxy.conf
[program:tropo-auth-proxy]
command=/opt/tropo-auth-proxy/bin/tropo-auth.linux
autostart=true
autorestart=true
startretries=10
user=tropo
directory=/opt/tropo-auth-proxy
redirect_stderr=true
stdout_logfile=/opt/tropo-auth-proxy/logs/tropo-auth-proxy.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
EOF

sudo service supervisor restart
supervisorctl status tropo-auth-proxy
```


#### Examples


###### POST

```log
$ curl -v -X POST -u user:pass "http://127.0.0.1:9082/connect-auth" -d "hostname=fs1&section=directory&tag_name=domain&key_name=name&key_value=connect.example.com&Event-Name=REQUEST_PARAMS&Core-UUID=cb1f-4b7c-b280-ba3434907a74&FreeSWITCH-Hostname=fs1&FreeSWITCH-Switchname=fs1&FreeSWITCH-IPv4=198.11.123.123&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2015-02-20%2014%3A08%3A45&Event-Date-GMT=Fri,%2020%20Feb%202015%2020%3A08%3A45%20GMT&Event-Date-Timestamp=1424462925448233&Event-Calling-File=sofia_reg.c&Event-Calling-Function=sofia_reg_parse_auth&Event-Calling-Line-Number=2741&Event-Sequence=518&action=sip_auth&sip_profile=internal&sip_user_agent=Blink%20Pro%204.1.0%20(MacOSX)&sip_auth_username=%2B111238510000&sip_auth_realm=connect.example.com&sip_auth_nonce=7cff7d80-afbb-4e0c-9455-c0d9d6baa3f7&sip_auth_uri=sip%3Aconnect.example.com&sip_contact_user=65341890&sip_contact_host=24.126.169.35&sip_to_user=%2B141238510000&sip_to_host=connect.example.com&sip_via_protocol=udp&sip_from_user=%2B141238510000&sip_from_host=connect.example.com&sip_call_id=XIiqYI.uxMNatGF6m76zb6HEFe0skr46&sip_request_host=connect.example.com&sip_auth_qop=auth&sip_auth_cnonce=mSfOYudMXnSRJsLwnrKiWJvhrChlrV00&sip_auth_nc=00000001&sip_auth_response=bcf3664ef7e5a84578967dde61d72c26&sip_auth_method=REGISTER&client_port=56510&key=id&user=%2B141238510000&domain=connect.example.com&ip=24.126.169.35"
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
    <domain name="connect.example.com">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"></param>
      </params>
      <groups>
        <group name="default">
          <users>
            <user id="+141238510000" cacheable="300000">
              <params>
                <param name="password" value="abc123"></param>
              </params>
            </user>
          </users>
        </group>
      </groups>
    </domain>
  </section>
* Connection #0 to host 127.0.0.1 left intact
</document>%
$
```

###### GET

```log
$ curl -v -u user:pass "http://127.0.0.1:9082/connect-auth?hostname=fs1&section=directory&tag_name=domain&key_name=name&key_value=connect.example.com&Event-Name=REQUEST_PARAMS&Core-UUID=0efa6954-cb1f-4b7c-b280-ba3434907a74&FreeSWITCH-Hostname=fs1&FreeSWITCH-Switchname=fs1&FreeSWITCH-IPv4=198.11.254.113&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2015-02-20%2014%3A08%3A45&Event-Date-GMT=Fri,%2020%20Feb%202015%2020%3A08%3A45%20GMT&Event-Date-Timestamp=1424462925448233&Event-Calling-File=sofia_reg.c&Event-Calling-Function=sofia_reg_parse_auth&Event-Calling-Line-Number=2741&Event-Sequence=518&action=sip_auth&sip_profile=internal&sip_user_agent=Blink%20Pro%204.1.0%20(MacOSX)&sip_auth_username=%2B141238510000&sip_auth_realm=connect.example.com&sip_auth_nonce=7cff7d80-afbb-4e0c-9455-c0d9d6baa3f7&sip_auth_uri=sip%3Aconnect.example.com&sip_contact_user=65341890&sip_contact_host=24.126.169.35&sip_to_user=%2B141238510000&sip_to_host=connect.example.com&sip_via_protocol=udp&sip_from_user=%2B141238510000&sip_from_host=connect.example.com&sip_call_id=XIiqYI.uxMNatGF6m76zb6HEFe0skr46&sip_request_host=connect.example.com&sip_auth_qop=auth&sip_auth_cnonce=mSfOYudMXnSRJsLwnrKiWJvhrChlrV00&sip_auth_nc=00000001&sip_auth_response=bcf3664ef7e5a84578967dde61d72c26&sip_auth_method=REGISTER&client_port=56510&key=id&user=%2B141238510000&domain=connect.example.com&ip=24.126.169.35"
* Hostname was NOT found in DNS cache
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 9082 (#0)
* Server auth using Basic with user 'tropo'
> POST /connect-auth?hostname=fs1&section=directory&tag_name=domain&key_name=name&key_value=connect.example.com&Event-Name=REQUEST_PARAMS&Core-UUID=0efa6954-cb1f-4b7c-b280-ba3434907a74&FreeSWITCH-Hostname=fs1&FreeSWITCH-Switchname=fs1&FreeSWITCH-IPv4=198.11.254.113&FreeSWITCH-IPv6=%3A%3A1&Event-Date-Local=2015-02-20%2014%3A08%3A45&Event-Date-GMT=Fri,%2020%20Feb%202015%2020%3A08%3A45%20GMT&Event-Date-Timestamp=1424462925448233&Event-Calling-File=sofia_reg.c&Event-Calling-Function=sofia_reg_parse_auth&Event-Calling-Line-Number=2741&Event-Sequence=518&action=sip_auth&sip_profile=internal&sip_user_agent=Blink%20Pro%204.1.0%20(MacOSX)&sip_auth_username=%2B141238510000&sip_auth_realm=connect.example.com&sip_auth_nonce=7cff7d80-afbb-4e0c-9455-c0d9d6baa3f7&sip_auth_uri=sip%3Aconnect.example.com&sip_contact_user=65341890&sip_contact_host=24.126.169.35&sip_to_user=%2B141238510000&sip_to_host=connect.example.com&sip_via_protocol=udp&sip_from_user=%2B141238510000&sip_from_host=connect.example.com&sip_call_id=XIiqYI.uxMNatGF6m76zb6HEFe0skr46&sip_request_host=connect.example.com&sip_auth_qop=auth&sip_auth_cnonce=mSfOYudMXnSRJsLwnrKiWJvhrChlrV00&sip_auth_nc=00000001&sip_auth_response=bcf3664ef7e5a84578967dde61d72c26&sip_auth_method=REGISTER&client_port=56510&key=id&user=%2B141238510000&domain=connect.example.com&ip=24.126.169.35 HTTP/1.1
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
    <domain name="connect.example.com">
      <params>
        <param name="dial-string" value="{presence_id=${dialed_user}@${dialed_domain}}${sofia_contact(${dialed_user}@${dialed_domain})}"></param>
      </params>
      <groups>
        <group name="default">
          <users>
            <user id="+141238510000" cacheable="300000">
              <params>
                <param name="password" value="abc123"></param>
              </params>
            </user>
          </users>
        </group>
      </groups>
    </domain>
  </section>
* Connection #0 to host 127.0.0.1 left intact
</document>%
$
```


#### Stats API

##### Go Stats

This API will return go virtual machine stats  [http://localhost:9082/stats/go](http://localhost:9082/stats/go)

```json
{
  "time": 1427915057163330600,
  "go_version": "go1.4",
  "go_os": "darwin",
  "go_arch": "amd64",
  "cpu_num": 8,
  "goroutine_num": 9,
  "gomaxprocs": 1,
  "cgo_call_num": 0,
  "memory_alloc": 163456,
  "memory_total_alloc": 224224,
  "memory_sys": 2885880,
  "memory_lookups": 9,
  "memory_mallocs": 1236,
  "memory_frees": 321,
  "memory_stack": 212992,
  "heap_alloc": 163456,
  "heap_sys": 835584,
  "heap_idle": 393216,
  "heap_inuse": 442368,
  "heap_released": 0,
  "heap_objects": 915,
  "gc_next": 215856,
  "gc_last": 1427915055894457600,
  "gc_num": 4,
  "gc_per_second": 0,
  "gc_pause_per_second": 0,
  "gc_pause": [
    0.131137,
    0.083353,
    0.075486,
    0.088565
  ]
}

```
##### HTTP Stats

This API will return http stats  [http://localhost:9082/stats/http](http://localhost:9082/stats/http)

```json
{
    "pid": 36719,
    "uptime": "5.041052691s",
    "uptime_sec": 5.041052691,
    "time": "2015-04-01 14:59:40.169469799 -0400 EDT",
    "unixtime": 1427914780,
    "status_code_count": {},
    "total_status_code_count": {
        "200": 1
    },
    "count": 0,
    "total_count": 1,
    "TotalResponseTime": "210.859µs",
    "total_response_time_sec": 0.000210859,
    "average_response_time": "210.859µs",
    "average_response_time_sec": 0.000210859
}
````
