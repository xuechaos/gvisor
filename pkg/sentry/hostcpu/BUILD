package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "hostcpu",
    srcs = [
        "getcpu_amd64.s",
        "hostcpu.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/hostcpu",
    visibility = ["//:sandbox"],
)

go_test(
    name = "hostcpu_test",
    size = "small",
    srcs = ["hostcpu_test.go"],
    embed = [":hostcpu"],
)
