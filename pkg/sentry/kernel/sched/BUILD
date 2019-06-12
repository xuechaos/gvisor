package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "sched",
    srcs = [
        "cpuset.go",
        "sched.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/kernel/sched",
    visibility = ["//pkg/sentry:internal"],
)

go_test(
    name = "sched_test",
    size = "small",
    srcs = ["cpuset_test.go"],
    embed = [":sched"],
)
