provider "helm" {
  kubernetes {
    config_path = "${path.module}/${var.kubeconfig_location}/kubeconfig"
  }
}

resource "helm_release" "argocd" {
  name       = "argocd"
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"
  create_namespace = true 
  namespace = "argocd"
  version = "5.20.4"
  
  values = [
    "${file("argocd-values.yaml")}"
  ]

  lifecycle {
    ignore_changes = [
      metadata
    ]
  }
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
    "${file("../../argocd/bootstrap-core-apps/values.yaml")}"
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
    "${file("../../argocd/bootstrap-core-apps/values.yaml")}"
  ]
}
