package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "anon",
    srcs = [
        "anon.go",
        "device.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/fs/anon",
    visibility = ["//pkg/sentry:internal"],
    deps = [
        "//pkg/abi/linux",
        "//pkg/sentry/context",
        "//pkg/sentry/device",
        "//pkg/sentry/fs",
        "//pkg/sentry/fs/fsutil",
        "//pkg/sentry/usermem",
    ],
)
