this folder containers a simple script that allows to change the external node ip of k3s by editing the systemd startup file.

Schedule this script as a cronjobn to run every X minutes, it will do 1 call to ifconfig.co get the external ipv4 and reload k3s.