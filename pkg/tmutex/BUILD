package(licenses = ["notice"])

go_library(
    name = "tmutex",
    srcs = ["tmutex.go"],
    importpath = "gvisor.dev/gvisor/pkg/tmutex",
    visibility = ["//:sandbox"],
)

go_test(
    name = "tmutex_test",
    size = "medium",
    srcs = ["tmutex_test.go"],
    embed = [":tmutex"],
)
