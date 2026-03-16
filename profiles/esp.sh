#!/usr/bin/env bash
# ESP-IDF profile: Espressif IoT Development Framework for ESP32/S2/S3/C3/C6 with USB flashing.

profile_ports() {
  echo "3333:3333:ESP OpenOCD"
}

profile_usb() { echo "true"; }

profile_provision() {
cat <<'PROVISION'
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
    grep -qF 'get_idf' /home/vagrant/.bashrc || \
      su - vagrant -c 'echo "alias get_idf=\"source ~/esp/esp-idf/export.sh\"" >> ~/.bashrc'
PROVISION
}
