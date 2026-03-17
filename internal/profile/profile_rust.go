package profile

func init() {
	register(&Profile{
		Name:        "rust",
		Description: "Rust profile: installs Rust via rustup.",
		Excludes:    []string{"target/"},
		Provision: func(projectDir string) string {
			return `
    if [ ! -d /home/vagrant/.cargo ]; then
      su - vagrant -c 'curl --proto "=https" --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y'
    fi
    grep -qF '.cargo/env' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "source $HOME/.cargo/env" >> ~/.bashrc'
`
		},
	})
}
