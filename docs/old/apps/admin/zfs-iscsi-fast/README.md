## democratic csi iscsi

#### preparation of truenas

add your ssds into a pool in truenas

#### customize your settings

```
cd deploy/mysecrets/
cp argocd-democratic-csi-iscsi-fast.yaml.example templates/argocd-democratic-csi-iscsi-fast.yaml
nano argocd-democratic-csi-iscsi-fast.yaml
```

now in here change the ip 172.16.137.13 with the ip of your truenas, update the root password for truenas and insert the id_ed25519_truenas private key, after that is done we can send this to kubeseal and apply the encrypted version

```
cat argocd-democratic-csi-iscsi-fast.yaml | kubeseal > templates/argocd-democratic-csi-iscsi-encrypted-fast.yaml 
git add templates/argocd-democratic-csi-iscsi-encrypted.yaml
git commit -m "added argocd-democratic-csi-iscsi-encrypted.yaml"
git push
```

#### cleanup

to cleanup the zvols in truenas that exist in truenas but do not exist in k8s you can use this:
https://github.com/democratic-csi/democratic-csi/blob/master/bin/k8s-csi-cleaner