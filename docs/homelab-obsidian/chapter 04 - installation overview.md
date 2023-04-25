# Installation overview

when running ./setup install there are quite a few flags to chose from

```bash
./setup install -h
Flags:
--authelia string                               enable argocd app authelia (default "false")
--bridge string                                 name of the bridge (default "vmbr0")
--cluster-issuer string                         when using nginx ingress, select cluster issuer ( staging / prod ) (default "staging")
--cores_k3s string                              the amount of virtual cores to pass to the k3s vm (default "10")
--disaster_recovery string                      disaster recovery type (default "k10")
--disksize string                               disk size + metric ( example: 100G ) (default "100G")
--domain string                                 the domain you want to use
--email string                                  the folder into which new_repo will be cloned into
--external_ip string                            your external ipv4 ( curl -4 ifconfig.co ) (default "1.2.3.4")
--externaldns string                            enable argocd app external-dns (default "false")
--ha string                                     enable argocd app home-assistant (default "false")
-h, --help                                          help for install
--ingress string                                which ingress to use ( nginx/cloudflaretunnel ) (default "cloudflaretunnel")
--interface string                              name of the primary interface (default "enp3s0")
--jellyfin string                               enable argocd app jellyfin (default "false")
--jellyseerr string                             enable argocd app jellyseerr (default "false")
--kasten-k10 string                             enable argocd app kasten-k10 (default "false")
--kubernetes_version string                     kubernetes version (default "v1.25.6+k3s1")
--local_path string                             the folder into which new_repo will be cloned into
--loki string                                   enable argocd app loki (default "false")
--macaddr string                                mac address used for the k3s vm (default "6E:1F:26:B6:DF:20")
--memory_k3s string                             amount of ram in MB to assign to the VM  (default "28672")
--new_repo string                               the name of your new repo
--nextcloud string                              enable argocd app nextcloud (default "false")
--nginxingress string                           enable argocd app nginx-ingress (default "false")
--nzbget string                                 enable argocd app nzbget (default "false")
--partition_external_shared_media_disk string   will partition --shared-media-disk-device (default "false")
--pci_device string                             the pci address of your gpu ( lspci |grep VGA ) (default "0000:02:00.0")
--pci_passthrough string                        prepare pci passthrough (default "intel")
--platform string                               server type: proxmox, minikube, baremetal (default "proxmox")
--prowlarr string                               enable argocd app prowlarr (default "false")
--proxmox_node_name string                      the name of the proxmox node ( hostname ) (default "beelink-sei12")
--proxmox_vm_name string                        name of the virtual machine in proxmox (default "k3s-beelink-01")
--radarr string                                 enable argocd app radarr (default "false")
--root_password string                          root password ( used for login to proxmox ) (default "topsecure")
--rtorrentflood string                          enable argocd app rtorrent-flood (default "false")
--shared_media_disk_device string               give the device name of your external shared media disk (default "sda")
--shared_media_disk_size string                 define the size of the shared media disk (default "100Gi")
--smtp_domain string                            the domain from which the email is sent from (default "example.com")
--smtp_host string                              the host of your email server (default "mail.example.com")
--smtp_port string                              the port of your email server (default "587")
--smtp_sender string                            the email address used to send emails (default "homelab@example.com")
--smtp_username string                          the username used to login to your email (default "homelab@example.com")
--sonarr string                                 enable argocd app sonarr (default "false")
--ssh_password string                           ssh password (default "demotime")
--ssh_private_key string                        location of ssh private key, id_ed25519 when generated with gh auth login (default "~/.ssh/id_ed25519")
--ssh_public_key string                         location of ssh public key, id_ed25519.pub when generated with gh auth login (default "~/.ssh/id_ed25519.pub")
--ssh_server_address string                     ip address of server for ssh connection (default "172.16.137.36")
--ssh_server_gateway string                     gateway of server ( example 172.16.137.254 ) (default "172.16.137.254")
--ssh_server_netmask string                     amount of ram in MB to assign to the VM  (default "24")
--ssh_username string                           ssh usernamer (default "loeken")
--storage string                                storage type ( democratic-csi, local-path, ceph ) (default "local-path")
--vaultwarden string                            enable argocd app vaultwarden (default "false")
--whoami string                                 enable argocd app whoami (default "true")
```


