# kasten k10
port-forward to k10 dashboard:

```
kubectl -n kasten-io port-forward svc/gateway 8082:8000
```

visit: http://localhost:8082/k10/#/dashboard