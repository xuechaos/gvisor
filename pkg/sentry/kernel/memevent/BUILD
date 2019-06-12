package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "proto_library")

go_library(
    name = "memevent",
    srcs = ["memory_events.go"],
    importpath = "gvisor.dev/gvisor/pkg/sentry/kernel/memevent",
    visibility = ["//:sandbox"],
    deps = [
        ":memory_events_go_proto",
        "//pkg/eventchannel",
        "//pkg/log",
        "//pkg/metric",
        "//pkg/sentry/kernel",
        "//pkg/sentry/usage",
    ],
)

proto_library(
    srcs = ["memory_events.proto"],
    visibility = ["//visibility:public"],
)
