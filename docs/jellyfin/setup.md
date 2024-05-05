## Jellyfin

login to website create user account
Import the two media folders we created in rtorrent ( there is a shared storage between those arr vms )


https://www.authelia.com/integration/openid-connect/jellyfin/

follow docs


in jellyfin:

- general settings -> branding custom login button:
```
<form action="https://jellyfin.example.com/sso/OID/start/authelia">
  <button class="raised block emby-button button-submit">
    Sign in with Authelia
  </button>
</form>
```

as OID endpoint set:
https://auth.example.com/.well-known/openid-configuration/


- custom css code:
```
a.raised.emby-button {
  padding: 0.9em 1em;
  color: inherit !important;
}

.disclaimerContainer {
  display: block;
}
```