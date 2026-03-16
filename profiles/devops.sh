#!/usr/bin/env bash
# DevOps profile: Kubernetes, Terraform, Ansible, AWS CLI.

profile_ports() { :; }

profile_provision() {
cat <<'PROVISION'
    # kubectl
    curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.31/deb/Release.key | gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg
    echo "deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.31/deb/ /" > /etc/apt/sources.list.d/kubernetes.list

    # Helm
    curl -fsSL https://baltocdn.com/helm/signing.asc | gpg --dearmor -o /etc/apt/keyrings/helm.gpg
    echo "deb [signed-by=/etc/apt/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" > /etc/apt/sources.list.d/helm.list

    apt-get update
    apt-get install -y kubectl helm ansible awscli

    # Terraform
    curl -fsSL https://releases.hashicorp.com/terraform/1.9.8/terraform_1.9.8_linux_arm64.zip -o /tmp/terraform.zip
    unzip -q /tmp/terraform.zip -d /usr/local/bin
    rm /tmp/terraform.zip
PROVISION
}
