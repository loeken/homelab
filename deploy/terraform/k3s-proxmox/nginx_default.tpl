upstream backend {
  server ${k3s_vm_default_ipv4_address}:30080;
}

upstream backend_ssl {
  server ${k3s_vm_default_ipv4_address}:30443;
}

server {
  listen 80;
  server_name _; # replace with your domain name

  location / {
    proxy_pass http://backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  }
}

server {
  listen 443;
  server_name _; # replace with your domain name

  location / {
    proxy_pass https://backend_ssl;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    # Enable SSL passthrough
    proxy_ssl_session_reuse off;
    proxy_ssl_name $host;
    proxy_ssl_server_name on;
  }
}