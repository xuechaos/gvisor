package(licenses = ["notice"])

go_library(
    name = "ports",
    srcs = ["ports.go"],
    importpath = "gvisor.dev/gvisor/pkg/tcpip/ports",
    visibility = ["//:sandbox"],
    deps = [
        "//pkg/tcpip",
    ],
)

go_test(
    name = "ports_test",
    srcs = ["ports_test.go"],
    embed = [":ports"],
    deps = [
        "//pkg/tcpip",
    ],
)
