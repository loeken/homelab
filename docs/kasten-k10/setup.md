# kasten k10
port-forward to k10 dashboard:

```
kubectl -n kasten-io port-forward svc/gateway 8082:8000
```

visit: http://localhost:8082/k10/#/dashboard

enter email/comapny name

## setup location profile

settings -> location -> new Profile

- S3 compatible
  - backblaze:
    - name = homelab-beelink
    - key ID = kasten s3 access key
    - application key = kasten S3 secret
    - bucket overview: Endpoint = Endpoint
    - region = for me eu
    - Bucket Name = obvious ;)
    - I enabled Immutable backups and set to 10 days of protection period

## turn on disaster recovery
settings -> K10 Disaster Recovery -> Enable K10DR
- select profile we created ( homelab-beelink )
- make sure to save the k10 ID

## manually kick off DR policy execution
