#!/usr/bin/env bash

apt-get update && apt-get upgrade -y
apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

# Install docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get install -y docker-ce docker-ce-cli containerd.io

usermod -aG docker ubuntu

# Install docker-compose
curl -L "https://github.com/docker/compose/releases/download/1.27.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

su ubuntu
git clone https://github.com/insight-infrastructure/prysm-docker-compose.git /home/ubuntu/prysm-docker-compose

echo '${keystore_password}' >> /home/ubuntu/prysm-docker-compose/validator/wallet-password
