# Homelab
This repo contains definitions for a k3s-based kubernetes cluster that uses argocd and GitOps principles to deploy the applications running in the cluster. 

## three choices to install
![three choices to install](/docs/Excalidraw/three-choices.png)

## features

- :luc_codepen: uses kubernetes (k3s) and containers
- :luc_lock: supports full disk encryption ( debian/cryptsetup )
- :luc_octagon: uses a mix of go/terraform/bash to create the cluster
- :luc_settings_2: ships with a ./setup helper ( written in go )
- 🏠 modular architecture, easy to add or remove features/components ( PRs welcome )
- :luc_server: bootstraps hypervisor ( proxmox ), with templates ( debian ) and kvms ( k3s )
- :luc_github: follows **gitops** principles
- :luc_lock: automated certificate management with cert-manager
- :luc_map: automated ip management with external-dns
- :luc_hard_drive: supports external storage ( iscsi/nfs ) such as truenas provides via democratic-csi
- ☁️ ingress via nginx or cloudflare tunnels
- ✈️ contineous delivery with **argocd**
- :luc_book_open: log solution with grafana's loki
- :luc_lock: secrets encrypted with **sealed secrets**
- :luc_hard_drive: backups & disaster recovery with **kasten k10**
- :luc_arrow_down_circle: uses this repo as upstream to receive updates
- :luc_cloud: replaces the need for a cloud
- :luc_lock: two factor / SSO with authelia
- :luc_tornado: covers your media needs with jellyfin/jellyseer/prowlarr/radarr/sonarr/nzbget/rtorrent
- 👓 support for linkerd service mesh ( already integrated with grafana/loki ) - WIP
- :luc_book_open: documentation with obsidian ( ./docs ) or [online](https://loeken.github.io/homelab)

## installation & docs
docs are written in obsidian markdown inside the docs folder, an online version can be found online: https://loeken.github.io/homelab/ you can also access it via [github](https://github.com/loeken/homelab/blob/main/docs/index.md)

