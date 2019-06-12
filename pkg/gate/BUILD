package(licenses = ["notice"])

go_library(
    name = "gate",
    srcs = [
        "gate.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/gate",
    visibility = ["//visibility:public"],
)

go_test(
    name = "gate_test",
    srcs = [
        "gate_test.go",
    ],
    deps = [
        ":gate",
    ],
)
