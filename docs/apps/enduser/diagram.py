from diagrams import Diagram, Cluster, Edge
from diagrams.custom import Custom
from diagrams.onprem.proxmox import Pve
from diagrams.k8s.network import Ingress, Service
from diagrams.k8s.clusterconfig import HPA
from diagrams.k8s.compute import Deployment, Pod, ReplicaSet
from diagrams.onprem.network import Opnsense, Internet
from diagrams.generic.network import Router,Switch
from diagrams.generic.storage import Storage



with Diagram("home", show=False):
    inet = Internet("Internet")
    router = Router("router")
    firewall = Opnsense("Opensense Firewall")
    switch = Switch("Switch")
   

    inet >> router >> firewall >> switch 
    
    with Cluster("Kubernetes Cluster"):
        ingr = Ingress("nginx ingress *.example.com") 
        svc1 = Service("nginx ingress 01")
        svc2 = Service("nginx ingress 02")
        svc3 = Service("nginx ingress 03")
        pve1 = Pve("Proxmox KVM 01")
        pve2 = Pve("Proxmox KVM 02")
        pve3 = Pve("Proxmox KVM 03")
        switch >> ingr
    
        ingr >> svc1 >> pve1
        ingr >> svc2 >> pve2
        ingr >> svc3 >> pve3

    with Cluster("Applications running in cluster"):
        pod1 = Pod("democratic csi for nfs/iscsi")
        pod2 = Pod("openldap")
        pod3 = Pod("authelia")
        pod4 = Pod("vaultwarden")
        pod5 = Pod("rtorrent")
        pod6 = Pod("nzbget")
        pod7 = Pod("plex")
        pod8 = Pod("sonarr")
        pod9 = Pod("radarr")
        pod10 = Pod("prowlarr")
        pod11 = Pod("overseer")
        pod12 = Pod("nextcloud")
        pod13 = Pod("home assistant")
        
        pve2 >> pod1
        pve2 >> pod2
        pve2 >> pod3
        pve2 >> pod4
        pve2 >> pod5
        pve2 >> pod6
        pve2 >> pod7
        pve2 >> pod8
        pve2 >> pod9
        pve2 >> pod10
        pve2 >> pod11
        pve2 >> pod12
        pve2 >> pod13

    with Cluster("redundant storage instance"):
        storage = Storage("truenas")
        pod1 >> storage
        pod2 >> storage
        pod3 >> storage
        pod4 >> storage
        pod5 >> storage
        pod6 >> storage
        pod7 >> storage
        pod8 >> storage
        pod9 >> storage
        pod10 >> storage
        pod11 >> storage
        pod12 >> storage        
        pod13 >> storage
    with Cluster("weekly offsite backups"):
        k10 = Custom("k10 backups", "./img/k10.png")
        storage >> k10