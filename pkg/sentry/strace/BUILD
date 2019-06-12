package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "proto_library")

go_library(
    name = "strace",
    srcs = [
        "capability.go",
        "clone.go",
        "futex.go",
        "linux64.go",
        "open.go",
        "poll.go",
        "ptrace.go",
        "signal.go",
        "socket.go",
        "strace.go",
        "syscalls.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/strace",
    visibility = ["//:sandbox"],
    deps = [
        ":strace_go_proto",
        "//pkg/abi",
        "//pkg/abi/linux",
        "//pkg/binary",
        "//pkg/bits",
        "//pkg/eventchannel",
        "//pkg/seccomp",
        "//pkg/sentry/arch",
        "//pkg/sentry/kernel",
        "//pkg/sentry/kernel/kdefs",
        "//pkg/sentry/socket/control",
        "//pkg/sentry/socket/epsocket",
        "//pkg/sentry/socket/netlink",
        "//pkg/sentry/syscalls/linux",
        "//pkg/sentry/usermem",
    ],
)

proto_library(
    srcs = ["strace.proto"],
    visibility = ["//visibility:public"],
)
