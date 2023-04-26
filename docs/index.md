## Three choices to install

![three choices to install](Excalidraw/three-choices.png)

## 1. Read this introduction
[index](index.md)

## 2. Gitops/Argocd Concept
[concepts of gitops and argocd](01_gitops_concepts.md)

## 3. Clone loeken/homelab github repo in your account
[setup private fork of github repository](02_github_repository.md)

## 4. Installation of debian 11 with full disk encryption
[installation of debian FDE](03_debian_install.md)

## 5. Installation Overview
[installation overview](04_installation_overview.md)

## 6. Installation of homelab with Proxmox with nginx
[installation with nginx](05_beelink_nginx.md)

## 7. Installation of homelab with Proxmox with cloudflare tunnels
[installation with cloudflare tunnel](06_beelink_cloudflare_tunnel.md)

# The stack from bottom to top

These are all the apps that are part of this repo and can be installed via the ./setup util. It can either be installed in baremetal debian, in proxmox ( setup converts debian to proxmox & creates templates/kvms ). This allows intalling debian with FullDiskEncryption :luc_lock: easily with the debian installer.

| Name                    | Type                    | Description                                      | Optional |
| ----------------------- | ----------------------- | ------------------------------------------------ | -------- |
| Proxmox                 | Operating System        | Allows running KVMs, ships with Webui            | yes      |
| Debian 11               | Template                | Template for Proxmox KVMs                        | yes      |
| KVM                     | Virtual Machine         | A virtual machine ( KVM ) in Proxmox             | yes      |
| k3s                     | Kubernetes Distribution | For self hosting Kubernetes                      | yes      |
| Minikube                | Kubernetes Distribution | Alternative to Baremetal/Proxmox                 | yes      |
| Argocd                  | Kubernetes Application  | Declarative Githubs CD                           | no       |
| Bootstrap Core apps     | Argocd App of Apps      | Used to load required Kubernetes Applications    | no       |
| Bootstrap Optional apps | Argocd App of Apps      | Used to load required Kubernetes Applications    | no       |
| Sealed Secrets          | Kubernetes Application  | Encryption of all application credentials        | no       |
| Cert Manager            | Kubernetes Application  | X.509 certificate controller ( letsencrypt )     | yes      |
| Import Sealed Secrets   | Kubernetes Application  | a helm chart to import stored sealed secrets     | yes      |
| Authelia                | Kubernetes Application  | SingleSignOn Multi-Factor portal                 | yes      |
| Cloudflare Tunnels      | Ingress                 | Alternative ingress using cloudflare tunnels     | yes      |
| Democractic CSI         | Kubernetes Application  | Storage class for iscsi/nfs works with truenas   | yes      |
| External DNS            | Kubernetes Application  | For updating DNS records such as cloudflare      | yes      |
| Home Assistant          | Kubernetes Application  | Open source home automation. Privacy focussed    | yes      |
| Jellyfin                | Kubernetes Application  | Software for streaming videos                    | yes      |
| Jellyseerr              | Kubernetes Application  | Portal for managing download requests            | yes      |
| Kasten K10              | Kubernetes Application  | Backup & Restore, DR and offsite - backblaze b2  | yes      |
| Loki                    | Kubernetes Application  | Grafana & Loki, loads dashboards for linkerd     | yes      |
| Nextcloud               | Kubernetes Application  | Self hosted Dropbox, apps for ios/android        | yes      |
| NFS Provisioner         | Kubernetes Application  | NFS provisioner for local-path for shared-media  | yes      |
| NGINX Ingress           | Ingress                 | Main recommended Ingress controller              | yes      |
| Nzbget                  | Kubernetes Application  | Usenet Downloader                                | yes      |
| Prowlarr                | Kubernetes Application  | Indexer for Usenet/Torrent Trackers              | yes      |
| Radarr                  | Kubernetes Application  | Movie Collection Manager for Usenet/Torrents     | yes      |
| Sonarr                  | Kubernetes Application  | TV Collection Manager for Usenet/Torrents        | yes      |
| Vaultwarden             | Kubernetes Application  | Rust based bitwarden, Password Manager           | yes      |
| Volume Snapshots        | Helm Chart              | Adds snapshot support for democratic-csi         | yes      |
| Whoami                  | Kubernetes Application  | a simple go app to display http requests/headers | yes      | 

