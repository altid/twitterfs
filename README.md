# Twitter

twitterfs is a file service, which presents a simple view of your Twitter feed.

`go get github.com/altid/twitterfs`

## Usage


`twitterfs [-p <dir>] [-s <srv>]`

 - `<dir>` fileserver path. Will default to /tmp/altid if none is given
 - `<srv>` service name to use. (Default `twitter`)

## Configuration

```
# altid/config - Place this in your operating systems' default configuration directory

service=twitter address=twitter.com auth=pass=hunter2
	user=myloginemail
	log=/usr/halfwit/log
	#listen_address=192.168.0.4
```
 - service matches the given servicename (default "twitter")

 - address is currently ignored
 - auth is the authentication method
   - pass will send the string following pass= as your user password to the Twitter server
   - factotum uses a local factotm (Plan9, plan9port) to find your password
 - user is your login email for Twitter
 - log is a location to store Twitter logs. A special value of `none` disables logging.
 - listen_address is a more advanced topic, explained here: [Using listen_address](https://altid.github.io/using-listen-address.html)

> See [altid configuration](https://altid.github.io/altid-configurations.html) for more information
