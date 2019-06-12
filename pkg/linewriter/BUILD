package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "linewriter",
    srcs = ["linewriter.go"],
    importpath = "gvisor.dev/gvisor/pkg/linewriter",
    visibility = ["//visibility:public"],
)

go_test(
    name = "linewriter_test",
    srcs = ["linewriter_test.go"],
    embed = [":linewriter"],
)
