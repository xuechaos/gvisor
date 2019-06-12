package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "proto_library")

go_library(
    name = "rpcinet",
    srcs = [
        "device.go",
        "rpcinet.go",
        "socket.go",
        "stack.go",
        "stack_unsafe.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/sentry/socket/rpcinet",
    visibility = ["//pkg/sentry:internal"],
    deps = [
        ":syscall_rpc_go_proto",
        "//pkg/abi/linux",
        "//pkg/binary",
        "//pkg/sentry/arch",
        "//pkg/sentry/context",
        "//pkg/sentry/device",
        "//pkg/sentry/fs",
        "//pkg/sentry/fs/fsutil",
        "//pkg/sentry/inet",
        "//pkg/sentry/kernel",
        "//pkg/sentry/kernel/kdefs",
        "//pkg/sentry/kernel/time",
        "//pkg/sentry/socket",
        "//pkg/sentry/socket/hostinet",
        "//pkg/sentry/socket/rpcinet/conn",
        "//pkg/sentry/socket/rpcinet/notifier",
        "//pkg/sentry/unimpl",
        "//pkg/sentry/usermem",
        "//pkg/syserr",
        "//pkg/syserror",
        "//pkg/tcpip",
        "//pkg/tcpip/buffer",
        "//pkg/unet",
        "//pkg/waiter",
    ],
)

proto_library(
    srcs = ["syscall_rpc.proto"],
    visibility = [
        "//visibility:public",
    ],
)
