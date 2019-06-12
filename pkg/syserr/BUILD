package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "syserr",
    srcs = [
        "host_linux.go",
        "netstack.go",
        "syserr.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/syserr",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/abi/linux",
        "//pkg/syserror",
        "//pkg/tcpip",
    ],
)
