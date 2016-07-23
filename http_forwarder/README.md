# HTTP Forwarder

Forwards HTTP requests to a list of backend hosts (tested against Squid proxy) depending on the Geo country of the request resource.

  - Properly tunnels tcp traffic after receiving CONNECT requests
  - Assumes backends are http proxies at the moment
