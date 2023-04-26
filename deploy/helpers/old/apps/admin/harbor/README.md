## Harbor

login with admin / harborAdminPassword you set in argocd-optional.yaml

#### ldap
Administration -> Configuration

- Auth Mode: Ldap
- LDAP Search DN: cn=admin,dc=ldap,dc=example,dc=com
- LDAP Search Password: homelab
- LDAP Base DN: ou=User,dc=ldap,dc=example,dc=com
- LDAP Filter: (&(objectClass=inetOrgPerson))
- LDAP UID: uid
- LDAP Scope: Subtree
- LDAP Group Base DN: ou=Group,dc=ldap,dc=example,dc=com
- LDAP Group Filter: (&(objectClass=groupOfNames))
- LDAP Group GID: cn
- LDAP Group Admin DN: cn=admins,ou=group,dc=ldap,dc=example,dc=com
- LDAP Group Membership: member
- LDAP Scope: Subtree
- LDAP verify Certificate: unchecked

#### Email
- smtp: mail.example.com
- Server Port: 465
- Email Username: harbor@example.com
- Email Password: homelab
- Email From: harbor <harbor@internetz.me>
- Email SSL: checked
- Verify Certificate: checked

#### notes
while login works the admin group does not work properly user remains in state "Unknown" manual promotion works

