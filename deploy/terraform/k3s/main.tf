terraform {
  required_providers {
    proxmox = {
      source = "telmate/proxmox"
    }
    argocd = {
      source = "oboukili/argocd"
    }
    kubectl = {
      source  = "gavinbunney/kubectl"
    }
  }
}
provider "proxmox" {
  pm_api_url = var.pm_api_url
  pm_password = var.pm_password
  pm_user = var.pm_user
  pm_tls_insecure = "true"
} 
resource "proxmox_vm_qemu" "k3s-master-01" {
  agent = 1
  onboot = true
  name = var.k3s_nodename_01
  target_node = var.proxmox_node_name
  clone = "template"
  full_clone = true
  os_type = "cloud-init"
  sockets = 1
  cores = var.cores_k3s_01
  memory = var.memory_k3s_01
  scsihw            = "virtio-scsi-pci"
  ipconfig0 = "ip=dhcp"
  sshkeys = var.ssh_public_key
  ciuser = "debian"
  disk {
    type            = "virtio"
    storage = "local"
    size = var.k3s_disksize_01
  }
  lifecycle {
     ignore_changes = [
       network
     ]
  }
  network {
    model = "virtio"
    bridge = var.k3s_bridge_01
    macaddr = var.k3s_macaddr_01
  }
  provisioner "local-exec" {
    command = <<EOT
    k3sup install \
      --ip ${proxmox_vm_qemu.k3s-master-01.default_ipv4_address} \
      --user debian \
      --cluster \
      --k3s-version ${var.kubernetes_version} \
      --k3s-extra-args '--no-deploy=traefik --node-external-ip=${var.k3s_external_ip_node_01} --advertise-address=${var.k3s_internal_ip_node_01} --node-ip=${var.k3s_internal_ip_node_01}'
    EOT
  }
}

resource "proxmox_vm_qemu" "k3s-master-02" {
  agent = 1
  onboot = true
  name = var.k3s_nodename_02
  target_node = var.proxmox_node_name
  clone = "template"
  full_clone = true
  os_type = "cloud-init"
  sockets = 1
  cores = var.cores_k3s_02
  memory = var.memory_k3s_02
  scsihw            = "virtio-scsi-pci"
  ipconfig0 = "ip=dhcp"
  sshkeys = var.ssh_public_key
  ciuser = "debian"
  disk {
    type            = "virtio"
    storage = "local"
    size = var.k3s_disksize_02
  }
  lifecycle {
     ignore_changes = [
       network
     ]
  }
  network {
    model = "virtio"
    bridge = var.k3s_bridge_02
    macaddr = var.k3s_macaddr_02
  }
  provisioner "local-exec" {
    command = <<EOT
    k3sup join \
      --ip  ${proxmox_vm_qemu.k3s-master-02.default_ipv4_address} \
      --user debian \
      --server \
      --server-user debian \
      --server-ip ${proxmox_vm_qemu.k3s-master-01.default_ipv4_address} \
      --k3s-version ${var.kubernetes_version} \
      --k3s-extra-args '--no-deploy=traefik --node-external-ip=${var.k3s_external_ip_node_02} --advertise-address=${var.k3s_internal_ip_node_02} --node-ip=${var.k3s_internal_ip_node_02}'
    EOT
  }
  depends_on = [
    proxmox_vm_qemu.k3s-master-01
  ]
}
resource "proxmox_vm_qemu" "k3s-master-03" {
  agent = 1
  onboot = true
  name = var.k3s_nodename_03
  target_node = var.proxmox_node_name
  clone = "template"
  full_clone = true
  os_type = "cloud-init"
  sockets = 1
  cores = var.cores_k3s_03
  memory = var.memory_k3s_03
  scsihw            = "virtio-scsi-pci"
  ipconfig0 = "ip=dhcp"
  sshkeys = var.ssh_public_key 
  ciuser = "debian"
  disk {
    type            = "virtio"
    storage = "local"
    size = var.k3s_disksize_03
  }
  lifecycle {
     ignore_changes = [
       network
     ]
  }
  network {
    model = "virtio"
    bridge = var.k3s_bridge_03
    macaddr = var.k3s_macaddr_03
  }
  provisioner "local-exec" {
    command = <<EOT
    k3sup join \
      --ip  ${proxmox_vm_qemu.k3s-master-03.default_ipv4_address} \
      --user debian \
      --server \
      --server-user debian \
      --server-ip ${proxmox_vm_qemu.k3s-master-01.default_ipv4_address} \
      --k3s-version ${var.kubernetes_version} \
      --k3s-extra-args '--no-deploy=traefik --node-external-ip=${var.k3s_external_ip_node_03} --advertise-address=${var.k3s_internal_ip_node_03} --node-ip=${var.k3s_internal_ip_node_03}'
    EOT
  }
  depends_on = [
    proxmox_vm_qemu.k3s-master-02
  ]
}
provider "helm" {
  kubernetes {
    config_path = "${path.module}/kubeconfig"
  }
}

resource "helm_release" "cert-manager" {
  name        = "cert-manager"
  chart       = "cert-manager"
  repository  = "https://charts.jetstack.io"
  create_namespace = true
  set {
    name = "installCRDs"
    value = "true"
  }
  namespace = "cert-manager"
  depends_on = [
    proxmox_vm_qemu.k3s-master-01,
    proxmox_vm_qemu.k3s-master-02,
    proxmox_vm_qemu.k3s-master-03
  ]
}

/*
//@TODO test if can be removed with proxmox has hypervisor iirc issues with ovh's managed kubernetes
resource "time_sleep" "wait_30_seconds" {
  depends_on = [helm_release.cert-manager]
  create_duration = "30s"
}

# This resource will create (at least) 30 seconds after null_resource.previous
resource "null_resource" "next" {
  depends_on = [time_sleep.wait_30_seconds]
}
*/
resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  create_namespace = true 
  namespace = "argocd"
  version = "4.4.1"
  
  values = [
    "${file("argocd-values.yaml")}"
  ]

  depends_on = [
    //null_resource.next,
    proxmox_vm_qemu.k3s-master-01,
    proxmox_vm_qemu.k3s-master-02,
    proxmox_vm_qemu.k3s-master-03
  ]
  lifecycle {
    ignore_changes = [
      metadata
    ]
  }
}
resource "helm_release" "sealed_secrets" {
  name       = "sealed-secrets-controller"
  repository = "https://bitnami-labs.github.io/sealed-secrets/"
  chart      = "sealed-secrets"
  create_namespace = true 
  namespace = "kube-system"
  version = "2.1.2"
  depends_on = [
    proxmox_vm_qemu.k3s-master-01,
    proxmox_vm_qemu.k3s-master-02
  ]
}

resource "helm_release" "bootstrap-core-apps" {
  name       = "bootstrap-core-apps"
  chart      = "./../../helm/bootstrap-core-apps"
  create_namespace = true 
  namespace = "argocd"

  depends_on = [
    helm_release.argocd
  ]
  values = [
    "${file("../../helm/bootstrap-core-apps/values.yaml")}"
  ]
}
resource "helm_release" "bootstrap-optional-apps" {
  name       = "bootstrap-optional-apps"
  chart      = "./../../helm/bootstrap-optional-apps"
  create_namespace = true 
  namespace = "argocd"

  depends_on = [
    helm_release.argocd
  ]
  values = [
    "${file("../../helm/bootstrap-optional-apps/values.yaml")}"
  ]
}
