## Authelia

The default configuration of authelia is to use a simple yaml file for configuration such as:

```
users:
  loeken:
    disabled: false
    displayname: luke filewalker
    email: loeken@example.com
    groups:
    - users
    - admins
    password: $argon2id$v=19$m=65536,t=3,p=4$hdf77UqbU+pUaDtc9K4oqg$sr9ZG+DeEHYNK/qcBQKm36nMd0y2/ML4/Mszbud3ymE
```

you can generate a password via docker:
```
❯ docker run --rm authelia/authelia:latest authelia hash-password topsecret
Unable to find image 'authelia/authelia:latest' locally
latest: Pulling from authelia/authelia
c158987b0551: Pull complete 
ad15638eea56: Pull complete 
79a91ce2b249: Pull complete 
2c361976b6a9: Pull complete 
59a81a9d2704: Pull complete 
81e0f975b3fc: Pull complete 
Digest: sha256:25fc5423238b6f3a1fc967fda3f6a9212846aeb4a720327ef61c8ccff52dbbe2
Status: Downloaded newer image for authelia/authelia:latest
Digest: $argon2id$v=19$m=65536,t=3,p=4$hdf77UqbU+pUaDtc9K4oqg$sr9ZG+DeEHYNK/qcBQKm36nMd0y2/ML4/Mszbud3ymE
```

then restart the authelia pod
```
kubectl rollout restart statefulset -n authelia authelia
```


by default it is not setup to send emails but output to a .txt file ( to change this set the smtp block )
```
❯ kubectl exec -n authelia authelia-0 cat /config/notification.txt
```

Once all is set update useAuthelia: false to useAuthelia: true in your 
```
deploy/argocd/bootstrap-optional-apps/values.yaml
```
and push changes to github