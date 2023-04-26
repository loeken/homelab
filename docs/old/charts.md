## Helm Charts
2 helm charts to bootstrap argocd applications, each argocd app installs a helm chart
### deploy/helm/bootstrap-core-apps

| application  | description                                                                                                 |
| ------------ | ----------------------------------------------------------------------------------------------------------- |
| cert-manager                      | uses letsencrypt to create tls certificates                                            |
| letsencrypt cluster issuers       | staging/prod                                                                           |
| external-dns                      | updates dns with services like cloudflare )                                            |
| sealed secrets controller         | ( used to encrypt secrets in the cluster, allows storing encrypted secrets in github   |
| nginx ingress controller          | redirects traffic to the containers )                                                  |
| deploy/helm/import_sealed_secrets | a small helm chart that can be used to import sealed secrets from github )             |
      
### deploy/helm/bootstrap-optional-apps

| application  | description                                                                                                           | linked with openldap |
| ------------ | --------------------------------------------------------------------------------------------------------------------- | -------------------- |
| authelia                                             | one/two factor auth, works via ingress annotations ) - ldap support           | yes                  
| democratic iscsi fast                                | storage class to be used with fast storage such as ssds                       | 
| democratic iscsi                                     | storage class to be used with spindle disks                                   |
| democratic nfs                                       | nfs shared storage among multiple containers                                  |
| heimdall                                             | application dashboard - to redirect to all services running                   | 
| home-assistant                                       | self hosted smart home management                                             | 
| jellyfin                                             | dot.net based media player, no costs, allows file downloads ) - ldap support  | yes
| jellyseer                                            | media request portal for sonarr/radarr ) - ldap support                       | yes
| kasten-k10 | works with democratic-csi based storage, uses volumesnapshots for backups, supports offsite backups via S3/Backblaze b2 |
| loki                                                 | log solution with prometheus and a few dashboards pre-loaded                  | yes
| nextcloud                                            | a self hosted dropbox alternative )                                           | yes
| nzbget                                               | a usenet downloader                                                           |
| openldap            | centralized auth services for all services that have ldap support                                              |
| prowlarr            | a usenet/torrent indexer integrates with sonarr/radarr                                                         |
| radarr              | a movie collection manager ( integrates with rtorrent/nzbget/prowlarr )                                        |
| rtorrent-flood      | rtorrent based torrent client with a js based frontent                                                         |
| sonarr              | a tv c ollection manager ( integrates with rtorrent/nzbget/prowlarr )                                          |
| vaultwarden         | a rewrite of bitwarden in rust, configured with postgresql as database                                         | yes
| volume-snapshots    | minor things to allow volumesnapshots for democratic iscsi/nfs                                                 |
| whoami              | a small go based server that helps setting up ingress ( forwarded for headers )                                |