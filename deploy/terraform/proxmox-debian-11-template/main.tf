terraform {
  required_providers {
    proxmox = {
      source = "telmate/proxmox"
    }
  }
}
resource "null_resource" "create_template" {
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }

  provisioner "remote-exec" {
    inline = [
      "sudo apt update -y",
      "sudo apt install -y libguestfs-tools",
      "sudo qm destroy 999",
      "sudo qm destroy 998",
      "sudo rm -rf debian-11-nocloud-amd64.qcow2",
      "sudo wget https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-nocloud-amd64.qcow2",
      "sudo qm create 998 -name template-vm -memory 1024 -net0 virtio,bridge=${var.bridge} -cores 1 -sockets 1",

      # Import the OpenStack disk image to Proxmox storage
      "sudo qm importdisk 998 debian-11-nocloud-amd64.qcow2 local",

      # Attach the disk to the virtual machine
      "sudo qm set 998 -scsihw virtio-scsi-pci -virtio0  /var/lib/vz/images/998/vm-998-disk-0.raw",

      # Add a serial output
      "sudo qm set 998 -serial0 socket",

      # configure cloud init networking
      "sudo qm set 998 --ipconfig0 ip=dhcp",

      # Set the bootdisk to the imported Openstack disk
      "sudo qm set 998 -boot c -bootdisk virtio0",

      # Enable the Qemu agent
      "sudo qm set 998 -agent 1",

      # Allow hotplugging of network, USB and disks
      "sudo qm set 998 -hotplug disk,network,usb",

      # Add a single vCPU (for now)
      "sudo qm set 998 -vcpus 1",

      # Add a video output
      "sudo qm set 998 -vga qxl",

      # Set a second hard drive, using the inbuilt cloudinit drive
      "sudo qm set 998 -ide2 local:cloudinit",

      # ssh key for cloud init TODO fix bug
      # "echo $(cat ${var.ssh_public_key}) > mykey.pub",
      # "sudo qm set 998 --sshkey mykey.pub",

      # Resize the primary boot disk (otherwise it will be around 2G by default)
      # This step adds another 8G of disk space, but change this as you need to
      "sudo qm resize 998 virtio0 +8G",

      ###
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --install qemu-guest-agent,cloud-init,cloud-initramfs-growroot,open-iscsi,lsscsi,sg3-utils,multipath-tools,scsitools,nfs-common,lvm2",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --append-line '/etc/multipath.conf:defaults {'",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --append-line '/etc/multipath.conf:    user_friendly_names yes'",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --append-line '/etc/multipath.conf:    find_multipaths yes'",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --append-line '/etc/multipath.conf:}'",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --run-command 'systemctl enable open-iscsi.service'",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --run-command 'echo fs.inotify.max_user_instances = 8192 >> /etc/sysctl.conf'",
      "sudo virt-customize -a /var/lib/vz/images/998/vm-998-disk-0.raw --run-command 'echo fs.inotify.max_user_watches = 524288 >> /etc/sysctl.conf'",

      # Convert the VM to the template
      "sudo qm clone 998 999 --name template --full true",

      # convert the clone to a template
      "sudo qm template 999",
      "sleep 10",
    ]
  }
}