package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "state",
    srcs = [
        "state.go",
        "state_metadata.go",
        "state_unsafe.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/state",
    visibility = ["//pkg/sentry:internal"],
    deps = [
        "//pkg/abi/linux",
        "//pkg/log",
        "//pkg/sentry/inet",
        "//pkg/sentry/kernel",
        "//pkg/sentry/watchdog",
        "//pkg/state/statefile",
        "//pkg/syserror",
    ],
)
