## setup authelia

#### notifier
```
cd deploy/mysecrets/
cp argocd-authelia.yaml.example argocd-authelia.yaml
nano argocd-authelia.yaml
```

```
cat argocd-authelia.yaml | kubeseal > argocd-authelia-encrypted.yaml
```