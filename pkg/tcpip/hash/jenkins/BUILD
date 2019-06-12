package(licenses = ["notice"])

go_library(
    name = "jenkins",
    srcs = ["jenkins.go"],
    importpath = "gvisor.dev/gvisor/pkg/tcpip/hash/jenkins",
    visibility = [
        "//visibility:public",
    ],
)

go_test(
    name = "jenkins_test",
    size = "small",
    srcs = [
        "jenkins_test.go",
    ],
    embed = [":jenkins"],
)
