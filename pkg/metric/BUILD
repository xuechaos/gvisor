package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test", "proto_library")

go_library(
    name = "metric",
    srcs = ["metric.go"],
    importpath = "gvisor.dev/gvisor/pkg/metric",
    visibility = ["//:sandbox"],
    deps = [
        ":metric_go_proto",
        "//pkg/eventchannel",
        "//pkg/log",
    ],
)

proto_library(
    srcs = ["metric.proto"],
    visibility = ["//:sandbox"],
)

go_test(
    name = "metric_test",
    srcs = ["metric_test.go"],
    embed = [":metric"],
    deps = [
        ":metric_go_proto",
        "//pkg/eventchannel",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)
