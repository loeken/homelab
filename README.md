# Homelab
This repo contains definitions for a k3s-based kubernetes cluster that uses argocd and GitOps principles to deploy the applications running in the cluster. 

## three choices to install
![three choices to install](/docs/Excalidraw/three-choices.png)

## features

:heavy_check_mark: uses kubernetes (k3s) and containers
:heavy_check_mark: supports full disk encryption ( debian/cryptsetup )
:heavy_check_mark: uses a mix of go/terraform/bash to create the cluster
:heavy_check_mark: ships with a ./setup helper ( written in go )
:heavy_check_mark: modular architecture, easy to add or remove features/components ( PRs welcome )
:heavy_check_mark: bootstraps hypervisor ( proxmox ), with templates ( debian ) and kvms ( k3s )
:heavy_check_mark: follows **gitops** principles
:heavy_check_mark: automated certificate management with cert-manager
:heavy_check_mark: automated ip management with external-dns
:heavy_check_mark: supports external storage ( iscsi/nfs ) such as truenas provides via democratic-csi
:heavy_check_mark: ingress via nginx or cloudflare tunnels
:heavy_check_mark: contineous delivery with **argocd**
:heavy_check_mark: log solution with grafana's loki
:heavy_check_mark: secrets encrypted with **sealed secrets**
:heavy_check_mark: backups & disaster recovery with **kasten k10**
:heavy_check_mark: uses this repo as upstream to receive updates
:heavy_check_mark: replaces the need for a cloud
:heavy_check_mark: two factor / SSO with authelia
:heavy_check_mark: covers your media needs with jellyfin/jellyseer/prowlarr/radarr/sonarr/nzbget/rtorrent
:heavy_check_mark: support for linkerd service mesh ( already integrated with grafana/loki ) - WIP
:heavy_check_mark: documentation with obsidian ( ./docs ) or [online](https://loeken.github.io/homelab)

## installation & docs
docs are written in obsidian markdown inside the docs folder, an online version can be found online: https://loeken.github.io/homelab/ you can also access it via [github](https://github.com/loeken/homelab/blob/main/docs/index.md)

