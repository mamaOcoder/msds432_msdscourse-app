#!/bin/bash
# Install Docker
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install -y docker-ce
sudo apt-get install unzip

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# (Optional) Authenticate with GCP to pull images from a private Artifact Registry
# echo $DOCKER_AUTH_CONFIG | sudo docker login -u _json_key --password-stdin https://us-east1-docker.pkg.dev

APP_LOCATION=$(curl -s "http://metadata.google.internal/computeMetadata/v1/instance/attributes/app-location" -H "Metadata-Flavor: Google")
gsutil cp "$APP_LOCATION" msdscourse-app.zip

unzip msdscourse-app.zip
# Navigate to the directory containing your docker-compose.yml
cd /msdscourse-app

# Run docker-compose up to start your application
sudo docker-compose up -d
