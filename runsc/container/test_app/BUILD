package(licenses = ["notice"])

load("//tools:defs.bzl", "go_binary")

go_binary(
    name = "test_app",
    testonly = 1,
    srcs = ["test_app.go"],
    pure = "on",
    visibility = ["//runsc/container:__pkg__"],
    deps = [
        "//runsc/test/testutil",
        "@com_github_google_subcommands//:go_default_library",
    ],
)
