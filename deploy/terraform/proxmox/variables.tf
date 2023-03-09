# variables.tf
variable "ssh_server_address" {
  type = string
  default = "localhost"
  description = "The address of the server to connect to via SSH."
}
variable "ssh_server_gateway" {
  type = string
  default = "172.16.137.254"
  description = "This is the gateway of the server."
}
variable "ssh_server_netmask" {
  type = string
  default = "24"
  description = "This is the netmask of the server`s network."
}
variable "ssh_username" {
  type = string
  default = "user"
  description = "The username to use when connecting to the server via SSH."
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
variable "ssh_private_key" {
  type = string
  default = "~/.ssh/id_rsa"
  description = "The path to the private key to use when connecting to the server via SSH."
}
variable "root_password" {
  type = string
  default = "topsecure"
  description = "This is the root password, not used for ssh - mostly to login to proxmox's webui."
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