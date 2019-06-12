package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "server",
    srcs = ["server.go"],
    importpath = "gvisor.dev/gvisor/pkg/control/server",
    visibility = ["//:sandbox"],
    deps = [
        "//pkg/log",
        "//pkg/unet",
        "//pkg/urpc",
    ],
)
