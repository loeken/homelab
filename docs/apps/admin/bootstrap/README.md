## bootstrapping argocd applications
all argocd apps defined in this repo are loaded via 2 charts. First the bootstrap-core-apps, which installs the apps that other apps rely upon. these apps. the rest of the apps are loaded via the bootrap-optional-apps chart afterwards.



```bash
cd ~/Projects/private/homelab-private/deploy/helm/


cp bootstrap-core-apps/values.yaml.example bootstrap-core-apps/values.yaml

# for proxmox based installs
cp bootstrap-optional-apps/values.yaml.example bootstrap-optional-apps/values.yaml

# for minikube based installs
cp bootstrap-optional-apps/values.yaml.example-minikube bootstrap-optional-apps/values.yaml


nano bootstrap-core-apps/values.yaml
nano bootstrap-optional-apps/values.yaml
```
now edit these two values.yaml files to match your project then we git add/push to our private repo. these values.yaml's shouldnt contain passwords, whenever passwords are involved we aim to use a sealed secret ( existingSecret used ).


we still need to add 1 further values.yaml file - what i call "the global values.yaml for the optional apps". inside deploy/argocd/bootstrap-optional-apps you will find a values.yaml.example again we
```bash
cd ~/Projects/private/homelab-private/deploy/argocd/bootstrap-optional-apps
cp values.yaml.example values.yaml
nano values.yaml
```
the templates folder inside this folder are part of a local helm chart and contain the manifest for all "optional apps", the values.yaml in this folder can be used to pass variables such as your "domain name" to all other helm charts. this is the main file where you configure your apps, if you need to configure something else which you cannot find in this values.yaml create an issue on github.com/loeken/homelab

we've added a variables.tf before which was responsible for creating the proxmox template. now we edit the contents of the one responsible for bootstrapping k3s

```bash
cd ~/Projects/private/homelab-private/deploy/terraform/k3s
nano variables.tf
```

as this variables.tf contains the proxmox password and other information we won't add it to github ( it is ignored via .gitignore )
```bash
cd ~/Projects/private/homelab-private

git status
On branch main
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Untracked files:
  (use "git add <file>..." to include in what will be committed)
	deploy/argocd/bootstrap-optional-apps/values.yaml
	deploy/helm/bootstrap-core-apps/values.yaml
	deploy/helm/bootstrap-optional-apps/values.yaml
```
as you can see we created 3 files let's add them to git
```bash
git add .
git status
On branch main
Your branch is ahead of 'origin/main' by 1 commit.
  (use "git push" to publish your local commits)

Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
	new file:   deploy/argocd/bootstrap-optional-apps/values.yaml
	new file:   deploy/helm/bootstrap-core-apps/values.yaml
	new file:   deploy/helm/bootstrap-optional-apps/values.yaml
```
now we commit the changes
```bash
git commit -m "added my own configs"
[main 60c43d5] added my own configs
 4 files changed, 137 insertions(+)
 create mode 100644 deploy/argocd/bootstrap-optional-apps/values.yaml
 create mode 100644 deploy/helm/bootstrap-core-apps/values.yaml
 create mode 100644 deploy/helm/bootstrap-optional-apps/values.yaml
```
and send them to our private repository
```bash
git push
Krypton ▶ Requesting SSH authentication from phone
Krypton ▶ Success. Request Allowed ✔
Enumerating objects: 23, done.
Counting objects: 100% (23/23), done.
Delta compression using up to 8 threads
Compressing objects: 100% (15/15), done.
Writing objects: 100% (15/15), 1.34 KiB | 1.34 MiB/s, done.
Total 15 (delta 8), reused 0 (delta 0), pack-reused 0
remote: Resolving deltas: 100% (8/8), completed with 6 local objects.
To github.com:loeken/homelab-private
   53676ca..60c43d5  main -> main
```

now we can run the terraform scripts to create the k3s cluster
```bash
cd ~/Projects/private/homelab-private/deploy/terraform/k3s
terraform init
terraform plan
terraform apply
```

terraform should have completed by now and we can attempt to access argocd's webui
```
cd ~/Projects/private/homelab-private/deploy/terraform/k3s
export KUBECONFIG=$PWD/kubeconfig
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
kubectl -n argocd port-forward svc/argocd-server 8081:443
```

### 5.1. argocd sync bootstrap-core-app of apps
now we can visit http://localhost:8081 in our browser. at this stage the bootstrap-core-apps app will be in a failed state as it needs the private key of the deploy we created above to be able to pull from github. 

error message in argocd's webui:
```
rpc error: code = Unknown desc = error creating SSH agent: "SSH agent requested but SSH_AUTH_SOCK not-specified"
```
this is basically the first secret we'll send to the cluster. We don't want to send it to the cluster unencrypted, this is where kubeseal comes in, we basically send the secret yaml to the kubeseal controller, which encrypts it and returns us the encrypted file, we then apply the encrypted format to the cluster, whenever argocd rolls out code and needs this information the cluster will decrypt and access the secret.

```
cd ~/Projects/private/homelab/deploy/mysecrets
cp argocd-bootstrap-core-apps-repo.yaml.example argocd-bootstrap-core-apps-repo.yaml
nano argocd-bootstrap-core-apps-repo.yaml
```

replace the ssh key section with the private key we created above ( id_ed25519_homelab_private_deploy_key ).

```
cd ~/Projects/private/homelab-private/deploy/mysecrets
cat argocd-bootstrap-core-apps-repo.yaml | kubeseal > argocd-bootstrap-core-apps-repo-encrypted.yaml
kubectl apply -f argocd-bootstrap-core-apps-repo-encrypted.yaml
```
we display the contents of argocd-bootstrap-core-apps-repo.yaml use kubeseal to send it to the cluster and then save the encrypted output to argocd-bootstrap-core-apps-repo-encrypted.yaml then we use kubectl apply -f to send it to the cluster. You can also git add/commit/push the -encrypted.yaml to the repo

the rest of the core apps can by synced in one go, only external-dns will need another secret, containing the credentials ( such as cloudflare`s api key/token ) so the cluster can update external dns.

once this is done create a backup of your sealed secrets ( the main encryption keys )
```
cd deploy/mysecrets
kubectl -n kube-system get secret -l sealedsecrets.bitnami.com/sealed-secrets-key=active -o yaml | kubectl neat > allsealkeys.yml
```

[external-dns documents on creating api key/token](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/cloudflare.md)

```
cd ~/Projects/private/homelab-private/deploy/mysecrets
cp argocd-external-dns-cloudflare.yaml.example argocd-external-dns-cloudflare.yaml
cat argocd-external-dns-cloudflare.yaml | kubeseal > argocd-external-dns-cloudflare-encrypted.yaml
kubectl apply -f argocd-external-dns-cloudflare-encrypted.yaml
```

the docs for the rest of the applications can be found here:

[application docs for admin](docs/apps/admin/README.md)