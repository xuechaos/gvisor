package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "amutex",
    srcs = ["amutex.go"],
    importpath = "gvisor.dev/gvisor/pkg/amutex",
    visibility = ["//:sandbox"],
)

go_test(
    name = "amutex_test",
    size = "small",
    srcs = ["amutex_test.go"],
    embed = [":amutex"],
    tags = ["flaky"],
)
