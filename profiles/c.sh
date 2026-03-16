#!/usr/bin/env bash
# C/C++ profile: installs compilers, debuggers, and analyzers.

profile_ports() { :; }

profile_provision() {
cat <<'PROVISION'
    apt-get install -y gcc g++ make cmake ninja-build autoconf automake libtool \
      gdb valgrind clang clang-format clang-tidy cppcheck doxygen \
      libboost-all-dev libcmocka-dev lcov libncurses5-dev libncursesw5-dev
PROVISION
}
