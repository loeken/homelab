## deploy/terraform/proxmox-debian-11-template section
variable "ssh_server_address" {
  type = string
  default = "localhost"
  description = "The address of the server to connect to via SSH."
}
variable "ssh_username" {
  type = string
  default = "user"
  description = "The username to use when connecting to the server via SSH."
}
variable "root_password" {
  type = string
  default = "topsecure"
  description = "This is the root password, not used for ssh - mostly to login to proxmox's webui."
}
variable "ssh_public_key" {
  type = string
  default = "~/.ssh/id_rsa.pub"
  description = "The path to the public key to use when connecting to the server via SSH."
}
variable "ssh_private_key" {
  type = string
  default = "~/.ssh/id_rsa"
  description = "The path to the private key to use when connecting to the server via SSH."
}
variable "interface" {
  type = string
  default = "enp3s0"
  description = "This is the name of the network interface used by proxmox, this will be the base for the bridge vmbr0."
}
variable "bridge" {
  type = string
  default = "vmbr0"
  description = "This is the name of the bridge by default: vmbr0."
}