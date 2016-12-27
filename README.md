# tcp-paste share terminal output.

[![Build Status](https://travis-ci.org/gregory-m/tcp-paste.svg?branch=master)](https://travis-ci.org/gregory-m/tcp-paste)

Inspired by [fiche](https://github.com/solusipse/fiche), but with built in HTTP server.

##Installation
Download compiled binary form [releases page](http://github.com/gregory-m/tcp-paste/releases).

Or if you want to build from source:

```
$ go get -u github.com/gregory-m/tcp-paste
```

##Usage:
```
Usage of tcp-paste:
  -hostname string
    	Hostname to use in links (default "localhost:8080")
  -http-host string
    	Host and port for HTTP connections (default ":8080")
  -storage string
    	Storage directory (default "/tmp")
  -tcp-host string
    	Host and port for for TCP connections (default ":4343")
```

``-hostname`` Hostname to uses in links for example if you deploy app to example.com you want to set in to ``example.com``

``-http-host`` HTTP port and host (used to view saved outputs) to listen on in flowing format: ``host:port`` if host part is omitted ``0.0.0.0`` will be used.

``-tcp-host`` TCP port and host (used to saved output) to listen on in flowing format: ``host:port`` if host part is omitted ``0.0.0.0`` will be used.

``-storage`` Storage directory usually you want to set it to something different then ``/tmp`` to preserve saved files after reboot.

Example:

```
$ tcp-paste -hostname example.com -http-host= :80 -tcp-host=443 -storage=/opt/tcp-paste
```
Note: In this example we listen on ports 443 and 80 on linux you can use ``etcap 'cap_net_bind_service=+ep' $(where tcp-paste)``
