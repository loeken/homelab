package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type configOption struct {
	name         string
	defaultValue string
	usage        string
	options      []string
	tags         []string
}

var options = []configOption{
	//./setup update-secret <namespace> <name> <item> <value>
	{"namespace", "vaultwarden", "the namespace to update the secret in", nil, []string{"update-secret"}},
	{"name", "vaultwarden", "the name of the secret", nil, []string{"update-secret"}},
	{"key", "vaultwarden", "the key inside the secret for example: smtp_username", nil, []string{"update-secret"}},
	{"value", "topsecure", "the value of the key for example: topsecure", nil, []string{"update-secret"}},

	{"bridge", "vmbr0", "name of the bridge", nil, []string{"install"}},
	{"cores_k3s", "10", "the amount of virtual cores to pass to the k3s vm", nil, []string{"install"}},
	{"cluster-issuer", "staging", "when using nginx ingress, select cluster issuer ( staging / prod )", nil, []string{"install"}},
	{"cloudflare_api_token", "false", "used for external-dns and to destroy dns records", nil, []string{"install", "destroy"}},
	{"disksize", "100G", "disk size + metric ( example: 100G )", nil, []string{"install"}},
	{"domain", "", "the domain you want to use", nil, []string{"install", "destroy"}},
	{"email", "", "the email used for most configs", nil, []string{"install", "destroy"}},
	{"external_ip", "1.2.3.4", "your external ipv4 ( curl -4 ifconfig.co )", nil, []string{"install"}},
	{"helo_name", "mail.example.com", "when sending email this is sent to the mail server while logging on", nil, []string{"install"}},
	{"interface", "enp3s0", "name of the primary interface", nil, []string{"install"}},
	{"ingress", "cloudflaretunnel", "which ingress to use ( nginx/cloudflaretunnel )", []string{"nginx", "cloudflaretunnel"}, []string{"install"}},
	{"kubernetes_version", "v1.26.4+k3s1", "kubernetes version", nil, []string{"install"}},
	{"local_path", "", "the folder into which new_repo will be cloned into", nil, []string{"install", "github"}},
	{"macaddr", "6E:1F:26:B6:DF:20", "mac address used for the k3s vm", nil, []string{"install"}},
	{"memory_k3s", "28672", "amount of ram in MB to assign to the VM ", nil, []string{"install"}},
	{"new_repo", "", "the name of your new repo", nil, []string{"install", "github", "check-dependencies", "destroy"}},
	{"pci_passthrough", "intel", "prepare pci passthrough", []string{"null", "intel", "amd"}, []string{"install"}},
	{"pci_device", "0000:02:00.0", "the pci address of your gpu ( lspci |grep VGA )", nil, []string{"install"}},
	{"platform", "proxmox", "server type: proxmox, minikube, baremetal", []string{"proxmox", "minikube", "baremetal"}, []string{"install"}},
	{"proxmox_node_name", "beelink-sei12", "the name of the proxmox node ( hostname )", nil, []string{"install"}},
	{"proxmox_vm_name", "k3s-beelink-01", "name of the virtual machine in proxmox", nil, []string{"install"}},
	{"root_password", "topsecure", "root password ( used for login to proxmox )", nil, []string{"install"}},
	{"smtp_from", "smtp_sender", "the email address used by vaultwarden to send emails with defaults to using the value of smtp_sender", nil, []string{"install"}},
	{"smtp_host", "mail.example.com", "the host of your email server", nil, []string{"install"}},
	{"smtp_port", "587", "the port of your email server", nil, []string{"install"}},
	{"smtp_sender", "homelab@example.com", "the email address used to send emails", nil, []string{"install"}},
	{"smtp_username", "homelab@example.com", "the username used to login to your email", nil, []string{"install"}},
	{"smtp_domain", "example.com", "the domain from which the email is sent from", nil, []string{"install"}},
	{"smtp_password", "example.com", "the password for the email", nil, []string{"install"}},
	{"ssh_password", "demotime", "ssh password", nil, []string{"install"}},
	{"ssh_private_key", "~/.ssh/id_ed25519", "location of ssh private key, id_ed25519 when generated with gh auth login", nil, []string{"install"}},
	{"ssh_public_key", "~/.ssh/id_ed25519.pub", "location of ssh public key, id_ed25519.pub when generated with gh auth login", nil, []string{"install"}},
	{"ssh_server_address", "172.16.137.36", "ip address of server for ssh connection", []string{"proxmox"}, []string{"install"}},
	{"ssh_server_gateway", "172.16.137.254", "gateway of server ( example 172.16.137.254 )", nil, []string{"install"}},
	{"ssh_server_netmask", "24", "amount of ram in MB to assign to the VM ", nil, []string{"install"}},
	{"ssh_username", "loeken", "ssh usernamer", nil, []string{"install"}},
	{"storage", "local-path", "storage type ( democratic-csi, local-path)", []string{"democratic-csi", "local-path"}, []string{"install"}},

	// app section
	{"authelia", "false", "enable argocd app authelia", nil, []string{"enable-argocd-app", "install"}},
	{"externaldns", "false", "enable argocd app external-dns", nil, []string{"enable-argocd-app", "install"}},
	{"loki", "false", "enable argocd app loki", nil, []string{"enable-argocd-app", "install"}},
	{"ha", "false", "enable argocd app home-assistant", nil, []string{"enable-argocd-app", "install"}},
	{"partition_external_shared_media_disk", "false", "will partition --shared-media-disk-device", nil, []string{"enable-argocd-app", "install", "destroy"}},
	{"shared_media_disk_size", "100Gi", "define the size of the shared media disk", nil, []string{"enable-argocd-app", "install"}},
	{"shared_media_disk_device", "sda", "give the device name of your external shared media disk", nil, []string{"enable-argocd-app", "install"}},
	{"jellyfin", "false", "enable argocd app jellyfin", nil, []string{"enable-argocd-app", "install"}},
	{"jellyseerr", "false", "enable argocd app jellyseerr", nil, []string{"enable-argocd-app", "install"}},
	{"kasten-k10", "false", "enable argocd app kasten-k10", nil, []string{"enable-argocd-app", "install"}},
	{"nextcloud", "false", "enable argocd app nextcloud", nil, []string{"enable-argocd-app", "install"}},
	{"nzbget", "false", "enable argocd app nzbget", nil, []string{"enable-argocd-app", "install"}},
	{"prowlarr", "false", "enable argocd app prowlarr", nil, []string{"enable-argocd-app", "install"}},
	{"radarr", "false", "enable argocd app radarr", nil, []string{"enable-argocd-app", "install"}},
	{"rtorrentflood", "false", "enable argocd app rtorrent-flood", nil, []string{"enable-argocd-app", "install"}},
	{"sonarr", "false", "enable argocd app sonarr", nil, []string{"enable-argocd-app", "install"}},
	{"vaultwarden", "false", "enable argocd app vaultwarden", nil, []string{"enable-argocd-app", "install"}},
	{"whoami", "true", "enable argocd app whoami", nil, []string{"enable-argocd-app", "install"}},
}

