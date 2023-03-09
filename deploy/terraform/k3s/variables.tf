## deploy/terraform/k3s section

variable "ssh_username" {
  type = string
  default = "user"
  description = "The username to use when connecting to the server via SSH."
}
variable "ssh_private_key" {
  type = string
  default = "~/.ssh/id_rsa"
  description = "The path to the private key to use when connecting to the server via SSH."
}
variable "kubernetes_version" {
    type = string
    default = "v1.25.3+k3s1"
    description = "which version of k3s to install, usually 1 versions behind the latest"
}
variable "k3s_external_ip" {
    type = string
    default = "1.2.3.4"
    description = "sets the external ip address, a script to update ips and restart k3s is also uploaded to the vm"
}
variable "ssh_password" {
  type = string
  default = "demotime"
  description = "The password to use when connecting to the server via SSH."
}
variable "ssh_public_key" {
  type = string
  default = "~/.ssh/id_rsa.pub"
  description = "The path to the public key to use when connecting to the server via SSH."
}
variable "ssh_server_address" {
  type = string
  default = "localhost"
  description = "The address of the server to connect to via SSH."
}
variable "storage" {
  type = string
  default = "local-path"
  description = "the default storage class"
}
variable "pci_device" {
  type = string
  default = "0000:02:00.0"
  description = "the default storage class"
}
variable "shared_media_disk_device" {
  type = string
  default = "vdb"
  description = "the device name of the default media disk class"
}
variable "shared_media_disk_size" {
  type = string
  default = "400G"
  description = "size of the shared media disk"
}

