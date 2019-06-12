package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "context",
    srcs = ["context.go"],
    importpath = "gvisor.dev/gvisor/pkg/sentry/context",
    visibility = ["//pkg/sentry:internal"],
    deps = [
        "//pkg/amutex",
        "//pkg/log",
    ],
)
