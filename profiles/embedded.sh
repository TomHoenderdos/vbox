#!/usr/bin/env bash
# Embedded profile: ARM toolchain, serial debuggers, PlatformIO, with USB flashing.

profile_ports() { :; }

profile_usb() { echo "true"; }

profile_provision() {
cat <<'PROVISION'
    apt-get install -y gcc-arm-none-eabi gdb-multiarch openocd picocom minicom screen libusb-1.0-0

    # udev rules for common debug probes (ST-Link, J-Link, CMSIS-DAP)
    cat > /etc/udev/rules.d/99-embedded.rules <<'UDEV'
SUBSYSTEM=="usb", ATTR{idVendor}=="0483", MODE="0666"
SUBSYSTEM=="usb", ATTR{idVendor}=="1366", MODE="0666"
SUBSYSTEM=="usb", ATTR{idVendor}=="0d28", MODE="0666"
SUBSYSTEM=="tty", ATTRS{idVendor}=="10c4", MODE="0666"
SUBSYSTEM=="tty", ATTRS{idVendor}=="1a86", MODE="0666"
SUBSYSTEM=="tty", ATTRS{idVendor}=="0403", MODE="0666"
UDEV
    udevadm control --reload-rules

    # Add vagrant to dialout for serial access
    usermod -aG dialout vagrant

    # Install PlatformIO
    su - vagrant -c 'pip3 install platformio'
PROVISION
}
