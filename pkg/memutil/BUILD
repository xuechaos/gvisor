package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "memutil",
    srcs = ["memutil_unsafe.go"],
    importpath = "gvisor.dev/gvisor/pkg/memutil",
    visibility = ["//visibility:public"],
    deps = ["@org_golang_x_sys//unix:go_default_library"],
)