type Command struct {
	Name         string
	VersionArg   []string
	VersionRegex string
	ExpectedVer  string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "./setup",
		Short: "My command line app",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Capture executed command with flags
			commandWithFlags := cmd.CommandPath() // Get the command path
			flagSet := cmd.Flags()                // Get the set of flags associated with the command
			flagSet.VisitAll(func(flag *pflag.Flag) {
				commandWithFlags += fmt.Sprintf(" --%s=%s", flag.Name, flag.Value)
			})

			// Write the executed command with flags to .setup.sh
			writeExecutedCommand(commandWithFlags)
		},
	}
	dependencyCheckCmd := &cobra.Command{
		Use:   "check-dependencies",
		Short: "Check if all needed dependencies are installed on this system",
		Run: func(cmd *cobra.Command, args []string) {
			new_repo := viper.GetString("new_repo")
			test := strings.Split(new_repo, "/")
			if len(test) != 2 {
				color.Red("--new_repo should be in the format user/repository")
				os.Exit(0)
			}
			if new_repo == "" {
				color.Red("--new_repo should be set")
				os.Exit(0)
			}
			parts := strings.Split(new_repo, "/")
			checkDependencies(true, parts[1])

		},
	}
	updateSecrets := &cobra.Command{
		Use:   "update-secret",
		Short: "updates a value inside the secret",
		Run: func(cmd *cobra.Command, args []string) {

			namespace := viper.GetString("namespace")
			secret := viper.GetString("name")
			outputFile := fmt.Sprintf("../deploy/mysecrets/updated_secret_%s_%s.yaml", namespace, secret)

			// Execute kubectl command to fetch the secret data
			//fmt.Println("kubectl", "get", "secret", secret, "-n", namespace, "-o", "yaml")
			ExecuteCommand := exec.Command("kubectl", "get", "secret", secret, "-n", namespace, "-o", "yaml")
			out, err := ExecuteCommand.Output()
			if err != nil {
				fmt.Printf("Error executing kubectl command: %s\n", err.Error())
				os.Exit(1)
			}

			// Use yq to extract data section
			yqCmd := exec.Command("yq", ".data", "-j")
			yqCmd.Stdin = ioutil.NopCloser(strings.NewReader(string(out)))
			yqOut, err := yqCmd.Output()
			if err != nil {
				fmt.Printf("Error executing yq command: %s\n", err.Error())
				os.Exit(1)
			}

			// Parse JSON output to get data values
			var data map[string]string
			err = json.Unmarshal(yqOut, &data)
			if err != nil {
				fmt.Printf("Error decoding yq output: %s\n", err.Error())
				os.Exit(1)
			}

			// Update data values
			for key, value := range data {
				decoded, err := base64.StdEncoding.DecodeString(value)
				if err != nil {
					fmt.Printf("Error decoding data value for key %s: %s\n", key, err.Error())
					continue
				}
				data[key] = (string(decoded))
			}
			// Generate YAML output with updated data section
			yamlOut := mapToYaml(data)

			// Write YAML output to file
			err = ioutil.WriteFile(outputFile, []byte(yamlOut), 0644)
			if err != nil {
				fmt.Printf("Error writing output file: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Println("saved a new file: updated_secret_" + namespace + "_" + secret + ".yaml inside deploy/mysecrets with the new values")
			fmt.Println("to encrypt and apply the new secret use: cat updated_secret_" + namespace + "_" + secret + ".yaml | kubeseal | kubectl apply -f -")
			fmt.Println("to overwrite your encrypted version: cat updated_secret_" + namespace + "_" + secret + ".yaml | kubeseal > templates/argocd-" + secret + ".yaml")
		},
	}
	githubCmd := &cobra.Command{
		Use: "github",
		Example: `
		# to create a new repo called loeken/homelab-beelink in the /home/loeken/Projects/private folder:
		
		./setup github --new_repo loeken/homelab-beelink --local_path /home/loeken/Projects/private
		
		`,
		Short: "Create/clone/configure upstream of loeken/homelab in my github account by using the gh command line client",
		Run: func(cmd *cobra.Command, args []string) {
			// your GitHub-related
			checkDependencies(false, "")
			new_repo := viper.GetString("new_repo")
			if new_repo == "" {
				color.Red("--new_repo should be set")
				os.Exit(0)
			}
			test := strings.Split(new_repo, "/")
			if len(test) != 2 {
				color.Red("--new_repo should be in the format user/repository")
				os.Exit(0)
			}
			parts := strings.Split(new_repo, "/")
			local_path := viper.GetString("local_path")

			if local_path == "" {
				color.Red("--local_path should be set")
				os.Exit(0)
			}

			// removing former existence of files in ../tmp
			runCommand("../tmp", "rm", []string{"-rf", "homelab.git"})

			// cloning loeken/homelab into ../tmp
			runCommand("../tmp", "git", []string{"clone", "--bare", "https://github.com/loeken/homelab"})

			// creating new repo parts[1] as a private repo
			runCommand("../tmp", "gh", []string{"repo", "create", parts[1], "--private"})

			// pushing locally cloned repo to newly created repo: parts[1]
			runCommand("../tmp/homelab.git", "git", []string{"push", "--mirror", "git@github.com:" + new_repo})

			// creating local_path if it doesnt exist yet
			runCommand(".", "mkdir", []string{"-p", local_path})

			// pulling newly created repo: new_repo to local_path
			runCommand(local_path, "git", []string{"clone", "git@github.com:" + new_repo})

			// registering loeken/homelab as an upstream for the new repo
			runCommand(local_path+"/"+parts[1], "git", []string{"remote", "add", "upstream", "https://github.com/loeken/homelab.git"})

			// generating a ssh-keyring of type id_ed25519
			runCommand(local_path+"/"+parts[1]+"/tmp", "ssh-keygen", []string{"-t", "ed25519", "-f", "id_ed25519", "-C", "argocd@homelab"})

			// uploading the newly created key's public key to github as a deploy key so argocd will be able to pull from the repo
			runCommand(local_path+"/"+parts[1]+"/tmp", "gh", []string{"repo", "deploy-key", "add", "id_ed25519.pub", "--repo", new_repo})

		},
	}
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "install the stack",
		Example: `
1. in Proxmox
	# to install onto a debian 11 inside a kvm, using local-path storage and cloudflare tunnels for ingress
	
	minimal install
		./setup install --authelia false \
						--domain loeken.xyz \
						--email loeken@internetz.me \
						--external_ip 94.134.58.102 \
						--externaldns false \
						--ha false \
						--ingress cloudflaretunnel \
						--interface eno1 \
						--jellyfin false \
						--jellyseerr false \
						--kasten-k10 false \
						--loki false \
						--new_repo loeken/homelab-kubeflow \
						--nextcloud false \
						--nzbget false \
						--platform proxmox \
						--prowlarr false \
						--proxmox_node_name homeserver \
						--proxmox_vm_name kubeflow \
						--radarr false \
						--sonarr false \
						--rtorrentflood false \
						--ssh_password demotime \
						--ssh_private_key ~/.ssh/id_ed25519 \
						--ssh_public_key ~/.ssh/id_ed25519.pub \
						--ssh_server_address 172.16.137.250 \
						--ssh_server_gateway 172.16.137.254 \
						--ssh_server_netmask 24 \
						--ssh_username loeken \
						--shared_media_disk_size false \
						--shared_media_disk_device false \
						--smtp_domain internetz.me \
						--smtp_host mail.internetz.me \
						--smtp_port 587 \
						--smtp_sender homelab-beelink@internetz.me \
						--smtp_username homelab-beelink@internetz.me \
						--storage local-path \
						--vaultwarden false \
						--whoami false
	`,
		Run: func(cmd *cobra.Command, args []string) {
			checkDependencies(false, "")
			storage := viper.GetString("storage")
			platform := viper.GetString("platform")

			new_repo := viper.GetString("new_repo")
			test := strings.Split(new_repo, "/")
			if len(test) != 2 {
				color.Red("--new_repo should be in the format user/repository")
				os.Exit(0)
			}
			email := viper.GetString("email")
			domain := viper.GetString("domain")
			external_ip := viper.GetString("external_ip")
			git_parts := strings.Split(new_repo, "/")

			ingress := viper.GetString("ingress")
			clusterissuer := viper.GetString("cluster-issuer")

			smtp_host := viper.GetString("smtp_host")
			smtp_port := viper.GetString("smtp_port")
			smtp_username := viper.GetString("smtp_username")
			smtp_sender := viper.GetString("smtp_sender")
			smtp_domain := viper.GetString("smtp_domain")

			installAuthelia := viper.GetString("authelia")
			installExternalDns := viper.GetString("externaldns")
			installLoki := viper.GetString("loki")
			installHa := viper.GetString("ha")
			installNextcloud := viper.GetString("nextcloud")

			installPartitionSharedMediaDisk := viper.GetString("partition_external_shared_media_disk")
			installSharedMediaDiskSize := viper.GetString("shared_media_disk_size")
			installSharedMediaDevice := viper.GetString("shared_media_disk_device")
			installJellyfin := viper.GetString("jellyfin")
			installJellyseerr := viper.GetString("jellyseerr")
			installK10 := viper.GetString("kasten-k10")
			installRtorrentFlood := viper.GetString("rtorrentflood")
			installNzbget := viper.GetString("nzbget")
			installProwlarr := viper.GetString("prowlarr")
			installRadarr := viper.GetString("radarr")
			installSonarr := viper.GetString("sonarr")
			installVaultwarden := viper.GetString("vaultwarden")

			privateKey := viper.GetString("ssh_private_key")
			publicKey := viper.GetString("ssh_public_key")

			// pciPassthrough := viper.GetString("pci_passthrough")
			// pciDevice := viper.GetString("pci_device")

			filename1 := os.ExpandEnv(privateKey)
			_, err := os.Stat(filename1)
			if os.IsNotExist(err) {
				fmt.Printf("private key %s does not exist\n", filename1)
				os.Exit(0)
			}
			filename2 := os.ExpandEnv(publicKey)
			_, err = os.Stat(filename2)
			if os.IsNotExist(err) {
				fmt.Printf("public key %s does not exist\n", filename2)
				os.Exit(0)
			}

			u, err := user.Current()
			if err != nil {
				fmt.Println(err)
				return
			}

			if ingress != "nginx" && ingress != "cloudflaretunnel" {
				fmt.Printf("only valid inputs are nginx/cloudflaretunnel")
				os.Exit(0)
			}
			// validation section start
			checkRepo()
			color.Yellow("validating inputs")
			if email == "" {
				color.Red("you need to set an --email <youremail>")
				os.Exit(0)
			}
			if new_repo == "" {
				color.Red("you need to set a --new_repo <you/whatever>")
				os.Exit(0)
			}
			if domain == "" {
				color.Red("you need to set an --domain <internetz.me>")
				os.Exit(0)
			}
			if external_ip == "1.2.3.4" {
				color.Red("you need to setup an --external_ip <1.2.3.4> - that is not 1.2.3.4")
				os.Exit(0)
			}
			if installAuthelia == "true" {
				if smtp_host == "mail.example.com" || smtp_port == "" || smtp_username == "homelab@example.com" || smtp_sender == "homelab@example.com" {
					color.Red("when installing authelia you need to provide --smtp_host mail.example.com --smtp_port 587 --smtp_username homelab@example.com --smtp_sender homelab@example.com")
					os.Exit(0)
				}
				if ingress == "cloudflaretunnel" {
					color.Red("authelia currently only works with cloudflaretunnel untill somebody can show me how to redirect traffic for cloudflare tunnels to use authelia :)")
					os.Exit(0)
				}
			}
			if installNextcloud == "true" {
				if smtp_host == "mail.example.com" || smtp_port == "" || smtp_username == "homelab@example.com" || smtp_sender == "homelab@example.com" || smtp_domain == "example.com" {
					color.Red("when installing nextcloud you need to provide --smtp_host mail.example.com --smtp_port 587 --smtp_username homelab@example.com --smtp_sender homelab@example.com --smtp-domain example.com")
					os.Exit(0)
				}
			}
			if installJellyfin == "true" || installRtorrentFlood == "true" || installNzbget == "true" || installRadarr == "true" || installSonarr == "true" {
				if installSharedMediaDiskSize == "" || installSharedMediaDevice == "" {
					color.Red("jellyfin/rtorrent/nzbget needs a shared media disk device and size to be created: --shared_media_disk_size 100Gi")
					os.Exit(0)
				}
			}

			if platform == "minikube" {
				if installPartitionSharedMediaDisk == "true" {
					color.Red("cannot use external shared media disk with minikube")
					os.Exit(0)
				}
			}
			// validation section end

			color.Yellow("writing deploy/terraform/terraform.tfvars")
			var tfvarsContent string
			for _, opt := range options {
				tfvarsContent += fmt.Sprintf("%s = \"%s\"\n", opt.name, viper.GetString(opt.name))
			}

			filename := "../deploy/terraform/terraform.tfvars"
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				// Create the file if it doesn't exist
				if err = ioutil.WriteFile(filename, []byte(tfvarsContent), 0644); err != nil {
					color.Red(fmt.Sprintf("Error creating file %s: %s\n", filename, err))
				} else {
					color.Green(fmt.Sprintf("Created file %s\n", filename))
				}
			} else {
				// Update the file if it exists
				if err = ioutil.WriteFile(filename, []byte(tfvarsContent), 0644); err != nil {
					color.Red(fmt.Sprintf("Error updating file %s: %s\n", filename, err))
				} else {
					color.Green(fmt.Sprintf("Updated file %s\n", filename))
				}
			}

			for _, opt := range options {
				color.Green(fmt.Sprintf("%s: %s\n", opt.name, viper.GetString(opt.name)))
			}
			confirmContinue()

			if platform == "proxmox" {

				color.Green("terraform proxmox")
				runTerraformCommand("proxmox")

				if installPartitionSharedMediaDisk == "true" {
					color.Green("terraform partition external disk")
					runTerraformCommand("external-disk")
				}

				color.Green("terraform template")
				runTerraformCommand("proxmox-debian-11-template")

				color.Green("terraform k3s proxmox")
				runTerraformCommand("k3s-proxmox")

			}
			if platform == "baremetal" {
				if installPartitionSharedMediaDisk == "true" {
					color.Green("terraform partition external disk")
					runTerraformCommand("external-disk")
				}

				color.Green("terraform k3s")
				runTerraformCommand("k3s")
			}

			// prepare the bootstrap values files
			runCommand(".", "cp", []string{"../deploy/argocd/bootstrap-core-apps/values.yaml.example", "../deploy/argocd/bootstrap-core-apps/values.yaml"})

			// set values for core chart
			runCommand(".", "sed", []string{"-i", "s/loeken/" + git_parts[0] + "/", "../deploy/argocd/bootstrap-core-apps/values.yaml"})
			runCommand(".", "sed", []string{"-i", "s/homelab-example/" + git_parts[1] + "/", "../deploy/argocd/bootstrap-core-apps/values.yaml"})
			runCommand(".", "sed", []string{"-i", "s/youremail@example.com/" + email + "/", "../deploy/argocd/bootstrap-core-apps/values.yaml"})

			data, err := ioutil.ReadFile("../deploy/argocd/bootstrap-optional-apps/values.yaml.example")
			if err != nil {
				log.Fatal(err)
			}

			// Unmarshal the YAML data into a map
			config := make(map[string]interface{})
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				log.Fatal(err)
			}

			// Modify some values
			config["domain"] = domain
			config["githubUser"] = git_parts[0]
			config["githubRepo"] = git_parts[1]

			if ingress == "cloudflaretunnel" {
				config["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
				config["nginxingress"].(map[interface{}]interface{})["enabled"] = false
			}
			if ingress == "nginx" {
				config["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				config["nginxingress"].(map[interface{}]interface{})["enabled"] = true
				config["clusterIssuer"] = clusterissuer
			}
			if installAuthelia == "true" {
				autheliaConfig := config["authelia"].(map[interface{}]interface{})
				autheliaConfig["enabled"] = true
				autheliaConfig["smtp"] = make(map[interface{}]interface{})
				autheliaConfig["smtp"].(map[interface{}]interface{})["host"] = smtp_host
				autheliaConfig["smtp"].(map[interface{}]interface{})["port"] = smtp_port
				autheliaConfig["smtp"].(map[interface{}]interface{})["sender"] = smtp_sender
				autheliaConfig["smtp"].(map[interface{}]interface{})["username"] = smtp_username
				config["authelia"].(map[interface{}]interface{})["smtp"] = autheliaConfig["smtp"]
			}
			if installExternalDns == "true" {
				externaldnsConfig := config["externaldns"].(map[interface{}]interface{})
				externaldnsConfig["enabled"] = true
				config["externaldns"] = externaldnsConfig
			}
			if installLoki == "true" {
				lokiConfig := config["loki"].(map[interface{}]interface{})
				lokiConfig["enabled"] = true
				config["loki"] = lokiConfig
			}
			if installHa == "true" {
				haConfig := config["ha"].(map[interface{}]interface{})
				haConfig["enabled"] = true
				config["ha"] = haConfig
			}
			if installSharedMediaDiskSize != "false" {
				// config["jellyfin"].(map[interface{}]interface{})["enabled"] = true
				sharedMediaConfig := make(map[interface{}]interface{})
				sharedMediaConfig["sharedmedia"] = make(map[interface{}]interface{})
				sharedMediaConfig["sharedmedia"].(map[interface{}]interface{})["enabled"] = true
				sharedMediaConfig["sharedmedia"].(map[interface{}]interface{})["size"] = installSharedMediaDiskSize
				if storage == "local-path" {
					sharedMediaConfig["sharedmedia"].(map[interface{}]interface{})["storageClass"] = "nfs-client"
				} else {
					sharedMediaConfig["sharedmedia"].(map[interface{}]interface{})["storageClass"] = storage
				}
			}
			if installNextcloud == "true" {
				nextcloudConfig := config["nextcloud"].(map[interface{}]interface{})
				// nextcloudConfig := make(map[interface{}]interface{})
				nextcloudConfig["enabled"] = true
				nextcloudConfig["authelia"] = installAuthelia

				nextcloudConfig["smtp"] = make(map[interface{}]interface{})
				nextcloudConfig["smtp"].(map[interface{}]interface{})["host"] = smtp_host
				nextcloudConfig["smtp"].(map[interface{}]interface{})["port"] = smtp_port
				nextcloudConfig["smtp"].(map[interface{}]interface{})["sender"] = smtp_sender
				nextcloudConfig["smtp"].(map[interface{}]interface{})["username"] = smtp_username
				nextcloudConfig["smtp"].(map[interface{}]interface{})["domain"] = smtp_domain

				nextcloudConfig["data"] = make(map[interface{}]interface{})
				nextcloudConfig["data"].(map[interface{}]interface{})["storageClass"] = storage
				nextcloudConfig["data"].(map[interface{}]interface{})["size"] = "10Gi"

				nextcloudConfig["ingress"] = make(map[interface{}]interface{})
				nextcloudConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				nextcloudConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					nextcloudConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					nextcloudConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					nextcloudConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					nextcloudConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["nextcloud"] = nextcloudConfig
			}
			if storage == "local-path" {

				config["nfsprovisioner"].(map[interface{}]interface{})["enabled"] = true
				if platform == "baremetal" {
					config["nfsprovisioner"].(map[interface{}]interface{})["ip"] = viper.GetString("ssh_server_address")
				}

			}
			if installJellyfin == "true" {
				jellyfinConfig := config["jellyfin"].(map[interface{}]interface{})
				jellyfinConfig["enabled"] = true
				jellyfinConfig["sharedmedia"] = make(map[interface{}]interface{})
				jellyfinConfig["sharedmedia"].(map[interface{}]interface{})["size"] = installSharedMediaDiskSize
				jellyfinConfig["sharedmedia"].(map[interface{}]interface{})["enabled"] = true
				if storage == "local-path" {
					jellyfinConfig["sharedmedia"].(map[interface{}]interface{})["storageClass"] = "nfs-client"
					jellyfinConfig["sharedmedia"].(map[interface{}]interface{})["existingClaim"] = "shared-media"
				} else {
					jellyfinConfig["sharedmedia"].(map[interface{}]interface{})["storageClass"] = storage
				}

				jellyfinConfig["config"] = make(map[interface{}]interface{})
				jellyfinConfig["config"].(map[interface{}]interface{})["size"] = "1Gi"
				jellyfinConfig["config"].(map[interface{}]interface{})["enabled"] = true
				jellyfinConfig["config"].(map[interface{}]interface{})["storageClass"] = "local-path"

				jellyfinConfig["cache"] = make(map[interface{}]interface{})
				jellyfinConfig["cache"].(map[interface{}]interface{})["size"] = "10Gi"
				jellyfinConfig["cache"].(map[interface{}]interface{})["enabled"] = true
				jellyfinConfig["cache"].(map[interface{}]interface{})["storageClass"] = "local-path"

				jellyfinConfig["ingress"] = make(map[interface{}]interface{})
				jellyfinConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				jellyfinConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					jellyfinConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					jellyfinConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					jellyfinConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					jellyfinConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["jellyfin"] = jellyfinConfig
			}
			if installRtorrentFlood == "true" {
				rtorrentConfig := config["rtorrentflood"].(map[interface{}]interface{})
				rtorrentConfig["enabled"] = true
				rtorrentConfig["useAuthelia"] = false
				rtorrentConfig["linkerd"] = false
				rtorrentConfig["storageClass"] = storage

				rtorrentConfig["ingress"] = make(map[interface{}]interface{})
				rtorrentConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				rtorrentConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					rtorrentConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					rtorrentConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					rtorrentConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					rtorrentConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["rtorrentflood"] = rtorrentConfig
			}
			if installNzbget == "true" {
				nzbgetConfig := config["nzbget"].(map[interface{}]interface{})
				nzbgetConfig["enabled"] = true
				nzbgetConfig["useAuthelia"] = false
				nzbgetConfig["linkerd"] = false
				nzbgetConfig["size"] = "10Gi"
				nzbgetConfig["storageClass"] = storage
				// storage section

				// ingress section
				nzbgetConfig["ingress"] = make(map[interface{}]interface{})
				nzbgetConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				nzbgetConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					nzbgetConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					nzbgetConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					nzbgetConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					nzbgetConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["nzbget"] = nzbgetConfig
			}
			if installRadarr == "true" {
				radarrConfig := config["radarr"].(map[interface{}]interface{})
				radarrConfig["enabled"] = true
				radarrConfig["useAuthelia"] = false
				radarrConfig["linkerd"] = false
				radarrConfig["size"] = "1Gi"
				radarrConfig["storageClass"] = storage
				// storage section

				// ingress section
				radarrConfig["ingress"] = make(map[interface{}]interface{})
				radarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				radarrConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					radarrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					radarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					radarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					radarrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["radarr"] = radarrConfig
			}
			if installSonarr == "true" {
				sonarrConfig := config["sonarr"].(map[interface{}]interface{})
				sonarrConfig["enabled"] = true
				sonarrConfig["useAuthelia"] = false
				sonarrConfig["linkerd"] = false
				sonarrConfig["size"] = "1Gi"
				sonarrConfig["storageClass"] = storage
				// storage section

				// ingress section
				sonarrConfig["ingress"] = make(map[interface{}]interface{})
				sonarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				sonarrConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					sonarrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					sonarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					sonarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					sonarrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}
				config["sonarr"] = sonarrConfig
			}
			if installVaultwarden == "true" {
				vaultwardenConfig := config["vaultwarden"].(map[interface{}]interface{})
				vaultwardenConfig["enabled"] = true
				vaultwardenConfig["useAuthelia"] = false
				vaultwardenConfig["linkerd"] = false
				vaultwardenConfig["size"] = "1Gi"
				vaultwardenConfig["storageClass"] = storage
				// storage section

				// ingress section
				vaultwardenConfig["ingress"] = make(map[interface{}]interface{})
				vaultwardenConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				vaultwardenConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					vaultwardenConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					vaultwardenConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					vaultwardenConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					vaultwardenConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["vaultwarden"] = vaultwardenConfig
			}
			if installProwlarr == "true" {
				prowlarrConfig := config["prowlarr"].(map[interface{}]interface{})
				prowlarrConfig["enabled"] = true
				prowlarrConfig["useAuthelia"] = false
				prowlarrConfig["linkerd"] = false
				prowlarrConfig["size"] = "1Gi"
				prowlarrConfig["storageClass"] = storage
				// storage section

				// ingress section
				prowlarrConfig["ingress"] = make(map[interface{}]interface{})
				prowlarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				prowlarrConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					prowlarrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					prowlarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					prowlarrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					prowlarrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["prowlarr"] = prowlarrConfig
			}
			if installJellyseerr == "true" {
				jellyseerrConfig := config["jellyseerr"].(map[interface{}]interface{})
				jellyseerrConfig["enabled"] = true
				jellyseerrConfig["useAuthelia"] = false
				jellyseerrConfig["linkerd"] = false

				jellyseerrConfig["ingress"] = make(map[interface{}]interface{})
				jellyseerrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"] = make(map[interface{}]interface{})
				jellyseerrConfig["ingress"].(map[interface{}]interface{})["nginx"] = make(map[interface{}]interface{})
				if ingress == "nginx" {
					jellyseerrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = true
					jellyseerrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = false
				}
				if ingress == "cloudflaretunnel" {
					jellyseerrConfig["ingress"].(map[interface{}]interface{})["cloudflaretunnel"].(map[interface{}]interface{})["enabled"] = true
					jellyseerrConfig["ingress"].(map[interface{}]interface{})["nginx"].(map[interface{}]interface{})["enabled"] = false
				}

				config["jellyseerr"] = jellyseerrConfig
			}
			if installK10 == "true" {
				config["k10"].(map[interface{}]interface{})["enabled"] = true
				config["k10"].(map[interface{}]interface{})["storageClass"] = storage
			}
			// Marshal the modified map into YAML
			modifiedData, err := yaml.Marshal(&config)
			if err != nil {
				log.Fatal(err)
			}

			// Save the modified YAML to a new file
			err = ioutil.WriteFile("../deploy/argocd/bootstrap-optional-apps/values.yaml", modifiedData, 0644)
			if err != nil {
				log.Fatal(err)
			}

			runCommand(".", "git", []string{"add", "../deploy/argocd/bootstrap-core-apps/values.yaml"})
			runCommand(".", "git", []string{"add", "../deploy/argocd/bootstrap-optional-apps/values.yaml"})
			runCommand(".", "git", []string{"commit", "-m", "initial commit of values.yaml for argocd bootstrap apps"})
			runCommand(".", "git", []string{"pull"})
			runCommand(".", "git", []string{"push"})
			path := os.ExpandEnv("~/.kube/config")
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				fmt.Println("file does not exist")
			} else {
				runCommand("~", "cp", []string{"-n", ".kube/config", ".kube/config.bak"})
			}
			runCommand("../tmp", "sed", []string{"-i", "s/default/" + git_parts[1] + "/g", "kubeconfig"})

			fmt.Println(u.HomeDir + "/.kube/config.tmp")

			if os.IsNotExist(err) {
				runCommand("..", "mkdir", []string{"-p", u.HomeDir + "/.kube/"})
				runCommand("../tmp", "cp", []string{"kubeconfig", u.HomeDir + "/.kube/config"})
			} else {
				MergeConfigs("../tmp/kubeconfig", u.HomeDir+"/.kube/config", u.HomeDir+"/.kube/config")
			}
			runCommand("~", "kubectl", []string{"config", "use-context", git_parts[1]})

			// loading secret before argocd so it can sync straight away

			runTerraformCommand("sealed-secrets")
			loadSecretFromTemplate("argocd", "repo")
			fmt.Println("terraform bootstrap argocd")

			runTerraformCommand("bootstrap-argocd")
			waitForPodReady("argocd", "app.kubernetes.io/component=server")
			color.Green("---")
			color.Green("argocd is now up you can follow the rest of the installation")
			color.Green("---")
			color.Green("	kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath=\"{.data.password}\" | base64 -d && echo")
			color.Green("	kubectl port-forward service/argocd-server -n argocd 8080:443")
			color.Green("---")
			color.Green("starting installation of additional apps")
			color.Green("supports multiline input, if single line input press enter twice")

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Once you logged into argocd, Press any key to continue...")
			_, confirmationErr := reader.ReadByte()

			if confirmationErr != nil {
				fmt.Println("Error reading input:", confirmationErr)
			} else {
				fmt.Println("Continuing...")
			}

			if ingress == "cloudflaretunnel" {
				fmt.Println("you selected cloudflare ingress please login")
				runCommand("../tmp", "cloudflared", []string{"login"})
				runCommand("../tmp", "kubectl", []string{"create", "namespace", "cloudflaretunnel"})
				runCommand("../tmp", "cloudflared", []string{"tunnel", "create", "homelab-tunnel_" + new_repo})
				runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "whoami." + domain})
				// runCommand(".", "kubectl", []string{"wait", "--for=condition=ready", "pod", "-n", "kube-system", "-l", "app.kubernetes.io/instance=sealed-secrets-controller", "--timeout=300s"})

				cfTunnelId := cloudflaretunnel("homelab-tunnel_" + new_repo)
				cloudflaresecret(cfTunnelId, *u)

				waitForPodReady("cloudflaretunnel", "app=cloudflared")
			}

			// wave 12
			if installExternalDns == "true" {
				color.Green("you need to create an api token @ https://dash.cloudflare.com/profile/api-tokens")
				color.Green("The token should be granted Zone Read, DNS Edit privileges, and access to All zones. Example config:")
				color.Green("")
				color.Green("Section Permissions:")
				color.Green("Zone 		Zone 			Read")
				color.Green("Zone 		DNS 			Edit")
				color.Green("")
				color.Green("Zone Resource:")
				color.Green("Include	Specific Zone	loeken.xyz")
				color.Green("")
				color.Green("the helm chart is also setup for external dns to only use the --domain you specified")
				color.Blue("\033[1m input settings for external-dns:\033[0m")
				loadSecretFromTemplate("external-dns", "externaldns")
				waitForPodReady("external-dns", "app.kubernetes.io/instance=external-dns")
			}

			// wave 13
			if installAuthelia == "true" {
				// Generate password hash
				color.Blue("\033[1m input settings for authelia:\033[0m")
				_, err := os.Stat("../tmp/authelia_users_database.yml")
				if err != nil {
					fmt.Println("no errors")
					loadSecretFromTemplate("authelia", "authelia")
					fmt.Println("loaded Secrets")
					generateAutheliaUsersDatabase()
					fmt.Println("generated authelia")
					color.Green("waiting for authelia to be up, to upload /config/users_database.yml")
				}

				waitForPodReady("authelia", "app.kubernetes.io/instance=authelia")
				runCommand("../tmp", "kubectl", []string{"cp", "authelia_users_database.yml", "authelia/authelia-0:/config/users_database.yml"})

				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "auth." + domain})
				}
			}
			// wave 4
			if installVaultwarden == "true" {
				color.Blue("\033[1m input settings for vaultwarden:\033[0m")
				loadSecretFromTemplate("vaultwarden", "vaultwarden")
				waitForPodReady("vaultwarden", "app.kubernetes.io/name=vaultwarden")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "vaultwarden." + domain})
				}
			}
			// wave 14
			if installNextcloud == "true" {
				color.Blue("\033[1m input settings for nextcloud:\033[0m")
				loadSecretFromTemplate("nextcloud", "nextcloud")
				waitForPodReady("nextcloud", "app.kubernetes.io/name=nextcloud")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "nextcloud." + domain})
				}
			}

			// wave none
			if installSharedMediaDiskSize != "false" {
				color.Blue("\033[1m input settings for media disk:\033[0m")
				runCommand(".", "kubectl", []string{"create", "namespace", "media"})
				err := createPVC("shared-media", "media", "nfs-client", installSharedMediaDiskSize)
				if err != nil {
					fmt.Println("error creating shared media disk:", err)
				} else {
					color.Green("created shared media pvc")
				}
			}

			// wave 20
			if installJellyfin == "true" {
				color.Blue("\033[1m input settings for jellyfin:\033[0m")
				waitForPodReady("media", "app.kubernetes.io/instance=jellyfin")
				podName, err := runCommand(".", "kubectl", []string{"get", "pods", "-n", "media", "-l", "app.kubernetes.io/instance=jellyfin", "-o", "jsonpath='{.items[0].metadata.name}'"})
				if err != nil {
					// handle error
					fmt.Println("error: ", err)
					os.Exit(3)
				}
				createFolderJellyfin(podName, "/media/tv")
				createFolderJellyfin(podName, "/media/movie")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "jellyfin." + domain})
				}
			}
			// wave 21
			if installJellyseerr == "true" {
				color.Blue("\033[1m input settings for jellyseerr:\033[0m")
				waitForPodReady("media", "app.kubernetes.io/instance=jellyseerr")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "jellyseerr." + domain})
				}
			}

			// wave 22
			if installRtorrentFlood == "true" {
				color.Blue("\033[1m input settings for rtorrent:\033[0m")
				waitForPodReady("media", "app.kubernetes.io/instance=rtorrent-flood")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "rtorrent." + domain})
				}
			}

			// wave 22
			if installNzbget == "true" {
				color.Blue("\033[1m input settings for nzbget:\033[0m")
				loadSecretFromTemplate("media", "nzbget")
				waitForPodReady("media", "app.kubernetes.io/instance=nzbget")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "rtorrent." + domain})
				}
			}

			// wave 23
			if installProwlarr == "true" {
				color.Blue("\033[1m input settings for prowlarr:\033[0m")
				waitForPodReady("media", "app.kubernetes.io/instance=prowlarr")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "prowlarr." + domain})
				}
			}
			// wave 23
			if installRadarr == "true" {
				color.Blue("\033[1m input settings for radarr:\033[0m")
				waitForPodReady("media", "app.kubernetes.io/instance=radarr")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "radarr." + domain})
				}
			}
			// wave 23
			if installSonarr == "true" {
				color.Blue("\033[1m input settings for sonarr:\033[0m")
				waitForPodReady("media", "app.kubernetes.io/instance=sonarr")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "sonarr." + domain})
				}
			}

			// wave 30
			if installLoki == "true" {
				color.Blue("\033[1m input settings for loki:\033[0m")
				loadSecretFromTemplate("loki", "loki")
				waitForPodReady("loki", "statefulset.kubernetes.io/pod-name=loki-stack-charts-0")
				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "grafana." + domain})
				}
			}
			// wave 30
			if installHa == "true" {
				color.Blue("\033[1m input settings for home-assistant:\033[0m")
				loadSecretFromTemplate("home-assistant", "home-assistant")

				color.Blue("we ll now wait for home assistant to be up this can take a bit of time - expect errors to be displayed untill its up")

				waitForPodReady("home-assistant", "statefulset.kubernetes.io/pod-name=home-assistant-0")
				time.Sleep(5 * time.Second)
				runCommand("../tmp", "kubectl", []string{"cp", "../deploy/helpers/ha_configuration.yml", "home-assistant/home-assistant-0:/config/configuration.yaml"})
				runCommand("../tmp", "echo", []string{"kubectl", "cp", "../deploy/helpers/ha_configuration.yaml", "home-assistant/home-assistant-0:/config/configuration.yaml"})
				runCommand(".", "kubectl", []string{"rollout", "restart", "statefulset", "home-assistant", "-n", "home-assistant"})

				if ingress == "cloudflaretunnel" {
					runCommand("../tmp", "cloudflared", []string{"tunnel", "route", "dns", "homelab-tunnel_" + new_repo, "ha." + domain})
				}
			}

			color.Green("installation finished!")
		},
	}

	destroyCmd := &cobra.Command{
		Use:   "destroy",
		Short: "destroy the stack ( DANGER DANGER! :) )",
		Run: func(cmd *cobra.Command, args []string) {
			mycmd := exec.Command("sh", "-c", "cat ../.git/config | grep url |grep -v loeken/homelab.git| cut -d' ' -f 3")
			cloudflare_api_token := viper.GetString("cloudflare_api_token")
			domain := viper.GetString("domain")
			var out bytes.Buffer
			mycmd.Stdout = &out
			err := mycmd.Run()
			if err != nil {
				fmt.Println("Error running command:", err)
				return
			}
			new_repo := strings.TrimSpace(out.String())
			parts := strings.Split(new_repo, "/")
			installPartitionSharedMediaDisk := viper.GetString("partition_external_shared_media_disk")
			checkDependencies(true, parts[1])

			color.Red("this will run terraform destroys on all terraform folders")
			confirmContinue()
			fmt.Println("dont forget to wipefs disk: ", installPartitionSharedMediaDisk)
			runCommand("../deploy/terraform/bootstrap-argocd", "terraform", []string{"init"})
			runCommand("../deploy/terraform/bootstrap-argocd", "terraform", []string{"destroy", "--auto-approve", "-var-file=../terraform.tfvars"})
			runCommand("../deploy/terraform/k3s-proxmox", "terraform", []string{"init"})
			runCommand("../deploy/terraform/k3s-proxmox", "terraform", []string{"destroy", "--auto-approve", "-var-file=../terraform.tfvars"})
			runCommand("../deploy/terraform/k3s", "terraform", []string{"init"})
			runCommand("../deploy/terraform/k3s", "terraform", []string{"destroy", "--auto-approve", "-var-file=../terraform.tfvars"})
			runCommand("../deploy/terraform/proxmox-debian-11-template", "terraform", []string{"init"})
			runCommand("../deploy/terraform/proxmox-debian-11-template", "terraform", []string{"destroy", "--auto-approve", "-var-file=../terraform.tfvars"})
			runCommand("../deploy/terraform/proxmox", "terraform", []string{"init"})
			runCommand("../deploy/terraform/proxmox", "terraform", []string{"destroy", "--auto-approve", "-var-file=../terraform.tfvars"})
			runCommand("../deploy/terraform/external-disk", "terraform", []string{"init"})
			runCommand("../deploy/terraform/external-disk", "terraform", []string{"destroy", "--auto-approve", "-var-file=../terraform.tfvars"})

			runCommand("../tmp", "cloudflared", []string{"tunnel", "cleanup", "homelab-tunnel_" + new_repo})
			runCommand("../tmp", "cloudflared", []string{"tunnel", "delete", "homelab-tunnel_" + new_repo})
			runCommand("../tmp", "rm", []string{"-rf", "authelia_users_database.yml"})
			runCommand("../deploy/mysecrets/templates", "find", []string{".", "-type", "f", "-delete"})
			runCommand("../deploy/terraform", "find", []string{".", "-name", "*terraform.tfstate*", "-type", "f", "-delete"})
			if cloudflare_api_token != "false" && domain != "" {
				deleteAllDNSRecords(cloudflare_api_token, domain)
			}
		},
	}

	// Add subcommands to root command
	rootCmd.AddCommand(githubCmd)
	rootCmd.AddCommand(updateSecrets)
	//rootCmd.AddCommand(enableArgocdApp)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(destroyCmd)
	rootCmd.AddCommand(dependencyCheckCmd)

	// Add options to subcommand based on command line argument

	for _, opt := range options {

		for _, tag := range opt.tags {

			//fmt.Println(opt.name, " ", tag)
			if tag == os.Args[1:][0] {
				var cmdToAddOptions *cobra.Command

				switch os.Args[1:][0] {
				case "github":
					cmdToAddOptions = githubCmd
				case "install":
					cmdToAddOptions = installCmd
				case "check-dependencies":
					cmdToAddOptions = dependencyCheckCmd
				// case "enable-argocd-app":
				// 	cmdToAddOptions = enableArgocdApp
				case "destroy":
					cmdToAddOptions = destroyCmd
				case "update-secret":
					cmdToAddOptions = updateSecrets
				}
				cmdToAddOptions.Flags().String(
					opt.name,
					opt.defaultValue,
					opt.usage,
				)
				cmdToAddOptions.RegisterFlagCompletionFunc(opt.name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
					return opt.tags, cobra.ShellCompDirectiveNoFileComp
				})

				viper.BindPFlag(opt.name, cmdToAddOptions.Flags().Lookup(opt.name))
			}
		}
	}
	// Read in config file if it exists
	viper.SetConfigName(".setup")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Handle errors reading the config file
		} else {
			// Create the config file if it doesn't exist
			configPath := filepath.Join(".", ".setup.yaml")
			f, err := os.Create(configPath)
			if err != nil {
				// Handle errors creating the file
				fmt.Println("error creating config file")
			}
			f.Close()
		}
	}

	// Write config file with flags
	err = viper.WriteConfig()
	if err != nil {
		// Handle errors writing the config file
		fmt.Println("error writing config: ", err)
	}
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func runTerraformCommand(folder string) {
	runCommand("../deploy/terraform/"+folder, "terraform", []string{"init"})
	runCommand("../deploy/terraform/"+folder, "terraform", []string{"apply", "-auto-approve", "-var-file=../terraform.tfvars"})
}
func runCommand(folder string, command string, args []string) (string, error) {
	// Get the absolute path of the folder
	var absFolder string
	if folder[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("error getting absolute path for folder %s: %w", folder, err)
		}
		folder = filepath.Join(u.HomeDir, folder[1:])

		absFolder, err = filepath.Abs(folder)
		if err != nil {
			return "", fmt.Errorf("error getting absolute path for folder %s: %w", folder, err)
		}

	} else {
		var err error
		absFolder, err = filepath.Abs(folder)
		if err != nil {
			return "", fmt.Errorf("error getting absolute path for folder %s: %w", folder, err)
		}

	}

	// Create the command

	cmd := exec.Command(command, args...)
	fmt.Println("---")
	color.Blue("Executing command in: " + absFolder)
	color.Yellow("Raw command: %v\n", cmd.String())
	cmd.Dir = absFolder

	// Capture the stdout and stderr of the command
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting command %s: %w", command, err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command step 2 %s: %s\n", cmd.Args, err)
		// if command == "terraform" {
		// 	if cmd.Args[1] != "destroy" {
		// 		fmt.Println("terraform failed executing halting operations")
		// 		os.Exit(3)
		// 	}
		// }
	}

	return stdout.String(), nil
}
func runCommandWithRetries(folder string, command string, args []string, maxRetries int, retryTimeout time.Duration) (string, error) {
	var out string
	var err error

	for i := 0; i <= maxRetries; i++ {
		out, err = runCommand(folder, command, args)
		if err == nil {
			return out, nil
		}
		fmt.Printf("Error executing command: %v\nRetrying in %v...\n", err, retryTimeout)
		time.Sleep(retryTimeout)
	}
	return out, fmt.Errorf("command execution failed after %d attempts: %v", maxRetries, err)
}
func checkRepo() {

	cmd := exec.Command("sh", "-c", "cat ../.git/config | grep url |grep -v loeken/homelab.git| cut -d' ' -f 3")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
	repo := strings.TrimSpace(out.String())
	if strings.HasPrefix(repo, "git@github.com:") {
		repo = strings.TrimPrefix(repo, "git@github.com:")
	} else {
		repo = strings.TrimPrefix(repo, "https://github.com/")
	}
	repo = strings.TrimSuffix(repo, ".git")

	if repo == "loeken/homelab" {
		fmt.Println("this is the repo loeken/homelab, you should run this tool from a cloned repo of loeken/homelab which uses loeken/homelab as an upstream")
		fmt.Println("./setup github -h")
		os.Exit(0)
	}
}
func checkDependencies(verbose bool, repoName string) {
	// Define the commands to check
	color.Blue("there is a helper script for ubuntu ./dependencies_ubuntu2204.sh to install dependencies")
	commands := []string{"gh", "cloudflared", "git", "terraform", "kubectl", "sshpass", "kubeseal", "k3sup", "docker"}
	// Loop through the commands and check if they're available
	for _, cmd := range commands {
		// Check if the command is available
		_, err := exec.LookPath(cmd)
		if err != nil {
			color.Red("%s is not available", cmd)
			continue
		}

		if verbose {
			color.Green("%s is available", cmd)
		}
	}

	cmd := exec.Command("gh", "auth", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("you first need to login to github via the github cli:\n  gh auth login")
		fmt.Println("when you login chose github.com, and select a private key ( no passphrase required )")
		os.Exit(0)
	}
	if strings.Contains(string(output), "Logged in to github.com") {
		color.Green("You are logged in with the GitHub CLI (gh)")
	} else {
		color.Red("You are not logged in with the GitHub CLI (gh), please login gh auth login")
		os.Exit(0)
	}
	if checkGitHubSSHKey() {
		color.Green("we can write to github via ssh")
	} else {
		color.Red("you need an ssh key ( ssh-keygen ) and the .pub key needs to be added to github so we can auth and push changes")
		os.Exit(0)
	}

	checkGitAccount()

	if repoName != "" {
		if isContextActive(repoName) {
			color.Green("The current Kubernetes context is %s\n", repoName)
		} else {
			color.Red("The current Kubernetes context is not %s\n", repoName)
		}
	}

	canRun, err := canRunDocker()
	if canRun {
		color.Green("Docker is installed, and the user can run Docker containers.")
	} else {
		color.Red("Docker may not be installed, or the user may not have the necessary permissions.")
		if err != nil {
			color.Red("Error:", err)
		}
	}
}
func checkGitHubSSHKey() bool {
	cmd := exec.Command("ssh", "-T", "git@github.com")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "You've successfully authenticated") {
			return true
		} else {
			return false
		}
	}
	if strings.Contains(string(output), "You've successfully authenticated") {
		return true
	} else {
		return false
	}
}

