package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "root",
    srcs = ["root.go"],
    importpath = "gvisor.dev/gvisor/runsc/test/root",
)

go_test(
    name = "root_test",
    size = "small",
    srcs = [
        "cgroup_test.go",
        "chroot_test.go",
        "crictl_test.go",
    ],
    embed = [":root"],
    tags = [
        # Requires docker and runsc to be configured before the test runs.
        # Also test only runs as root.
        "manual",
        "local",
    ],
    deps = [
        "//runsc/cgroup",
        "//runsc/specutils",
        "//runsc/test/root/testdata",
        "//runsc/test/testutil",
        "@com_github_syndtr_gocapability//capability:go_default_library",
    ],
)
