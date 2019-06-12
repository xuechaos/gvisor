package(licenses = ["notice"])

load("//tools:defs.bzl", "go_library", "proto_library")

go_library(
    name = "eventchannel",
    srcs = [
        "eventchannel.go",
    ],
    importpath = "gvisor.dev/gvisor/pkg/eventchannel",
    visibility = ["//:sandbox"],
    deps = [
        ":eventchannel_go_proto",
        "//pkg/log",
        "//pkg/unet",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
    ],
)

proto_library(
    srcs = ["eventchannel.proto"],
    visibility = ["//:sandbox"],
)
