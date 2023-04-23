# main.tf

## turns debian 11 (with full disk encryption) into proxmox 

provider "null" { }

resource "null_resource" "copy_public_key" {
  provisioner "local-exec" {
    command = <<EOF
     SSHPASS='${var.ssh_password}' sshpass -e ssh -o 'StrictHostKeyChecking no' ${var.ssh_username}@${var.ssh_server_address} "mkdir -p ~/.ssh && echo '$(cat ${var.ssh_public_key})' > ~/.ssh/authorized_keys"
    EOF
  }
}
resource "null_resource" "sudo_setup" {
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  provisioner "remote-exec" {
    inline = [
      "echo '${var.ssh_password}' | sudo -S -k sh -c 'echo \"${var.ssh_username} ALL=(ALL) NOPASSWD:ALL\" >/etc/sudoers.d/${var.ssh_username}'",
    ]
  }
  depends_on = [
    null_resource.copy_public_key
  ]
}
resource "null_resource" "migrate" {
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  provisioner "remote-exec" {
    inline = [
      "sudo apt update -y",
      "sudo sed -i 's/127.0.1.1/${var.ssh_server_address}/g' /etc/hosts",
      "hostname --ip-address",
      "echo \"deb [arch=amd64] http://download.proxmox.com/debian/pve bullseye pve-no-subscription\" | sudo tee /etc/apt/sources.list.d/pve-install-repo.list",
      "sudo wget https://enterprise.proxmox.com/debian/proxmox-release-bullseye.gpg -O /etc/apt/trusted.gpg.d/proxmox-release-bullseye.gpg",
      "sudo apt update -y",
      "sudo apt full-upgrade -y",
      "sudo apt install -y pve-kernel-5.19",
      "sudo sed -i 's/quiet/quiet intel_iommu=on/g' /etc/default/grub",
      "echo vfio | sudo tee -a /etc/modules",
      "echo vfio_iommu_type1 | sudo tee -a /etc/modules",
      "echo vfio_pci | sudo tee -a /etc/modules",
      "echo vfio_virqfd | sudo tee -a /etc/modules",
      "sudo DEBIAN_FRONTEND=noninteractive apt remove 'linux-image-5.*' linux-image-amd64 -y",
      "sudo update-initramfs -u -k 5.19",
      "sudo update-grub",
      "if ! uname -r | grep -q pve; then sudo shutdown -r now; fi",
    ]
  }
  depends_on = [
    null_resource.sudo_setup
  ]
}

resource "null_resource" "wait_for_reboot" {
  provisioner "local-exec" {
    command = "sleep 3 && until ping -c1 ${var.ssh_server_address}; do sleep 1; done"
  }
  depends_on = [
    null_resource.migrate
  ]
}
resource "null_resource" "pve_setup" {
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  provisioner "remote-exec" {
    inline = [
      "uname -a",
      "echo \"postfix postfix/mailname string my.hostname.example\" | sudo debconf-set-selections",
      "echo \"postfix postfix/main_mailer_type string 'Internet Site'\" | sudo debconf-set-selections",
      "DEBIAN_FRONTEND=noninteractive sudo apt install -y lvm2 proxmox-ve postfix open-iscsi ifupdown2 bridge-utils",
      "sudo apt remove -y os-prober",
      "sudo apt autoremove -y",
      "sudo rm -rf /etc/apt/sources.list.d/pve-enterprise.list"
    ]
  }
  depends_on = [
    null_resource.wait_for_reboot
  ]
}
resource "null_resource" "set_root_password" {
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  provisioner "remote-exec" {
    inline = [
      "echo 'root:${var.root_password}' | sudo chpasswd",
      "sudo rm -rf /etc/network/interfaces",
      "echo source /etc/network/interfaces.d/* | sudo tee -a /etc/network/interfaces",
      "echo auto lo | sudo tee -a /etc/network/interfaces",
      "echo iface lo inet loopback | sudo tee -a /etc/network/interfaces",
      "echo auto ${var.interface} | sudo tee -a /etc/network/interfaces",
      "echo iface ${var.interface} inet static | sudo tee -a /etc/network/interfaces",
      "echo auto ${var.bridge} | sudo tee -a /etc/network/interfaces",
      "echo iface ${var.bridge} inet static | sudo tee -a /etc/network/interfaces",
      "echo '  address ${var.ssh_server_address}/${var.ssh_server_netmask}' | sudo tee -a /etc/network/interfaces",
      "echo '  gateway ${var.ssh_server_gateway}' | sudo tee -a /etc/network/interfaces",
      "echo '  bridge-ports ${var.interface}' | sudo tee -a /etc/network/interfaces",
      "echo '  bridge-stp off' | sudo tee -a /etc/network/interfaces",
      "echo '  bridge-fd 0' | sudo tee -a /etc/network/interfaces",
      "echo '  bridge-vlan-aware yes' | sudo tee -a /etc/network/interfaces",
      "echo '  bridge-vids 2-4094' | sudo tee -a /etc/network/interfaces",
      "sudo rm -rf /etc/network/interfaces.new",
      "sudo ifreload --all"
    ]
  }
  depends_on = [
    null_resource.pve_setup
  ]
}
