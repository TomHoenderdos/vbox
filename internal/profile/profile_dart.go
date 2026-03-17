package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "dart",
		Description: "Dart profile: installs Dart SDK and Flutter for web/server development.",
		Ports:       []config.Port{{Guest: 8080, Host: 8080, Label: "Flutter web"}},
		Excludes:    []string{".dart_tool/", "build/"},
		Provision: func(projectDir string) string {
			return `
    # Dart/Flutter dependencies
    apt-get install -y clang cmake ninja-build pkg-config libgtk-3-dev liblzma-dev libstdc++-12-dev

    # Install Flutter (includes Dart SDK)
    if [ ! -d /home/vagrant/.flutter ]; then
      su - vagrant -c 'git clone https://github.com/flutter/flutter.git ~/.flutter -b stable --depth 1'
      su - vagrant -c 'export PATH=$HOME/.flutter/bin:$PATH && flutter precache --web'
    fi
    grep -qF '.flutter/bin' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "export PATH=\$HOME/.flutter/bin:\$HOME/.flutter/bin/cache/dart-sdk/bin:\$PATH" >> ~/.bashrc'
`
		},
	})
}
