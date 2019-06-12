package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "urpc",
    srcs = ["urpc.go"],
    importpath = "gvisor.dev/gvisor/pkg/urpc",
    visibility = ["//:sandbox"],
    deps = [
        "//pkg/fd",
        "//pkg/log",
        "//pkg/unet",
    ],
)

go_test(
    name = "urpc_test",
    size = "small",
    srcs = ["urpc_test.go"],
    embed = [":urpc"],
    deps = ["//pkg/unet"],
)
