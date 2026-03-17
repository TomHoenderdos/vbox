package profile

func init() {
	register(&Profile{
		Name:        "c",
		Description: "C/C++ profile: installs compilers, debuggers, and analyzers.",
		Excludes:    []string{"build/"},
		Provision: func(projectDir string) string {
			return `
    apt-get install -y gcc g++ make cmake ninja-build autoconf automake libtool \
      gdb valgrind clang clang-format clang-tidy cppcheck doxygen \
      libboost-all-dev libcmocka-dev lcov libncurses5-dev libncursesw5-dev
`
		},
	})
}
