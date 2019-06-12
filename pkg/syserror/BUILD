package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "syserror",
    srcs = ["syserror.go"],
    importpath = "gvisor.dev/gvisor/pkg/syserror",
    visibility = ["//visibility:public"],
)

go_test(
    name = "syserror_test",
    srcs = ["syserror_test.go"],
    deps = [
        ":syserror",
    ],
)
