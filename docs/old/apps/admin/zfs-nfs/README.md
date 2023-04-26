## democratic csi iscsi

#### basic truenas setup
make sure you followed the zfs-iscsi setup instructions for truenas
[zfs-iscsi](../zfs-iscsi/README.md)


#### customize your settings

```
cd deploy/mysecrets/
cp argocd-democratic-csi-nfs.yaml.example argocd-democratic-csi-nfs.yaml
nano argocd-democratic-csi-nfs.yaml
```

now in here change the ip 172.16.137.13 with the ip of your truenas, update the root password for truenas and insert the id_ed25519_truenas private key, after that is done we can send this to kubeseal and apply the encrypted version

```
cat argocd-democratic-csi-nfs.yaml | kubeseal > argocd-democratic-csi-nfs-encrypted.yaml 
kubectl apply -f argocd-democratic-csi-nfs-encrypted.yaml 
```
#### cleanup

to cleanup the zvols in truenas that exist in truenas but do not exist in k8s you can use this:
https://github.com/democratic-csi/democratic-csi/blob/master/bin/k8s-csi-cleaner