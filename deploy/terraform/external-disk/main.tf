resource "null_resource" "nfs_server" {
  count = var.partition_external_shared_media_disk == "true" ? 1 : 0
  connection {
    type     = "ssh"
    host     = var.ssh_server_address
    user     = var.ssh_username
    private_key = file("${var.ssh_private_key}")
  }
  
  provisioner "remote-exec" {
    inline = [
      "if ! grep -qs '/dev/sda' /proc/mounts; then",
      "  sudo apt update -y",
      "  DEBIAN_FRONTEND=noninteractive sudo apt install -y parted",
      "  sudo partprobe /dev/${var.shared_media_disk_device}",
      "  sudo parted --script /dev/${var.shared_media_disk_device} mklabel gpt",
      "  sudo parted --script /dev/${var.shared_media_disk_device} mkpart primary ext4 0% 100%",

      "  sudo mkfs.ext4 -F /dev/${var.shared_media_disk_device}1",
      "  sudo mkdir -p /mnt/data",
      "  echo '/dev/${var.shared_media_disk_device}1 /mnt/data ext4 rw,discard,errors=remount-ro 0 1' | sudo tee -a /etc/fstab",
      "  sudo mount -a",
      "  sudo pvesm add dir external-disk --path /mnt/data --content iso,rootdir,images,backup",
      "fi"
    ]
  }
}

