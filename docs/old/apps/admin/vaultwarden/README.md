## setup vaultwarden

#### rollout default application

by default i turn on the ADMIN_TOKEN, so we can do the inital config of the .env vars via a nice webui.

#### create secrets
```
kubectl create namespace vaultwarden
cd deploy/mysecrets
cp argocd-vaultwarden.yaml.example argocd-vaultwarden.yaml
cp argocd-vaultwarden-postgresql.yaml.example argocd-vaultwarden-postgresql.yaml
nano argocd-vaultwarden.yaml
nano argocd-vaultwarden-postgresql.yaml
cat argocd-vaultwarden.yaml | kubeseal > argocd-vaultwarden-encrypted.yaml
kubectl apply -f argocd-vaultwarden-encrypted.yaml
cat argocd-vaultwarden-postgresql.yaml | kubeseal > argocd-vaultwarden-postgresql-encrypted.yaml 
kubectl apply -f argocd-vaultwarden-postgresql-encrypted.yaml 
```

now rollout the argocd application