func isContextActive(contextName string) bool {
	// Get the name of the current Kubernetes context
	currentContextBytes, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return false
	}
	currentContext := strings.TrimSpace(string(currentContextBytes))

	// Check if the current context matches the desired context
	return currentContext == contextName
}
func confirmContinue() {
	fmt.Print("Do you want to continue? [Y/n] ")
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println(err)
		return
	}

	response = strings.TrimSpace(response)
	if response == "Y" || response == "y" || response == "" || response == "\\n" {
		// Continue with the program
		fmt.Println("let's go! ...")
	} else {
		// Stop the program
		fmt.Println("aborted ...")
		os.Exit(0)
	}
}

func MergeConfigs(cfg1 string, cfg2 string, configOut string) {
	// Load the first kubeconfig file.
	kubeconfig1, err := ioutil.ReadFile(cfg1)
	if err != nil {
		panic(err)
	}
	config1, err := clientcmd.Load(kubeconfig1)
	if err != nil {
		panic(err)
	}

	// Load the second kubeconfig file.
	kubeconfig2, err := ioutil.ReadFile(cfg2)
	if err != nil {
		panic(err)
	}
	config2, err := clientcmd.Load(kubeconfig2)
	if err != nil {
		panic(err)
	}

	// Merge the two kubeconfig files.
	mergedConfig := api.NewConfig()
	mergedConfig.AuthInfos = make(map[string]*clientcmdapi.AuthInfo)
	mergedConfig.Clusters = make(map[string]*clientcmdapi.Cluster)
	mergedConfig.Contexts = make(map[string]*clientcmdapi.Context)
	mergedConfig.Preferences = config1.Preferences
	mergedConfig.APIVersion = config1.APIVersion
	mergedConfig.CurrentContext = config1.CurrentContext

	mergedConfig.Clusters = mergeClusters(config1.Clusters, config2.Clusters)
	mergedConfig.AuthInfos = mergeAuthInfos(config1.AuthInfos, config2.AuthInfos)
	mergedConfig.Contexts = mergeContexts(config1.Contexts, config2.Contexts)

	// Write the merged kubeconfig file to disk.
	mergedKubeconfig, err := clientcmd.Write(*mergedConfig)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configOut, mergedKubeconfig, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("Merged kubeconfig file written to /path/to/merged-kubeconfig")
}

