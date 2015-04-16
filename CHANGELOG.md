# 0.5.0

* Support setting `allow_direct_sip_out` from PAPI resouce configration. Used to allow a user to place direct sip calls from FS without going through tropo-gateway - OPS-3615

# 0.4.0

* Update to include http stats using [stats](https://github.com/thoas/stats)

# 0.3.0

* Add stats api for checking process health restfully

# 0.2.4

* remove user_context debug var

# 0.2.3

* Support for enforcing domain when users register. Default to false

# 0.2.2

* User profile variables should be decendent of User

# 0.2.1

* Rename cache config options
* Properly use cache config
* Ensure all configs are in seconds

# 0.2.0

* Refactor PAPI request methods
* Support returning TollPlan for Directory responses

# 0.1.0

* Support caching

# 0.0.6

* Create users auth resource OPS-3505

# 0.0.5

* include version numbers in log output for health check endpoint

# 0.0.4

* Improved logging

# 0.0.3

* Require e164 address for lookup

# 0.0.2

* Add connectivty check to /health route

# 0.0.1

* Intial release

