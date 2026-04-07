package profile

import "github.com/TomHoenderdos/vbox/internal/config"

func init() {
	register(&Profile{
		Name:        "esp",
		Description: "ESP-IDF profile: Espressif IoT Development Framework for ESP32 with USB flashing.",
		Ports:       []config.Port{{Guest: 3333, Host: 3333, Label: "ESP OpenOCD"}},
		NeedsUSB:    true,
		Provision: func(projectDir string) string {
			return `
    # ESP-IDF prerequisites
    apt-get install -y git wget flex bison gperf python3 python3-pip python3-venv \
      cmake ninja-build ccache libffi-dev libssl-dev dfu-util libusb-1.0-0 \
      picocom minicom

    # udev rules for ESP USB devices (CP210x, CH340, FTDI)
    cat > /etc/udev/rules.d/99-esp.rules <<'UDEV'
SUBSYSTEM=="tty", ATTRS{idVendor}=="10c4", MODE="0666"
SUBSYSTEM=="tty", ATTRS{idVendor}=="1a86", MODE="0666"
SUBSYSTEM=="tty", ATTRS{idVendor}=="0403", MODE="0666"
SUBSYSTEM=="usb", ATTR{idVendor}=="303a", MODE="0666"
UDEV
    udevadm control --reload-rules

    # Add vagrant to dialout for serial access
    usermod -aG dialout vagrant

    # Install ESP-IDF (skip if already installed)
    if [ ! -d /home/vagrant/esp/esp-idf ]; then
      su - vagrant -c 'mkdir -p ~/esp && cd ~/esp && git clone --recursive https://github.com/espressif/esp-idf.git -b v5.4 --depth 1'
      su - vagrant -c 'cd ~/esp/esp-idf && ./install.sh all'
    fi
    grep -qF 'esp-idf/export.sh' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "source ~/esp/esp-idf/export.sh" >> ~/.bashrc'
`
		},
	})
}
