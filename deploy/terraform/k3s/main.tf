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
      "sudo apt install -y curl"
    ]
  }
  depends_on = [
    null_resource.copy_public_key
  ]
}
resource null_resource bootstrap-k3s {
  provisioner "local-exec" {
    command = <<EOT
    cd ../../../tmp
    k3sup install \
      --ip ${var.ssh_server_address} \
      --user ${var.ssh_username} \
      --cluster \
      --k3s-version ${var.kubernetes_version} \
      --k3s-extra-args '--disable=traefik --node-external-ip=${var.k3s_external_ip} --advertise-address=${var.ssh_server_address} --node-ip=${var.ssh_server_address}'
    EOT
  }
  depends_on = [
    null_resource.sudo_setup
  ]
}

resource "null_resource" "upload_ips" {
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = "${var.ssh_username}"
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
  depends_on = [
    null_resource.bootstrap-k3s
  ]
}

resource "null_resource" "nfs_server" {
  count = var.storage == "local-path" ? 1 : 0
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = "${var.ssh_username}"
    private_key = file("${var.ssh_private_key}")
  }
  
  provisioner "remote-exec" {
    inline = [
      "sudo apt update -y",
      "DEBIAN_FRONTEND=noninteractive sudo apt install -y nfs-kernel-server curl",
      "sudo mkdir -p /mnt/data",
      "echo '/mnt/data ${var.ssh_server_address}/32(rw,all_squash,anonuid=1000,anongid=1000)' | sudo tee /etc/exports",
      "sudo chown -R ${var.ssh_username}:${var.ssh_username} /mnt/data",
      "sudo systemctl restart nfs-kernel-server",
      "sudo sysctl fs.inotify.max_user_instances = 512"
    ]
  }

  depends_on = [
    null_resource.upload_ips
  ]
}