// Helper function to merge clusters.
func mergeClusters(c1 map[string]*clientcmdapi.Cluster, c2 map[string]*clientcmdapi.Cluster) map[string]*clientcmdapi.Cluster {
	result := make(map[string]*clientcmdapi.Cluster)
	for key, value := range c1 {
		result[key] = value.DeepCopy()
	}
	for key, value := range c2 {
		if _, ok := result[key]; !ok {
			result[key] = value.DeepCopy()
		}
	}
	return result
}

// Helper function to merge authentication information.
func mergeAuthInfos(a1 map[string]*clientcmdapi.AuthInfo, a2 map[string]*clientcmdapi.AuthInfo) map[string]*clientcmdapi.AuthInfo {
	result := make(map[string]*clientcmdapi.AuthInfo)
	for key, value := range a1 {
		result[key] = value.DeepCopy()
	}
	for key, value := range a2 {
		if _, ok := result[key]; !ok {
			result[key] = value.DeepCopy()
		}
	}
	return result
}

// Helper function to merge contexts.
func mergeContexts(c1 map[string]*clientcmdapi.Context, c2 map[string]*clientcmdapi.Context) map[string]*clientcmdapi.Context {
	result := make(map[string]*clientcmdapi.Context)
	for key, value := range c1 {
		result[key] = value.DeepCopy()
	}
	for key, value := range c2 {
		if _, ok := result[key]; !ok {
			result[key] = value.DeepCopy()
		}
	}
	return result
}
func cloudflaretunnel(repo_name string) (tunnelId string) {
	// Run `cloudflared tunnel list` and get its output
	cmd1 := exec.Command("cloudflared", "tunnel", "list")
	output1, err := cmd1.Output()
	if err != nil {
		panic(err)
	}

	// Find the line that contains repo_name and extract its ID (the first field)
	var id string
	for _, line := range strings.Split(string(output1), "\n") {
		if strings.Contains(line, repo_name) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				id = fields[0]
				break
			}
		}
	}

	if id == "" {
		fmt.Println("No tunnel found for " + repo_name)
		return
	}
	return id
}

type Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Type       string `yaml:"type"`
	Metadata   struct {
		Namespace string `yaml:"namespace"`
		Name      string `yaml:"name"`
	} `yaml:"metadata"`
	StringData map[string]string `yaml:"stringData"`
}

func loadSecretFromTemplate(namespace string, application string) {

	_, err := os.Stat("../deploy/mysecrets/templates/argocd-" + namespace + "-encrypted.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			// Handle file does not exist error
			color.Green("no existing sealed secret found, creating new one")
		} else {
			// Handle other errors
			color.Green("unexpected error in loadSecretFromTemplate ", err)
		}
	} else {
		// File exists
		color.Red("secret already exists for " + namespace + " if you want to create a new delete deploy/mysecrets/templates/argocd-" + namespace + "-encrypted.yaml")
		return
	}

	// Read the secret YAML template file
	data, err := ioutil.ReadFile("../deploy/mysecrets/argocd-" + application + ".yaml.example")

	if err != nil {
		fmt.Printf("Error loading secrets file: %v\n", err)
		return
	}

	// Parse YAML file
	var secrets map[string]interface{}
	if err := yaml.Unmarshal(data, &secrets); err != nil {
		fmt.Printf("Error parsing secrets file: %v\n", err)
		return
	}

	// Loop through secrets
	scanner := bufio.NewScanner(os.Stdin)
	for key, value := range secrets["stringData"].(map[interface{}]interface{}) {
		// Convert key and value to the appropriate types
		strKey := key.(string)
		strValue := value.(string)

		if strKey == "DOMAIN" {
			if !strings.HasPrefix(viper.GetString("domain"), "https://") {
				color.Green("found " + strKey + " value in arguments, reusing that as default, adding https:// prefix")
				strValue = "https://" + viper.GetString(strKey)
			} else {
				color.Green("found " + strKey + " value in arguments, reusing that as default")
				strValue = viper.GetString(strKey)
			}
		} else if strKey == "SMTP_FROM" && viper.GetString("SMTP_FROM") == "smtp_sender" {
			color.Green("found " + strKey + " value in arguments, reusing the value of --smtp_sender as default")
			strValue = viper.GetString("smtp_sender")
		} else if strKey == "URL" {
			cmd := exec.Command("bash", "-c", "cat ../.git/config|grep url|grep git@| cut -d' ' -f 3")

			// Run the command and capture its output
			output, err := cmd.Output()
			if err != nil {
				fmt.Println(err)
				return
			}

			// Convert the output to a string and remove any trailing newline characters
			url := strings.TrimSpace(string(output))
			strValue = url
		} else if strKey == "sshPrivateKey" {
			// Read the contents of the file "../tmp/id_ed25519"
			privateKeyBytes, err := ioutil.ReadFile("../tmp/id_ed25519")
			if err != nil {
				fmt.Printf("Error reading sshPrivateKey from file: %v\n", err)
				return
			}
			secrets["stringData"].(map[interface{}]interface{})[key] = string(privateKeyBytes)
			continue
		} else if strings.HasPrefix(strValue, "generate|") {
			parts := strings.Split(strValue, "generate|")
			if len(parts) > 1 {
				fmt.Println("Part after 'generate|':", parts[1])
				length, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("error converting length")
				}
				strValue = generatePassword(length)
			}
		} else if viper.GetString(strKey) != "" {
			color.Green("found " + strKey + " value in arguments, reusing that as default")
			strValue = viper.GetString(strKey)
		}
		// Print secret name and default value
		fmt.Printf("%s (%s):\n", strKey, strValue)
		var input string

		// Read input from user
		var inputBuffer bytes.Buffer
		for {
			if scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					break
				}
				inputBuffer.WriteString(line + "\n")
			} else {
				fmt.Printf("Error reading input: %v\n", scanner.Err())
				return
			}
		}

		// Use default value if input is empty
		input = strings.TrimSpace(inputBuffer.String())
		if input == "" {
			input = strValue
		}

		// Update secret value
		secrets["stringData"].(map[interface{}]interface{})[key] = input
	}

	// Convert secrets to YAML and save to file
	output, err := yaml.Marshal(secrets)
	if err != nil {
		fmt.Printf("Error converting secrets to YAML: %v\n", err)
		return
	}
	err = ioutil.WriteFile("../deploy/mysecrets/argocd-"+namespace+".yaml", output, 0644)
	if err != nil {
		fmt.Printf("Error saving secrets to file: %v\n", err)
		return
	}

	sealedSecret, err := sealSecret(string(output))
	if err != nil {
		fmt.Printf("Failed to seal secret1: %v\n", err)
		return
	}
	createSecret(namespace, string(sealedSecret))
}
func generatePassword(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	password := make([]byte, n)
	for i := range password {
		password[i] = letters[rand.Intn(len(letters))]
	}
	return string(password)
}

