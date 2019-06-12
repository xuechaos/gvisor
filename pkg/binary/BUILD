package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "binary",
    srcs = ["binary.go"],
    importpath = "gvisor.dev/gvisor/pkg/binary",
    visibility = ["//:sandbox"],
)

go_test(
    name = "binary_test",
    size = "small",
    srcs = ["binary_test.go"],
    embed = [":binary"],
)
