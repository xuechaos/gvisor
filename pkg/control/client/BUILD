package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "client",
    srcs = [
        "client.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/control/client",
    visibility = ["//:sandbox"],
    deps = [
        "//pkg/unet",
        "//pkg/urpc",
    ],
)
