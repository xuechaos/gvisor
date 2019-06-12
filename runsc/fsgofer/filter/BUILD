package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library")

go_library(
    name = "filter",
    srcs = [
        "config.go",
        "extra_filters.go",
        "extra_filters_msan.go",
        "extra_filters_race.go",
        "filter.go",
    ],
    importpath = "gvisor.dev/gvisor/runsc/fsgofer/filter",
    visibility = [
        "//runsc:__subpackages__",
    ],
    deps = [
        "//pkg/abi/linux",
        "//pkg/log",
        "//pkg/seccomp",
        "@org_golang_x_sys//unix:go_default_library",
    ],
)
