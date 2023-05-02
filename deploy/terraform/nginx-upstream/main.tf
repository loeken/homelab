resource "proxmox_vm_qemu" "upstream-vm" {
  agent = 1
  onboot = true
  name = "nginx-upstream"
  target_node = var.proxmox_node_name
  clone = "template"
  full_clone = true
  os_type = "cloud-init"
  sockets = 1
  cores = 1
  memory = 512
  scsihw            = "virtio-scsi-pci"
  ipconfig0 = "ip=dhcp"
  sshkeys = file("${var.ssh_public_key}")
  ciuser = var.ssh_username
  qemu_os = "l26"
  vcpus = var.cores_k3s
  disk {
    type    = "virtio"
    storage = "local"
    size = "5G"
  }
    # Second disk
  lifecycle {
    ignore_changes = [
      network
    ]
  }
  network {
    model = "virtio"
    bridge = var.bridge
  }
}
resource "null_resource" "k3s-installation" {
  # This resource will only be executed after the K3s virtual machine is up and running
  depends_on = [proxmox_vm_qemu.k3s-vm]

  provisioner "local-exec" {
    working_dir = "${path.module}/${var.kubeconfig_location}"
    command = <<EOT
      #!/bin/bash
      sudo apt update
      sudo apt install nginx -t
      echo "upstream backend {
    server 172.16.137.200:30080;
}

upstream backend_ssl {
    server 172.16.137.200:30443;
}

server {
    listen 80;
    server_name _; # replace with your domain name

    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}

server {
    listen 443;
    server_name _; # replace with your domain name

    location / {
        proxy_pass https://backend_ssl;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        # Enable SSL passthrough
        proxy_ssl_session_reuse off;
        proxy_ssl_name $host;
        proxy_ssl_server_name on;
    }
}" > /etc/nginx/server.conf
    EOT
  }
}