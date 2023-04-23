<# Installation of Debian FDE

we ll start by grabbing the last version of debian's netinstall:

```bash
wget https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-11.6.0-amd64-netinst.iso
```

and we prepare a bootable usb stick. first identify the stick. plug it in then run this command:
```bash
sudo dmesg -T |grep sda
[sudo] password for loeken: 
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] 30752848 512-byte logical blocks: (15.7 GB/14.7 GiB)
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Write Protect is off
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Mode Sense: 03 00 00 00
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] No Caching mode page found
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Assuming drive cache: write through
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Attached SCSI removable disk
```

so this is a 16 GB stick called sda so lets send the image to it
```bash
sudo dd if=debian-11.6.0-amd64-netinst.iso of=/dev/sda bs=4M status=progress
sudo sync
```

- During installation turn on disk encryption
- When asked to enter a root password do not enter one. This will setup user and add the non-root user to the sudoers group.