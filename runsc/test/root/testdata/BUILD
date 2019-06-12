package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "testdata",
    srcs = [
        "busybox.go",
        "containerd_config.go",
        "httpd.go",
        "httpd_mount_paths.go",
        "sandbox.go",
    ],
    importpath = "gvisor.dev/gvisor/runsc/test/root/testdata",
    visibility = [
        "//visibility:public",
    ],
)
