## deploy/terraform/k3s section
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
variable "kubernetes_version" {
    type = string
    default = "v1.25.3+k3s1"
    description = "which version of k3s to install, usually 1 versions behind the latest"
}
variable "external_ip" {
    type = string
    default = "1.2.3.4"
    description = "sets the external ip address, a script to update ips and restart k3s is also uploaded to the vm"
}
variable "cores_k3s" {
    type = string
    default = "10"
    description = "defines to the toal core count that are assigned to the k3s vm"
}
variable "memory_k3s" {
    type = string
    default = "28672"
    description = "defines to the toal MB that are assigned to the k3s vm"
}
variable "proxmox_node_name" {
    type = string
    default = "beelink-sei12"
    description = "the name of the proxmox node, equals the hostname of the proxmox server, in case of cluster defines on which node the vm is created"
}
variable "macaddr" {
    type = string
    default = "8E:AB:AB:4C:CE:A0"
    description = "you can use this mac address to achive static ips by setting static mappings in your dhcp server"
}
variable "disksize" {
    type = string
    default = "100G"
    description = "size of the virtual machine primary disk, you can add more disks later on"
}
variable "bridge" {
  type = string
  default = "vmbr0"
  description = "This is the name of the bridge by default: vmbr0."
}
variable "proxmox_vm_name" {
    type = string
    default = "k3s-beelink"
    description = "the name of the proxmox node, equals the hostname of the proxmox server, in case of cluster defines on which node the vm is created"
}
variable "kubeconfig_location" {
  type = string
  default = "../../../tmp/"
  description = "the location of the kubeconfig, can be overwriten to be used with minikube"
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