func createSecret(namespace string, yamlString string) ([]byte, error) {
	runCommand("../tmp", "kubectl", []string{"create", "namespace", namespace})
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(yamlString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("../deploy/mysecrets/templates/argocd-"+namespace+"-encrypted.yaml", []byte(yamlString), 0644)
	if err != nil {
		panic(err)
	}
	return output, nil
}

func sealSecret(secretYaml string) ([]byte, error) {
	// Use kubeseal to encrypt the secret data
	kubesealCmd := exec.Command("kubeseal", "--format=yaml")
	kubesealCmd.Stdin = strings.NewReader(secretYaml)

	var sealed bytes.Buffer
	kubesealCmd.Stdout = &sealed

	if err := kubesealCmd.Run(); err != nil {
		return nil, err
	}

	return sealed.Bytes(), nil
}

func cloudflaresecret(cfTunnelId string, u user.User) {
	// Define the command to create the secret
	createSecretCmd := exec.Command("kubectl", "create", "secret", "generic", "tunnel-credentials", "--from-file=credentials.json="+u.HomeDir+"/.cloudflared/"+cfTunnelId+".json", "-n", "cloudflaretunnel", "--dry-run=client", "-o", "yaml")

	// Run the command and capture its output
	var out bytes.Buffer
	createSecretCmd.Stdout = &out
	if err := createSecretCmd.Run(); err != nil {
		fmt.Printf("Failed to create secret: %v\n", err)
		return
	}

	// Get the YAML definition for the secret from the command's output
	secretYaml := out.String()

	// Seal the secret
	sealedSecret, err := sealSecret(string(secretYaml))
	if err != nil {
		fmt.Printf("Failed to seal secret2: %v\n", err)
		return
	}
	createSecret("cloudflaretunnel", string(sealedSecret))

	fmt.Println("Secret created and applied successfully")
}
func hashAutheliaPassword(password string) string {
	cmd := exec.Command("docker", "run", "--rm", "authelia/authelia:latest", "authelia", "hash-password", password)
	stdout, err := cmd.Output()
	if err != nil {
		// Handle error running docker command
		color.Red("error running docker command: ", err)
	}

	output := strings.TrimSpace(string(stdout))
	output = strings.TrimPrefix(output, "Digest: ")
	index := strings.Index(output, ":")
	if index != -1 {
		output = output[index+1:]
	}
	return output
}
func generateAutheliaUsersDatabase() {
	color.Green("---")
	color.Green("launching authelia/authelia via docker to generate a users file:")
	color.Green("---")

	// Prompt user to enter user details
	fmt.Print("Enter username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter display name: ")
	scanner.Scan()
	displayName := scanner.Text()

	fmt.Print("Enter email address: ")
	scanner.Scan()
	email := scanner.Text()

	fmt.Print("Enter groups (comma separated): ")
	scanner.Scan()
	groups := scanner.Text()
	groupList := strings.Split(groups, ",")

	fmt.Print("Enter password: ")
	scanner.Scan()
	password := scanner.Text()

	// Generate password hash
	hash := hashAutheliaPassword(password)

	// Create user data map
	userData := make(map[string]interface{})
	userData["disabled"] = false
	userData["displayname"] = displayName
	userData["password"] = hash
	userData["email"] = email
	userData["groups"] = groupList

	// Create users map
	users := make(map[string]interface{})
	users[username] = userData

	// Create YAML data
	data, err := yaml.Marshal(map[string]interface{}{
		"users": users,
	})
	if err != nil {
		// Handle error marshaling YAML
		color.Red("error marsheling authelia users_database.yml: ", err)
	}

	// Write YAML data to file
	err = ioutil.WriteFile("../tmp/authelia_users_database.yml", data, 0644)
	if err != nil {
		// Handle error writing file
		color.Red("error writing file for authelia users_database.yml: ", err)
	}
}
func createPVC(pvcName string, namespace string, storageClass string, storageSize string) error {
	args := []string{"apply", "-f", "-"}
	yamlTemplate := fmt.Sprintf(`
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: %s
  namespace: %s
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: %s
  storageClassName: %s`, pvcName, namespace, storageSize, storageClass)

	fmt.Println(yamlTemplate)
	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = strings.NewReader(yamlTemplate)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating PVC %s: %s: %w", pvcName, string(output), err)
	}

	return nil
}
func createFolderJellyfin(podName string, folderName string) {
	cmdArgs := []string{"exec", "-n", "media", strings.ReplaceAll(podName, "'", ""), "--", "mkdir", "-p", folderName}
	_, err := runCommand(".", "kubectl", cmdArgs)
	if err != nil {
		// handle error
		fmt.Println("error: ", err)
		os.Exit(3)
	}
}

