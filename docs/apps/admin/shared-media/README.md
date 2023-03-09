## setup shared-media persistent volume

first we ll need to manually create a persistent volume
```
kubectl create namespace media
cd /tmp
nano media-persistent-volume-claim.yaml
```
now we insert the following text into this file:
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: shared-media
  namespace: media
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1000Gi
  storageClassName: freenas-nfs-csi
```
adjust the size/storageClassName and make sure they match the values from the deploy/argocd/bootstrap-optional-apps/values.yaml
```
kubectl apply -f media-persistent-volume-claim.yaml
persistentvolumeclaim/shared-media created
```

preparing the folder, this can be done from any pod that mounts the shared volume, i ll do it from the rtorrent pod.

```
❯ k exec -it -n media rtorrent-flood-84f6c7c6d4-mlhq2 ash
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
/ $ cd /downloads/
/downloads $ mkdir movies
/downloads $ mkdir tv
/downloads $ mkdir music
/downloads $ mkdir games
```