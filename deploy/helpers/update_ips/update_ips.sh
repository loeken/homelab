#!/bin/bash
newip=`curl ifconfig.co -4`
if [ "$newip" == "" ]; then
  exit 
fi
oldip=`cat /opt/oldip`
if [ "$oldip" == "" ]; then
  echo $newip > /opt/oldip
fi
if [ "$oldip" != "$newip" ]; then
  sed -i "s/$oldip/$newip/g" /etc/systemd/system/k3s.service
  systemctl daemon-reload
  systemctl restart k3s
  echo $newip > /opt/oldip
fi
