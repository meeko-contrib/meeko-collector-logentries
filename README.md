# meeko-collector-logentries #

[![Build
Status](https://drone.io/github.com/meeko-contrib/meeko-collector-logentries/status.png)](https://drone.io/github.com/meeko-contrib/meeko-collector-logentries/latest)

Meeko collector for Logentries webhooks

## Meeko Variables ##

* `LISTEN_ADDRESS` - TCP network address to listen on; format `HOST:PORT`
* `ACCESS_TOKEN` - Token to be used for for webhook authentication. The token is
  expected to be set via a query parameter `token`, e.g.
  `https://example.com?token=secret`.

## License ##

GNU GPLv3, see the `LICENSE` file.
