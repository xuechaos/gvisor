package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "proto_library")

proto_library(
    srcs = ["unimplemented_syscall.proto"],
    visibility = ["//visibility:public"],
    deps = ["//pkg/sentry/arch:registers_proto"],
)

go_library(
    name = "unimpl",
    srcs = ["events.go"],
    importpath = "gvisor.dev/gvisor/pkg/sentry/unimpl",
    visibility = ["//:sandbox"],
    deps = [
        "//pkg/log",
        "//pkg/sentry/context",
    ],
)
