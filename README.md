# tcp-paste share terminal output.

Inspired by [fiche](https://github.com/solusipse/fiche), but with built in HTTP server and slack support.

![Screenshot](https://i.imgur.com/gVMPnTY.gif "Logo Title Text 1")


## Client
### Usage:
Just pipe output to netcat:

```
$ ps ax | nc tcp-paste.server.com 4343
```

Multiple commands also supported:

```
$ (ps ax && ls -la) | nc tcp-paste.server.com 4343
```

You can redirect stderr as usual:

```
ssh -vv some.problemactic.host.com ls 2>&1 | nc tcp-paste.server.com 4343
```

## Server
### Usage:
```
Usage of tcp-paste:
 -hostname string
    	Hostname to use in links (default "localhost:8080")
  -http
    	Run HTTP service (to server saved output from local disk) (default true)
  -http-host string
    	Host and port for HTTP service (default ":8080")
  -paste
    	Run paste service (save output on local disk) (default true)
  -paste-host string
    	Host and port for post service (default ":4343")
  -slack
    	Run Slack service (to post output to slack) (default true)
  -slack-chanel string
    	Slack API token (default "testa")
  -slack-host string
    	Host and port for slack service (default ":9393")
  -slack-token string
    	Slack API token
  -storage string
    	Storage directory (default "/tmp")
```

---
``-hostname`` Hostname to uses in links for example if you deploy app to example.com you want to set in to ``example.com``

``-storage`` Storage directory usually you want to set it to something different then ``/tmp`` to preserve saved files after reboot.

---

``-http`` Run http service to serve saved files

``-http-host`` HTTP service port and host in flowing format: ``host:port`` if host part is omitted ``0.0.0.0`` will be used.

---
``-paste`` Run paste service, listen to connections and save input to disk

``-paste-host`` Paste service port and host in flowing format: ``host:port`` if host part is omitted ``0.0.0.0`` will be used.

---
``-slack`` Run slack service to listen to connections and post input to slack.

``-slack-host`` Slack service port and host in flowing format: ``host:port`` if host part is omitted ``0.0.0.0`` will be used.

``-slack-token`` Slack API token.

``-slack-chanel`` Slack channel to post




Example:

```
$ -hostname example.com -storage=/opt/tcp-paste -http-host:80 -paste-host=443 -slack-host=993 -slack-token="very-secret" -slack-chanel=general
```
Note: In this example we listen on ports 443 and 80 on linux you can use ``etcap 'cap_net_bind_service=+ep' $(where tcp-paste)``

### Installation
Download compiled binary form [releases page](http://github.com/gregory-m/tcp-paste/releases).

Or if you want to build from source:

```
$ go get -u github.com/gregory-m/tcp-paste
```

Or you can use docker images.

#### Docker images:
You can use [docker image](https://hub.docker.com/r/gregorym/tcp-paste/).



For example to listen on ports 443 and 80, and to use host /opt/tcp-paste directory as data storage and example.com as hostname:

```
$ docker run -p 80:8080 -p 443:4343 -p 993:9393 -e HOSTNAME=gregory.beer \
-e SLACK_TOKEN=very-secret -e SLACK_CHANNEL=testa \
-v /opt/tcp-paste:/data gregorym/tcp-paste
```