/*
	func waitForPodReady(namespace string, podName string) {
		// First, check if the namespace has been created
		nsCreated, nsErr := runCommandWithRetries(".", "kubectl", []string{"get", "namespace", namespace}, 10, 5*time.Second)
		if nsErr != nil {
			fmt.Printf("Error checking for namespace existence: %v\n", nsErr)
			return
		}
		fmt.Printf("Namespace is created: %s\n", nsCreated)

		// Now wait for the pod to be ready
		out, err := runCommandWithRetries(".", "kubectl", []string{"wait", "--for=condition=ready", "pod", "-n", namespace, "-l", "app.kubernetes.io/instance=" + podName, "--timeout=300s"}, 10, 5*time.Second)
		if err != nil {
			fmt.Printf("Error waiting for pod to be ready: %v\n", err)
		} else {
			fmt.Printf("Pod is ready: %s\n", out)
		}
	}
*/
func waitForPodReady(namespace string, searchTag string) {
	// Wait for the namespace and pod to be ready
	err := waitWithRetries(10, 10*time.Second, func() (bool, error) {
		_, nsErr := runCommand(".", "kubectl", []string{"get", "namespace", namespace})
		if nsErr == nil {
			out, err := runCommand(".", "kubectl", []string{"get", "pods", "-n", namespace, "-l", searchTag})
			if err == nil && strings.Contains(out, "Running") {
				fmt.Printf("Pod is ready: %s\n", out)
				return true, nil
			} else {
				fmt.Printf("Error getting pod status: %v\n", err)
			}
		} else {
			fmt.Printf("Error checking for namespace existence: %v\n", nsErr)
		}
		return false, nil
	})
	if err != nil {
		fmt.Printf("Error waiting for pod to be ready: %v\n", err)
	}
}

