# run setup on windows

## install choco package manager
in powershell:
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

choco install gh cloudflared git terraform kubernetes-cli -y
```

# run setup in docker

```
docker-compose run setup go run /app/setup.go -h
```
