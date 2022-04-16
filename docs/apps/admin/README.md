# apps, recommended setup order

## 1.) Whoami
[whoami](whoami/README.md)

## 2.) volume-snapshots
[volume-snapshots](volume-snapshots/README.md)

## 3.) zfs-iscsi
[zfs-iscsi](zfs-iscsi/README.md)

## 4.) zfs-nfs
[zfs-nfs](zfs-nfs/README.md)

## 5.) Openldap
[openldap](openldap/README.md)

## 6.) authelia
[openldap](authelia/README.md)

## 7.) Vaultwarden
[vaultwarden](vaultwarden/README.md)

## 8.) Heimdall
[kuma](heimdall/README.md)

## 9.) shared-media disk
[shared-media](shared-media/README.md)

## 10.) Rtorrent-flood
[rtorrent-flood](rtorrent-flood/README.md)

## 11.) Plex
[plex](plex/README.md)

## 12.) nzbget
[nzbget](nzbget/README.md)

## 13.) Sonarr
[sonarr](sonarr/README.md)

## 14.) Radarr
[radarr](radarr/README.md)

## 15.) Lidarr
[lidarr](lidarr/README.md)

## 16.) Prowlarr
[prowlarr](prowlarr/README.md)

## 17.) Overseer
[overseerr](overseerr/README.md)

## 18.) Nextcloud
[nextcloud](nextcloud/README.md)

## 19.) harbor
[harbor](harbor/README.md)

## 20.) home assistant
[home-assistant](home-assistant/README.md)

## 21.) kasten k10
[k10](k10/README.md)


## General Update Process

#### general concept:

no forks needed, no private repos needed, just clone a copy of this repo, during the the setup process you basically create two files
- deploy/mysecrets/argocd-optional.yaml
- deploy/terraform/variables.tf

in the repo we only keep a .example copy of it. this allows you the end user to just pull the repo and not have your configs be overwritten. these files are also placed in the .gitignore file to not be accidently pushed to some repo. 

#### update this repo
```
git pull
```

#### updating argocd ( via terraform )
```
cd deploy/terraform
terraform apply
```

#### use the "app of apps" to send updates of all argocd apps
in the argocd ui head to
- bootstrap-core-apps
- bootstrap-optional-apps

click on Sync, select "all out of sync" and apply ( repeat for both folders )

still inside argocd's webui, go to the apps that received an update, then click on refresh ( but make sure you click that little arrow pointing down righe next to the text "REFRESH", a small option will appear to Hard Refresh, run this ). the hard refresh is required to also pull any updates made to the helm-values files ( inside the helm-values secret in the argocd namespace ).


#### pro tip:
you can also use the argocd cli to trigger these upgrades

```
argocd login localhost:8081

argocd app get bootstrap-core-apps --hard-refresh
argocd app get bootstrap-optional-apps --hard-refresh

argocd app sync bootstrap-core-apps
argocd app sync bootstrap-optional-apps

argocd app get authelia --hard-refresh
argocd app get cert-manager --hard-refresh
argocd app get cluster-issuer --hard-refresh
argocd app get external-dns --hard-refresh
argocd app get harbor --hard-refresh
argocd app get heimdall --hard-refresh
argocd app get home-assistant --hard-refresh
argocd app get k10 --hard-refresh
argocd app get kuma --hard-refresh
argocd app get lidarr --hard-refresh
argocd app get nextcloud --hard-refresh
argocd app get nginx-ingress --hard-refresh
argocd app get nzbget --hard-refresh
argocd app get openldap --hard-refresh
argocd app get overseerr --hard-refresh
argocd app get plex --hard-refresh
argocd app get prowlarr --hard-refresh
argocd app get radarr --hard-refresh
argocd app get rtorrent-flood --hard-refresh
argocd app get sealed-secrets-controller --hard-refresh
argocd app get sonarr --hard-refresh
argocd app get vaultwarden --hard-refresh
argocd app get volume-snapshots --hard-refresh
argocd app get whoami --hard-refresh
argocd app get zfs-iscsi --hard-refresh
argocd app get zfs-nfs --hard-refresh
```