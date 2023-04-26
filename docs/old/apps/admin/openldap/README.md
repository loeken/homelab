## setup openldap

#### Openldap

```
cp argocd-openldap.yaml.example argocd-openldap.yaml
cat argocd-openldap.yaml | kubeseal > argocd-openldap-encrypted.yaml 
kubectl apply -f argocd-openldap-encrypted.yaml

cp argocd-openldap-ltb-passwd.yaml.example argocd-openldap-ltb-passwd.yaml
cat argocd-openldap-ltb-passwd.yaml | kubeseal > argocd-openldap-ltb-passwd-encrypted.yaml 
kubectl apply -f argocd-openldap-ltb-passwd-encrypted.yaml
```

//@TODO section needs update

Now go to your ldap admin dashboard and login via
```
cn=admin,dc=ldap,dc=example,dc=com
```
and the admin password specified in the deploy/mysecrets/argocd-openldap.yaml
```
version: 1

# Entry 1: ou=Group,dc=ldap,dc=example,dc=com
dn: ou=Group,dc=ldap,dc=example,dc=com
objectclass: organizationalUnit
objectclass: top
ou: Group

# Entry 2: ou=User,dc=ldap,dc=example,dc=com
dn: ou=User,dc=ldap,dc=example,dc=com
objectclass: organizationalUnit
objectclass: top
ou: User


# Entry 3: uid=loeken,ou=User,dc=ldap,dc=example,dc=com
dn: uid=loeken,ou=User,dc=ldap,dc=example,dc=com
cn: loeken
displayname: loeken
mail: loeken@example.me
objectclass: inetOrgPerson
objectclass: top
sn: Filewalker
uid: loeken
userpassword: {CRYPT}$5$WtLavDZC$ScD.IMUJdgRhMZMtAlyYtbqezxsQ2qfWTVbQOFo5tg4

# Entry 4: cn=admins,ou=Group,dc=ldap,dc=example,dc=com
dn: cn=admins,ou=Group,dc=ldap,dc=example,dc=com
cn: admins
member: uid=loeken,ou=User,dc=ldap,dc=example,dc=com
objectclass: groupOfNames
objectclass: top

# Entry 5: cn=users,ou=Group,dc=ldap,dc=example,dc=com
dn: cn=users,ou=Group,dc=ldap,dc=example,dc=com
cn: users
member: uid=loeken,ou=User,dc=ldap,dc=example,dc=com
objectclass: groupOfNames
objectclass: top

```

the userpassword resolves to "homelab"

#### creating a new user

create a user in the Ou=User then head to the Ou=Group section to add the user to a group