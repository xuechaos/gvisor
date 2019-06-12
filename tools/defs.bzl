"""Wrappers that include standard generators.

The recommended way is to use the go_library rule defined below with mostly
identical configuration as the native go_library rule.

load("//tools:defs.bzl", "go_library")

go_library(
    name = "foo",
    srcs = ["foo.go"],
)
"""

load("//tools/go_stateify:defs.bzl", "go_stateify")
load("@io_bazel_rules_go//go:def.bzl", _go_binary = "go_binary", _go_library = "go_library", _go_test = "go_test")
load("@io_bazel_rules_go//proto:def.bzl", _go_proto_library = "go_proto_library")

def go_binary(name, **kwargs):
    """Wraps the standard go_binary.

    Args:
      name: the rule name.
      **kwargs: standard go_binary arguments.
    """
    _go_binary(
        name = name,
        **kwargs
    )

def go_library(name, srcs, deps = [], imports = [], stateify = True, **kwargs):
    """Wraps the standard go_library and does stateification.

    Args:
      name: the rule name.
      srcs: the library sources.
      deps: the library dependencies.
      imports: imports required for stateify.
      stateify: whether statify is enabled (default: true).
      **kwargs: standard go_library arguments.
    """
    if stateify:
        # Only do stateification for non-state packages without manual autogen.
        go_stateify(
            name = name + "_state_autogen",
            srcs = [src for src in srcs if src.endswith(".go")],
            imports = imports,
            package = name,
            out = name + "_state_autogen.go",
        )
        all_srcs = srcs + [name + "_state_autogen.go"]
        if "//pkg/state" not in deps:
            all_deps = deps + ["//pkg/state"]
        else:
            all_deps = deps
    else:
        all_deps = deps
        all_srcs = srcs
    _go_library(
        name = name,
        srcs = all_srcs,
        deps = all_deps,
        **kwargs
    )

def go_test(name, **kwargs):
    """Wraps the standard go_test.

    Args:
      name: the rule name.
      **kwargs: standard go_test arguments.
    """
    _go_test(
        name = name,
        **kwargs
    )

def proto_library(srcs, **kwargs):
    """Wraps the standard proto_library.

    Args:
      srcs: the proto sources; the first proto file should will be used to
        generate the rule name. If this file is "${src}.proto", then rules
        named "${src}_proto" and "${src}_go_proto" will be generated.
      **kwargs: standard proto_library and go_proto_library arguments.
    """

    import_path = kwargs.pop("import_path")
    name = srcs[0].split(".")[0]
    native.proto_library(
        name = name + "_proto",
        srcs = srcs,
        deps = kwargs.pop("deps", None),
        **kwargs
    )
    _go_proto_library(
        name = name + "_go_proto",
        proto = ":" + name + "_proto",
        import_path = import_path,
        **kwargs
    )
