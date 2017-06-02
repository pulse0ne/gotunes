# gotunes
A mpd-based music player/client/server written in Go and Angular

## Server
#### Starting the Server
```bash
./gotunes
```
Running the server with a different config
```bash
./gotunes -c some-config.json
```

#### Config Options
Config:
```json
{
    "mpdhost": "localhost",
    "mpdport": 6600,
    "httpport": 80,
    "webroot": "./public",
    "startmpd": true,
    "loglevel": "info"
}
```
* `mpdhost` and `mpdport` are pretty self-explanatory
* `httpport` is the port that the client will be served on (`80` will require root)
* `webroot` should point to the location on disk that contains the web files
* `startmpd` tells the server to attempt to start `mpd` if it isn't already started
* `loglevel` should be one of the following: "debug", "info", "warn", "error"

## Client
#### Client Settings
On the settings page, the following options are available:
* Theme - there are a number of color themes available
* Disabled Views - gives the ability to disable certain views (Artist, Album, Track)
* Idle View Enabled - enables/disables the idle view
* Idle Delay - number of seconds of inactivity until the idle view is triggered
* Dynamic Title Enabled - enables/disables the dynamic title (putting title/artist in browser title)
