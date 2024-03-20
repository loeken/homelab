## Home Assistant

### install HACS
a good guide can be found (here)[https://hacs.xyz/docs/configuration/start]

```
❯ kubectl exec -it -n home-assistant home-assistant-0 ash
kubectl exec [POD] [COMMAND] is DEPRECATED and will be removed in a future version. Use kubectl exec [POD] -- [COMMAND] instead.
/config # cat configuration.yaml 
# Loads default set of integrations. Do not remove.
default_config:

# Load frontend themes from the themes folder
frontend:
  themes: !include_dir_merge_named themes

automation: !include automations.yaml
script: !include scripts.yaml
scene: !include scenes.yaml

# my edits:
http:
  use_x_forwarded_for: true
  trusted_proxies:
    - 10.42.0.0/16  # Adjust the IP range according to your Kubernetes cluster's IP range
  ip_ban_enabled: false
/config # wget -O - https://get.hacs.xyz | bash -
Connecting to get.hacs.xyz (104.21.5.2:443)
Connecting to raw.githubusercontent.com (185.199.110.133:443)
writing to stdout
-                    100% |********************************************************************************************************************************************************************|  2473  0:00:00 ETA
written to stdout
INFO: Trying to find the correct directory...
INFO: Found Home Assistant configuration directory at '/config'
INFO: Creating custom_components directory...
INFO: Changing to the custom_components directory...
INFO: Downloading HACS
Connecting to github.com (140.82.121.4:443)
Connecting to github.com (140.82.121.4:443)
Connecting to objects.githubusercontent.com (185.199.110.133:443)
saving to 'hacs.zip'
hacs.zip             100% |********************************************************************************************************************************************************************|  730k  0:00:00 ETA
'hacs.zip' saved
INFO: Creating HACS directory...
INFO: Unpacking HACS...
INFO: Removing HACS zip file...
INFO: Installation complete.

INFO: Remember to restart Home Assistant before you configure it
```

either restart the pod, run kill 1 in this terminal or if you can login to the ui some other way, head to Developer Tools -> Restart

then over to Configuration -> Devices & Services

click Add Integration at bottom right. Search for HACS 
check all the stupid checkboxes and submit
head over to github and enter the key home assistant showed you, assign it to a room


while we have a shell open let's also add the redirect rules

```
cd /config
mkdir themes
apk add nano
nano /config/configuration.yaml

http:
    server_host: 0.0.0.0
    ip_ban_enabled: true
    login_attempts_threshold: 5
    use_x_forwarded_for: true
    trusted_proxies:
    # Pod CIDR
    - 10.42.0.0/16
    # Node CIDR
    - 172.16.137.0/24
```

without this block you will encouter a redirect loop.
