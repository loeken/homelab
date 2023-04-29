# 06 - installation cloudflare tunnel
I have a beelink sei12 with 32GB of ram, this comes with an interal ssd ( 500GB ) and ive added a 4TB 2.5" evo ssd in there. This one i ll be setting up with proxmox as it will run a few small other vms besides the kvm that runs homelab. Of that external 4TB ill use 3TB for homelab

![Beelink](Excalidraw/beelink.png)

since i want to partition the /dev/sda i ll ssh into the server and prepare it
```bash
sudo wipefs --all /dev/sda
```

now I can get ready to start the install 

```bash
./setup install --platform proxmox \
				--email loeken@internetz.me \
				--external_ip 46.142.60.140 \
				--ingress nginx \
				--domain loeken.xyz \
				--cores_k3s 10 \
				--memory_k3s 28672 \
				--cert-manager true \
				--cluster-issuer letsencrypt-staging \
				--disksize 100GB \
				--helo_name mail.internetz.me \
				--interface enp3s0 \
				--kubernetes_version v1.26.4+k3s1 \
				--macaddr 6E:1F:26:B6:DF:20 \
				--proxmox_node_name beelink-sei12 \
				--proxmox_vm_name k3s-beelink-01 \
				--authelia false \
				--externaldns true \
				--ha false \
				--jellyfin false \
				--jellyseerr false \
				--kasten-k10 false \
				--loki false \
				--new_repo loeken/homelab-beelink \
				--nextcloud false \
				--nzbget false \
				--prowlarr false \
				--radarr false \
				--sonarr false \
				--rtorrentflood false \
				--vaultwarden false \
				--ssh_password demotime \
				--ssh_private_key ~/.ssh/id_ed25519 \
				--ssh_public_key ~/.ssh/id_ed25519.pub \
				--ssh_server_address 172.16.137.36 \
				--ssh_server_gateway 172.16.137.254 \
				--ssh_server_netmask 24 \
				--ssh_username loeken \
				--shared_media_disk_size 3000G \
				--shared_media_disk_device sda \
				--partition_external_shared_media_disk true \
				--root_password topsecure \
				--storage local-path \
				--cloudflare_api_token REPLACE

```

