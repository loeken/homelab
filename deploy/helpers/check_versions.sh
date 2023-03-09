#!/bin/bash

function lastversion() {
    version=$(curl -s https://artifacthub.io/api/v1/packages/helm/$1/$2/feed/rss | xmlstarlet sel -t -m "/rss/channel/item[1]" -v "title")
    echo "$version"
    mkdir -p versions/$1
    currentVersion=`cat versions/$1/$2`
    if [ "$currentVersion" == "" ]; then
        echo $version > versions/$1/$2
    fi
}

jellyfinversion=`lastversion k8s-at-home jellyfin`
jellyfinCurrentVersion=`cat versions/k8s-at-home/jellyfin`
echo jellyfin current: $jellyfinCurrentVersion remote: $jellyfinversion

# argocd
argocdVersion=`lastversion argo argo-cd`
argocdVersionCurrent=`cat versions/argo/argo-cd`
echo argo-cd current: $argocdVersionCurrent remote: $argocdVersion

# certManager
certManagerVersion=`lastversion cert-manager cert-manager`
certManagerVersionCurrent=`cat versions/cert-manager/cert-manager`
echo cert-manager current: $certManagerVersionCurrent remote: $certManagerVersion

# externalDns
externalDnsVersion=`lastversion bitnami external-dns`
externalDnsVersionCurrent=`cat versions/bitnami/external-dns`
echo external-dns current: $externalDnsVersionCurrent remote: $externalDnsVersion

# nginx-ingress
nginxIngresVersion=`lastversion bitnami nginx-ingress-controller`
nginxIngresVersionCurrent=`cat versions/bitnami/nginx-ingress-controller`
echo nginx-ingress-controller current: $nginxIngresVersionCurrent remote: $nginxIngresVersion

# sealed-secrets
sealedSecretsVersion=`lastversion bitnami-labs sealed-secrets`
sealedSecretsVersionCurrent=`cat versions/bitnami-labs/sealed-secrets`
echo sealed-secrets current: $sealedSecretsVersionCurrent remote: $sealedSecretsVersion

# openldap-stack-ha
openldapVersion=`lastversion helm-openldap openldap-stack-ha`
openldapVersionCurrent=`cat versions/helm-openldap/openldap-stack-ha`
echo openldap-stack-ha current: $openldapVersionCurrent remote: $openldapVersion

# authelia
echo authelia check on https://charts.authelia.com/

# vaultwarden
vaultwardenVersion=`lastversion k8s-at-home vaultwarden`
vaultwardenVersionCurrent=`cat versions/k8s-at-home/vaultwarden`
echo vaultwarden current: $vaultwardenVersionCurrent remote: $vaultwardenVersion

# heimdall
heimdallVersion=`lastversion k8s-at-home heimdall`
heimdallVersionCurrent=`cat versions/k8s-at-home/heimdall`
echo heimdall current: $heimdallVersionCurrent remote: $heimdallVersion

# rtorrent-flood
rtorrentFloodVersion=`lastversion k8s-at-home rtorrent-flood`
rtorrentFloodVersionCurrent=`cat versions/k8s-at-home/rtorrent-flood`
echo rtorrent-flood current: $rtorrentFloodVersionCurrent remote: $rtorrentFloodVersion

# overseerr
rtorrentFloodVersion=`lastversion k8s-at-home overseerr`
rtorrentFloodVersionCurrent=`cat versions/k8s-at-home/overseerr`
echo overseerr current: $rtorrentFloodVersionCurrent remote: $rtorrentFloodVersion

# sonarr
version=`lastversion k8s-at-home sonarr`
versionCurrent=`cat versions/k8s-at-home/sonarr`
echo sonarr current: $versionCurrent=`cat versions/k8s-at-home/sonarr` remote: $version

# radarr
version=`lastversion k8s-at-home radarr`
versionCurrent=`cat versions/k8s-at-home/radarr`
echo radarr current: $versionCurrent=`cat versions/k8s-at-home/radarr` remote: $version

# prowlarr
version=`lastversion k8s-at-home prowlarr`
versionCurrent=`cat versions/k8s-at-home/prowlarr`
echo prowlarr current: $versionCurrent=`cat versions/k8s-at-home/prowlarr` remote: $version

# nzbget
version=`lastversion k8s-at-home nzbget`
versionCurrent=`cat versions/k8s-at-home/nzbget`
echo nzbget current: $versionCurrent=`cat versions/k8s-at-home/nzbget` remote: $version

# nextcloud
version=`lastversion nextcloud nextcloud`
versionCurrent=`cat versions/nextcloud/nextcloud`
echo nextcloud current: $versionCurrent=`cat versions/nextcloud/nextcloud` remote: $version
