provider "helm" {
  kubernetes {
    config_path = "${path.module}/${var.kubeconfig_location}/kubeconfig"
  }
}

resource "helm_release" "sealed-secrets" {
  name       = "sealed-secrets-controller"
  chart      = "sealed-secrets"
  create_namespace = true 
  namespace = "kube-system"
  repository = "https://bitnami-labs.github.io/sealed-secrets/"
  version = "2.8.2"
}
