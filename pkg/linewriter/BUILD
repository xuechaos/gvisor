load("//tools/go_stateify:defs.bzl", "go_library", "go_test")

package(licenses = ["notice"])

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
