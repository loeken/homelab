# Gitops & Argocd Concepts

![gitops.svg](Excalidraw/gitops.svg)

GitOps is a method of managing and updating computer systems using Git, which helps track changes in files. In GitOps, you define the desired state of your system, like the configuration and settings, in a Git repository.

Argo CD is a tool that helps implement GitOps for Kubernetes, a platform for running and managing containers. With Argo CD, you store the desired state of your Kubernetes environment in a Git repository, and Argo CD continuously monitors this repository. When there's a change in the repository, Argo CD automatically updates your Kubernetes environment to match the new desired state.

In summary, GitOps with Argo CD means using Git to store and manage the desired state of your Kubernetes environment, and Argo CD ensures your environment is always in sync with the repository. This setup makes it easier to collaborate, track changes, and maintain a consistent system state.