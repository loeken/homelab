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
variable "ssh_private_key" {
  type = string
  default = "~/.ssh/id_ed25519"
  description = "The path to the private key to use when connecting to the server via SSH."
}
variable "shared_media_disk_device" {
  type = string
  default = "vdb"
  description = "the device name of the default media disk class"
}
variable "partition_external_shared_media_disk" {
  type = string
  default = "false"
  description = "if we use/format the external media disk"
}
