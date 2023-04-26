## setup kasten k10

## initial setup:

```
kubectl -n kasten-io port-forward svc/gateway 8080:8000
```

then visit http://localhost:8080/k10/#/ enter email/company name and enter.

head to settings -> location and create a new profile

-> s3 bucket

- profile name home
- aws access key
- aws secret key
- region ( ireland cheapest for me in europe )
- bucketname: pickwhatyouwant


#### disaster recovery
in order to backup it's own state and starta recovery from offsite location ( S3 ) we ll enable the disaster recovery with the newly added profile ( home )

make sure to save your k10 backup id
592067ee-8284-459c-89da-1395af87eec4