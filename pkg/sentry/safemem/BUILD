package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "go_test")

go_library(
    name = "safemem",
    srcs = [
        "block_unsafe.go",
        "io.go",
        "safemem.go",
        "seq_unsafe.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/safemem",
    visibility = ["//pkg/sentry:internal"],
    deps = [
        "//pkg/sentry/platform/safecopy",
    ],
)

go_test(
    name = "safemem_test",
    size = "small",
    srcs = [
        "io_test.go",
        "seq_test.go",
    ],
    embed = [":safemem"],
)
