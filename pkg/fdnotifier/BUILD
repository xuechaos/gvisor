package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "fdnotifier",
    srcs = [
        "fdnotifier.go",
        "poll_unsafe.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/fdnotifier",
    visibility = ["//:sandbox"],
    deps = ["//pkg/waiter"],
)
