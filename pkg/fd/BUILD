package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "fd",
    srcs = ["fd.go"],
    importpath = "gvisor.dev/gvisor/pkg/fd",
    visibility = ["//visibility:public"],
)

go_test(
    name = "fd_test",
    size = "small",
    srcs = ["fd_test.go"],
    embed = [":fd"],
)
