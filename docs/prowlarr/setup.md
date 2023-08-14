## setup prowlarr

login to the webui

#### Indexers

Head over to Indexers -> Add

- add indexers to your liking

#### Download Clients
Settings -> Download Clients

###### rtorrent-flood
- select option flood
- Host: rtorrent-flood
- UrlBase: /
- Username/Password ( as you set in ../rtorrent-flood/README.md )

###### nzbget 
- select option nzbget
- Host: nzbget.media
- Category: MOVIES
- username ( what you configured in nzbget )
- password ( what you configured in nzgbet )

- test
- save

#### Apps
Settings -> Apps

###### Radarr
- prowlar host: http://prowlarr.media:9696
- radar host:   http://radarr.media:7878
- api key - copy from radarr Settings -> General
- full sync

###### Sonarr
- prowlar host: http://prowlarr.media:9696
- radar host:   http://sonarr.media:8989
- api key - copy from sonarr Settings -> General
- full sync

###### Lidarr
- prowlar host: http://prowlarr.media:9696
- radar host:   http://lidarr.media:8989
- api key - copy from lidarr Settings -> General
- full sync