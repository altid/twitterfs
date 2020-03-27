# Twitter

twitterfs is a file service, which presents a simple view of your Twitter feed.

`go get github.com/altid/twitterfs`

## Usage

Ensure you have done the following (for now, this will change)
To acquire these keys, log in to https://developer.twitter.com/en/apps and create a new "App" (using oauth1).
The consumer key and consumer secret will be provided upon completion.

```sh
export TWITTER_CONSUMER_KEY=mytwitterconsumerkey
export TWITTER_CONSUMER_SECRET=mytwitterconsumersecret

twitterfs [-p <dir>] [-s <srv>] | -t | -conf

```

- `<dir>` fileserver path. Will default to /tmp/altid if none is given
- `<srv>` service name to use. (Default `twitter`)

## Configuration

### Automatic

To create interactively, run `twitterfs -conf` and follow the steps

[![asciicast](https://asciinema.org/a/hi5pQKjk9NwVrAGU2kzm4hLZF.svg)](https://asciinema.org/a/hi5pQKjk9NwVrAGU2kzm4hLZF)


### Manual

```ndb

# altid/config - Place this in your operating systems' default configuration directory

service=twitter
	secret=myusersecret
	token=myusertoken
	log=/usr/halfwit/log
	listen_address=192.168.0.4

```

- service matches the given servicename (default "twitter")
- log is a location to store Twitter logs. A special value of `none` disables logging.
- listen_address is a more advanced topic, explained here: [Using listen_address](https://altid.github.io/using-listen-address.html)

> See [altid configuration](https://altid.github.io/altid-configurations.html) for more information
