package profile

func init() {
	register(&Profile{
		Name:        "docker",
		Description: "Docker profile: installs Docker Engine inside the VM.",
		Provision: func(projectDir string) string {
			return `
    # Install Docker via official repo
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" > /etc/apt/sources.list.d/docker.list
    apt-get update
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

    # Allow vagrant user to use Docker without sudo
    usermod -aG docker vagrant
`
		},
	})
}
