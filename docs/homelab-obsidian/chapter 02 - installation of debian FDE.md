# Installation of Debian FDE

we ll start by grabbing the last version of debian's netinstall, we ll use the non free variant which should increase your chance of having all drivers/firmware.

```bash
wget https://cdimage.debian.org/cdimage/unofficial/non-free/cd-including-firmware/11.6.0+nonfree/amd64/iso-cd/firmware-11.6.0-amd64-netinst.iso
```

and we prepare a bootable usb stick. first identify the stick. plug it in then run this command:
```bash
sudo dmesg -T |grep sd
[sudo] password for loeken: 
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] 30752848 512-byte logical blocks: (15.7 GB/14.7 GiB)
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Write Protect is off
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Mode Sense: 03 00 00 00
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] No Caching mode page found
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Assuming drive cache: write through
[Sat Apr 22 23:17:57 2023] sd 0:0:0:0: [sda] Attached SCSI removable disk
```

so this is a 16 GB USB stick called /dev/sda so lets write the image to it
```bash
sudo dd if=firmware-11.6.0-amd64-netinst.iso of=/dev/sda bs=4M status=progress
sudo sync
```

- During installation turn on disk encryption
- When asked to enter a root password do not enter one. This will setup user and add the non-root user to the sudoers group.