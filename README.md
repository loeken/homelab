# Homelab
This repo contains definitions for a k3s-based kubernetes cluster that uses argocd and GitOps principles to deploy the applications running in the cluster. 

## three choices to install
![three choices to install](/docs/Excalidraw/three-choices.png)

## features

- uses kubernetes (k3s) and containers
- supports full disk encryption ( debian/cryptsetup )
- uses a mix of go/terraform/bash to create the cluster
- ships with a ./setup helper ( written in go )
- modular architecture, easy to add or remove features/components ( PRs welcome )
- bootstraps hypervisor ( proxmox ), with templates ( debian ) and kvms ( k3s )
- follows **gitops** principles
- automated certificate management with cert-manager
- automated ip management with external-dns
- supports external storage ( iscsi/nfs ) such as truenas provides via democratic-csi
- ingress via nginx or cloudflare tunnels
- contineous delivery with **argocd**
- log solution with grafana's loki
- secrets encrypted with **sealed secrets**
- backups & disaster recovery with **kasten k10**
- uses this repo as upstream to receive updates
- replaces the need for a cloud
- two factor / SSO with authelia
- covers your media needs with jellyfin/jellyseer/prowlarr/radarr/sonarr/nzbget/rtorrent
- support for linkerd service mesh ( already integrated with grafana/loki ) - WIP
- documentation with obsidian ( ./docs ) or [online](https://loeken.github.io/homelab)

## installation & docs
docs are written in obsidian markdown inside the docs folder, an online version can be found online: https://loeken.github.io/homelab/ you can also access it via [github](https://github.com/loeken/homelab/blob/main/docs/index.md)

