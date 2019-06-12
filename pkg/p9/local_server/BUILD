package(licenses = ["notice"])

load("//tools:defs.bzl", "go_binary")

go_binary(
    name = "local_server",
    srcs = ["local_server.go"],
    deps = [
        "//pkg/fd",
        "//pkg/log",
        "//pkg/p9",
        "//pkg/unet",
    ],
)