Here are the most important ones

## --platform

This defines what installation you want to end up with. The first option is "proxmox". This expects a physical server with debian 11 pre-installed, the ./setup will then turn this debian into a proxmox server, create a kvm template, then a kvm based on this template and install k3s inside that vm. The second option is baremetal, where it skips turning debian into proxmox but installs k3s directly into debian. the proxmox option is recommended when you want to run other vms on the same server, if you want to purely enjoy k3s use the baremetal way. The third option is minikube, this will skip debian/proxmox/k3s setup, but install argocd directly into your minikube

## --storage
two main options here, the first is "local-path" this will use k3s built in "local-path" storageclass. This is most likely the way most ppl will install it. the second option is democratic-csi - this is recommended if you have an external storage such as a truenas. democratic-csi supports iscsi/nfs

## --shared_media*

the arr tools tools such as radarr/sonarr/prowlarr etc all work with the concept of having a shared storageclass (ReadWriteMany)) Out of the box k3s's local-path does not support ReadWriteMany. When using --platform proxmox or --platform baremetal the setup will install a nfs server in debian/proxmox and also setup a nfs-provisioner. this allows the arr tools to sue a shared media disk.

--shared-media-disk can be set to false or to a disk. When you only have 1 disk in your server set it to false. if you have 1 disk set it to that disk for example --shared-media-disk /dev/sdb. If you have more then 1 disk create a software array and set it to --shared-media-disk /dev/md0 ( or whatever your raid device is named ).

--shared_media_disk_size will specify how much of the data is allocated to the nfs server. When using --platform proxmox, it will partition the external disk and mount under /mnt/data. it will then add that to proxmox as a datastore and when creating the kvm it will add a second virtual disk and place it inside /mnt/data, you can set how big that virtual disk image should be here.

--partition_external_shared_media_disk can be used to make terraform partition your disk and create filesystems on them.

to reduce the risk of picking the wrong disk etc you need to specifically destroy the records on the external disks first you can do so via
```bash
wipefs --all /dev/yourdisk
```

## --ssh_*

--ssh_password string                           ssh password (default "demotime")
--ssh_private_key string                        location of ssh private key, id_ed25519 when generated with gh auth login (default "~/.ssh/id_ed25519")
--ssh_public_key string                         location of ssh public key, id_ed25519.pub when generated with gh auth login (default "~/.ssh/id_ed25519.pub")
--ssh_server_address string                     ip address of server for ssh connection (default "172.16.137.36")
--ssh_server_gateway string                     gateway of server ( example 172.16.137.254 ) (default "172.16.137.254")
--ssh_server_netmask string                     amount of ram in MB to assign to the VM  (default "24")
--ssh_username string                           ssh usernamer (default "loeken")
--root_password

the --ssh_* flags specify how to connect to your server. you specify the non root user created during install and the password, ./setup will then login over ssh and set the root password ( as we lll need that later on to login to the proxmox webui ).

If you already have an (older) ssh keypair in your system of the type id_ed25519 you need to point the private/public key to their location. if you followed the github video and used the "gh auth login" to generate a command then you should already have ~/.ssh/id_ed25519 and ~/.ssh/id_ed25519.pub and dont need to specify any flags

## --external_ip 

external ips is the external ip of the server, this is used when installing k3s and sets the external_node_ip. all ingresses will inherit of this ip. If you have a dynamic ip on your server that changes take a look into the helpers folder, there is a bash script you can use to regularly check the ip and update/restart k3s if the ip changes, then external-dns would kick in and update the cluster's dns records

## --other
therer are plenty of other flags, pretty much every app that can run in the cluster has a flag for example --vaultwarden which can be set to true or false. true will make ./setup configure/install those, setting those to false will not install those parts. you can still turns on other components after your install

## --ingress
this defines which ingress to use, my recommendation is to use nginx ( as we also have external-dns/cert-manager ). nginx is also required if you plan on using authelia.
Cloudflare tunnels is a simple alternative that create a tunnel to cloudflare and route traffic through those tunnels. this avoids having to open ports locally.
