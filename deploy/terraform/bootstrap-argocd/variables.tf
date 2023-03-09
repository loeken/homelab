## deploy/terraform/k3s section

variable "kubeconfig_location" {
  type = string
  default = "../../../tmp/"
  description = "the location of the kubeconfig, can be overwriten to be used with minikube"
}