func waitWithRetries(maxRetries int, retryTimeout time.Duration, conditionFunc func() (bool, error)) error {
	for i := 0; i <= maxRetries; i++ {
		success, err := conditionFunc()
		if success {
			return nil
		}
		if err == nil {
			fmt.Printf("Retrying in %v...\n", retryTimeout)
		} else {
			fmt.Printf("Error: %v\nRetrying in %v...\n", err, retryTimeout)
		}
		time.Sleep(retryTimeout)
	}
	return fmt.Errorf("condition not met after %d attempts", maxRetries)
}
func checkGitAccount() {
	rebaseStrategy, err := exec.Command("git", "config", "pull.rebase").Output()
	if err != nil {
		// Handle error
		color.Red("Error checking rebase strategy: " + err.Error())
		color.Red("Git rebase strategy not set suggestion: git config pull.rebase false")
		return
	}
	if len(strings.TrimSpace(string(rebaseStrategy))) == 0 {
		color.Red("Git rebase strategy not set suggestion: git config pull.rebase false")
	} else {
		color.Green("Git rebase strategy set to: " + string(rebaseStrategy))
	}

	// Check if the user has set a Git user email
	userEmail, err := exec.Command("git", "config", "user.email").Output()
	if err != nil {
		// Handle error
		color.Red("Error checking user email:" + err.Error())
		color.Red("Git user email not set: git config user.email your.em@il.org")
		return
	}
	if len(strings.TrimSpace(string(userEmail))) == 0 {
		color.Red("Git user email not set: git config user.email your.em@il.org")
	} else {
		color.Green("Git user email set to: " + string(userEmail))
	}

	// Check if the user has set a Git user name
	userName, err := exec.Command("git", "config", "user.name").Output()
	if err != nil {
		// Handle error
		color.Red("Error checking user name: " + err.Error())
		color.Red("Git user name not set: git config user.name anonymous")
		return
	}
	if len(strings.TrimSpace(string(userName))) == 0 {
		color.Red("Git user name not set: gitconfig user.name anonymous")
	} else {
		color.Green("Git user name set to: " + string(userName))
	}
}
func mapToYaml(m map[string]string) string {
	yaml := "apiVersion: v1\n"
	yaml += "kind: Secret\n"
	yaml += "metadata:\n"
	yaml += "  name: " + viper.GetString("secret") + "\n"
	yaml += "  namespace: " + viper.GetString("namespace") + "\n"
	yaml += "type: Opaque\n"
	yaml += "stringData:\n"
	for key, value := range m {
		yaml += fmt.Sprintf("  %s: %s\n", key, value)
	}
	return yaml
}

const (
	cloudflareAPIBase = "https://api.cloudflare.com/client/v4/"
)

type DNSRecord struct {
	ID string `json:"id"`
}

type DNSRecordsResponse struct {
	Result []DNSRecord `json:"result"`
}

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ZonesResponse struct {
	Result []Zone `json:"result"`
}

func getZoneID(apiToken, domainName string) (string, error) {
	fmt.Println("getZoneID() called...")
	client := &http.Client{}

	req, err := http.NewRequest("GET", cloudflareAPIBase+"zones?name="+domainName, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var zonesResponse ZonesResponse
	err = json.Unmarshal(body, &zonesResponse)
	if err != nil {
		return "", err
	}

	if len(zonesResponse.Result) == 0 {
		return "", fmt.Errorf("no zone found with domain name %s", domainName)
	}

	return zonesResponse.Result[0].ID, nil
}

func deleteAllDNSRecords(apiToken, domainName string) error {
	zoneID, err := getZoneID(apiToken, domainName)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", cloudflareAPIBase+"zones/"+zoneID+"/dns_records", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dnsRecords DNSRecordsResponse
	err = json.Unmarshal(body, &dnsRecords)
	if err != nil {
		return err
	}

	for _, record := range dnsRecords.Result {
		req, err := http.NewRequest("DELETE", cloudflareAPIBase+"zones/"+zoneID+"/dns_records/"+record.ID, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Deleted DNS record:", record.ID)
		} else {
			return fmt.Errorf("failed to delete DNS record %s, status code: %d", record.ID, resp.StatusCode)
		}
	}

	return nil
}
func canRunDocker() (bool, error) {
	cmd := exec.Command("docker", "version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false, err
	}
	return true, nil
}
func writeExecutedCommand(commandWithFlags string) {
	configFileName := ".setup.log"
	configPath := filepath.Join(".", configFileName)

	// Open or create the config file, appending to it if it exists
	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// Handle errors opening or creating the file
		fmt.Println("error opening or creating config file:", err)
		return
	}
	defer f.Close()

	// Add a timestamp and a newline to separate commands
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	commandWithFlags = fmt.Sprintf("%s %s\n", timestamp, commandWithFlags)

	_, err = f.WriteString(commandWithFlags)
	if err != nil {
		// Handle errors writing to the file
		fmt.Println("error writing to config file:", err)
		return
	}
}
