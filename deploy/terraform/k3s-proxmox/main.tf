terraform {
  required_providers {
    proxmox = {
      source = "telmate/proxmox"
    }
  }
}
provider "proxmox" {
  pm_api_url = "https://${var.ssh_server_address}:8006/api2/json" 
  pm_password = var.root_password
  pm_user = "root@pam"
  pm_tls_insecure = "true"
} 
resource "proxmox_vm_qemu" "k3s-vm" {
  agent = 1
  onboot = true
  name = var.proxmox_vm_name
  target_node = var.proxmox_node_name
  clone = "template"
  full_clone = true
  os_type = "cloud-init"
  sockets = 1
  cores = var.cores_k3s
  memory = var.memory_k3s
  scsihw            = "virtio-scsi-pci"
  ipconfig0 = "ip=dhcp"
  sshkeys = file("${var.ssh_public_key}")
  ciuser = var.ssh_username
  qemu_os = "l26"
  vcpus = var.cores_k3s
  disk {
    type    = "virtio"
    storage = "local"
    size = "50G"
  }
    # Second disk
  dynamic "disk" {
    for_each = var.partition_external_shared_media_disk ? [1] : []
    content {
      type    = "virtio"
      storage = var.partition_external_shared_media_disk ? "external-disk" : "local"
      size    = var.shared_media_disk_size
    }
  }
  lifecycle {
    ignore_changes = [
      network
    ]
  }
  network {
    model = "virtio"
    bridge = var.bridge
    macaddr = var.macaddr
  }
  provisioner "local-exec" {
    working_dir = "${path.module}/${var.kubeconfig_location}"
    command = <<EOT
      k3sup install \
      --ip ${proxmox_vm_qemu.k3s-vm.default_ipv4_address} \
      --ssh-key ${var.ssh_private_key} \
      --user ${var.ssh_username} \
      --cluster \
      --k3s-version ${var.kubernetes_version} \
      --k3s-extra-args '--disable=traefik --node-external-ip=${var.external_ip} --advertise-address=${proxmox_vm_qemu.k3s-vm.default_ipv4_address} --node-ip=${proxmox_vm_qemu.k3s-vm.default_ipv4_address}'
    EOT
  }
}

resource "null_resource" "upload_ips" {
  connection {
    type     = "ssh"
    host     = proxmox_vm_qemu.k3s-vm.default_ipv4_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  provisioner "file" {
    source     = "../../helpers/update_ips/update_ips.sh"
    destination = "/tmp/update-ips"
  }
  provisioner "remote-exec" {
    inline = [
      "sudo cp /tmp/update-ips /usr/local/bin/update-ips",
      "sudo chmod +x /usr/local/bin/update-ips"
    ]
  }
}

resource "null_resource" "nfs_server" {
  count = var.storage == "local-path" ? 1 : 0
  connection {
    type     = "ssh"
    host     = proxmox_vm_qemu.k3s-vm.default_ipv4_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  
  provisioner "remote-exec" {
    inline = [
      "sudo apt update -y",
      "DEBIAN_FRONTEND=noninteractive sudo apt install -y nfs-kernel-server parted",
      "sudo parted /dev/vdb mklabel msdos",
      "sudo parted /dev/vdb mkpart primary ext4 0% 100%",
      "sudo mkfs.ext4 -F /dev/vdb1",
      "sudo mkdir -p /mnt/data",
      "echo '/dev/vdb1 /mnt/data ext4 rw,discard,errors=remount-ro 0 1' | sudo tee -a /etc/fstab",
      "sudo mount -a",
      "echo '/mnt/data ${proxmox_vm_qemu.k3s-vm.default_ipv4_address}/32(rw,all_squash,anonuid=1000,anongid=1000)' | sudo tee /etc/exports",
      "sudo chown -R ${var.ssh_username}:${var.ssh_username} /mnt/data",
      "sudo systemctl restart nfs-kernel-server",
    ]
  }
}
