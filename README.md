# Homelab
This repo contains definitions for a k3s-based kubernetes cluster that uses argocd and GitOps principles to deploy the applications running in the cluster. 

# note
the ./setup is still a work in progress, not all options have been finished/tested

## setup
This repository contains a command line app that allows you to easily set up a homelab environment. The ./setup script provides a number of commands that can be used to perform various tasks, including installing the homelab stack, checking system dependencies, and configuring your GitHub repository.


```bash
homelab: ❯ ./setup -h
My command line app

Usage:
  ./setup [command]

Available Commands:
  check-dependencies Check if all needed dependencies are installed on this system
  completion         Generate the autocompletion script for the specified shell
  destroy            destroy the stack ( DANGER DANGER! :) )
  enable-argocd-app  little helper to load sealed secrets
  github             Create/clone/configure upstream of loeken/homelab in my github account by using the gh command line client
  help               Help about any command
  install            install the stack

Flags:
  -h, --help   help for ./setup

Use "./setup [command] --help" for more information about a command.
```

## 1. check dependencies
this setup.go tool uses a few tools ( gh, cloudflared, git, terraform, code ) to ensure you have all dependencies run this command:
```bash
❯ ./setup check-dependencies --new_repo loeken/homelab-beelink
gh version: 2.23.0
cloudflared version: 2023.2.1
git version: 2.39.1
terraform version: 1.3.8
code version: 1.75.1
The current Kubernetes context is not homelab-beelink
```
you can ignore the last line for now as we havent setup a kubernetes context for this cluster yet.


## 2. github
this function will create a new repo in your github account, mirror the contents of this repo and configure this repo as an upstream ( so you can pull changes/updates in the future from upstream ).
```bash
./setup github --new_repo loeken/homelab-beelink --local_path /home/loeken/Projects/private
```

this would create github.com/loeken/homelab-beelink and clone it into /home/loeken/Projects/private and then open it up in vscode for further usage.


## 3. minimal install
```bash
./setup install --email loeken@internetz.me --new_repo loeken/homelab-beelink --domain loeken.xyz --ingress cloudflaretunnel --external_ip 1.2.3.4
```
this expects a debian 11 ( full disk encryption! ) you can define connection via the --ssh* flags. it will take this debian turn it into proxmox, create a kvm, install k3s, bootstrap argocd. as a default ingress it will setup a cloudflare tunnel using the cloudflared api, configure the domain loeken.xyz and also enable the service whoami.loeken.xyz.

At this stage you ll have a working k3s cluster with argocd, an ingress ( without having to open ports, setup load balancers ).

## 4. all installation flags:
```bash
❯ ./setup install -h
install the stack

Usage:
  ./setup install [flags]

Flags:
      --authelia string             enable argocd app authelia (default "false")
      --bridge string               name of the bridge (default "vmbr0")
      --cloudflaretunnel string     enable argocd app cloudflaretunnel (default "true")
      --cluster-issuer string       when using nginx ingress, select cluster issuer ( staging / prod ) (default "staging")
      --cores_k3s string            the amount of virtual cores to pass to the k3s vm (default "8")
      --disaster_recovery string    disaster recovery type (default "k10")
      --disksize string             disk size + metric ( example: 400G ) (default "400G")
      --domain string               the domain you want to use
      --email string                the folder into which new_repo will be cloned into
      --external-dns string         enable argocd app external-dns (default "false")
      --external_ip string          your external ipv4 ( curl -4 ifconfig.co ) (default "1.2.3.4")
      --grafana string              enable argocd app grafana (default "false")
      --ha string                   enable argocd app home-assistant (default "false")
      --heimdall string             enable argocd app heimdall (default "false")
  -h, --help                        help for install
      --iface string                name of the primary interface (default "enp3s0")
      --ingress string              which ingress to use ( nginx/cloudflaretunnel ) (default "cloudflaretunnel")
      --jellyfin string             enable argocd app jellyfin (default "false")
      --jellyseerr string           enable argocd app jellyseerr (default "false")
      --kasten-k10 string           enable argocd app kasten-k10 (default "false")
      --kubernetes_version string   kubernetes version (default "v1.25.6+k3s1")
      --local_path string           the folder into which new_repo will be cloned into
      --macaddr string              mac address used for the k3s vm (default "6E:1F:26:B6:DF:20")
      --memory_k3s string           amount of ram in MB to assign to the VM  (default "28672")
      --new_repo string             the name of your new repo
      --nextcloud string            enable argocd app nextcloud (default "false")
      --nginxingress string         enable argocd app nginx-ingress (default "false")
      --nzbget string               enable argocd app nzbget (default "false")
      --pci_passthrough string      prepare pci passthrough (default "intel")
      --platform string             storage type (default "proxmox")
      --prowlarr string             enable argocd app prowlarr (default "false")
      --proxmox_node_name string    the name of the proxmox node ( hostname ) (default "beelink-sei12")
      --proxmox_vm_name string      name of the virtual machine in proxmox (default "k3s-beelink-01")
      --radarr string               enable argocd app radarr (default "false")
      --root_password string        root password ( used for login to proxmox ) (default "topsecure")
      --rtorrentflood string        enable argocd app rtorrent-flood (default "false")
      --sonarr string               enable argocd app sonarr (default "false")
      --ssh_password string         ssh password (default "demotime")
      --ssh_private_key string      location of ssh private key (default "~/.ssh/id_rsa")
      --ssh_public_key string       location of ssh public key (default "~/.ssh/id_rsa.pub")
      --ssh_server_address string   ip address of server for ssh connection (default "172.16.137.36")
      --ssh_server_gateway string   gateway of server ( example 172.16.137.254 ) (default "172.16.137.254")
      --ssh_server_netmask string   amount of ram in MB to assign to the VM  (default "24")
      --ssh_username string         ssh usernamer (default "loeken")
      --storage string              storage type ( democratic-csi, local, ceph ) (default "local")
      --vaultwarden string          enable argocd app vaultwarden (default "false")
      --volumesnapshots string      enable argocd app volume-snapshots (default "false")
      --whoami string               enable argocd app whoami (default "true")

```