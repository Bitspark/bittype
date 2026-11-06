#!/usr/bin/env python3
"""Capture existing Go test invocations without changing their result or selection.

The only result files are invocation.json, stdout.jsonl, and stderr.log. Context
is captured independently by the context command for the reporting host. This is
CI observation capture, not a source-equivalence or profile-attestation service.
"""

import argparse
import json
import os
from pathlib import Path
import shlex
import signal
import subprocess
import sys
import uuid


def command_json(args, cwd):
    result = subprocess.run(args, cwd=cwd, capture_output=True, text=True)
    if result.returncode:
        raise RuntimeError(f"{args[0]} {args[1]} failed: {result.stderr.strip()}")
    decoder = json.JSONDecoder()
    values = []
    text = result.stdout
    while text.strip():
        text = text.lstrip()
        value, end = decoder.raw_decode(text)
        values.append(value)
        text = text[end:]
    return values


def selected_build_flags(args):
    if len(args) < 2 or Path(args[0]).name not in ("go", "go.exe") or args[1] != "test":
        raise ValueError("capture requires an explicit go test command")
    switches = {"-race", "-msan", "-asan", "-trimpath", "-linkshared", "-cover"}
    valued = {
        "-tags", "-mod", "-modfile", "-overlay", "-compiler", "-installsuffix",
        "-pkgdir", "-toolexec", "-buildmode", "-asmflags", "-gcflags", "-ldflags",
        "-covermode", "-coverpkg",
    }
    flags = []
    index = 2
    while index < len(args):
        argument = args[index]
        if argument == "-args":
            break
        name, equal, _ = argument.partition("=")
        if name in switches:
            flags.append(argument)
        elif name in valued:
            flags.append(argument)
            if not equal:
                index += 1
                if index == len(args):
                    raise ValueError(f"missing value for {name}")
                flags.append(args[index])
        index += 1
    return flags


def module_revision(module):
    revision = (module.get("Origin") or {}).get("Hash", "")
    mod_file = module.get("GoMod", "")
    if not revision and mod_file.endswith(".mod"):
        info = Path(mod_file).with_suffix(".info")
        if info.is_file():
            revision = (json.loads(info.read_text(encoding="utf-8")).get("Origin") or {}).get("Hash", "")
    return revision


def source_context(root, cwd, harness, build_flags):
    result = {
        "revision": "", "module_path": "", "module_revision": "", "harness": harness,
        "go_version": "", "goos": "", "goarch": "", "cgo_enabled": "",
        "go_flags": "", "build_flags": build_flags, "workspace": "",
        "modules": [], "replacements": [], "input_state": "unknown",
        "input_changes": [], "complete": False, "problems": [],
    }
    problems = result["problems"]
    try:
        revision = subprocess.check_output(["git", "-C", str(root), "rev-parse", "HEAD"], text=True).strip()
        result["revision"] = revision
        # The harness source comes from this checkout. Independently installed
        # modules retain their selected versions/revisions below and a different
        # harness/module identity; they cannot match repository declarations.
        result["module_revision"] = revision
    except (OSError, subprocess.SubprocessError) as error:
        problems.append(f"checkout revision unavailable: {error}")
    try:
        environment = command_json([
            "go", "env", "-json", "GOVERSION", "GOOS", "GOARCH", "CGO_ENABLED", "GOFLAGS", "GOWORK",
        ], cwd)[0]
        for name, target in (
            ("GOVERSION", "go_version"), ("GOOS", "goos"), ("GOARCH", "goarch"),
            ("CGO_ENABLED", "cgo_enabled"), ("GOFLAGS", "go_flags"), ("GOWORK", "workspace"),
        ):
            result[target] = environment[name]
        effective = build_flags + shlex.split(result["go_flags"])
        if any(flag.partition("=")[0] == "-overlay" for flag in effective):
            problems.append("overlay source inputs are not qualified by checkout identity")
        modules = command_json(["go", "list", *build_flags, "-m", "-json", "all"], cwd)
        for module in modules:
            if module.get("Main"):
                if result["module_path"]:
                    problems.append("multiple main workspace modules require an explicit source model")
                result["module_path"] = module["Path"]
            item = {"path": module["Path"], "version": module.get("Version", "")}
            item["revision"] = result["revision"] if module.get("Main") else module_revision(module)
            result["modules"].append(item)
            replacement = module.get("Replace")
            if replacement:
                result["replacements"].append({
                    "path": module["Path"], "version": module.get("Version", ""),
                    "replacement": replacement["Path"], "with_version": replacement.get("Version", ""),
                    "revision": module_revision(replacement),
                })
                if not replacement.get("Version"):
                    problems.append(f"local replacement source is not qualified: {module['Path']}")
        if not result["module_path"]:
            problems.append("main module was not selected")
        result["modules"].sort(key=lambda item: item["path"])
        result["replacements"].sort(key=lambda item: (item["path"], item["version"]))
    except (OSError, ValueError, KeyError, IndexError, RuntimeError) as error:
        problems.append(f"Go build context unavailable: {error}")
    try:
        # Report existing modifications, including dependency/formatter changes.
        # This introduces no checksum, equivalent-revision, or clean-up workflow.
        status = subprocess.check_output([
            "git", "-C", str(root), "status", "--porcelain=v1", "--untracked-files=all", "-z",
        ])
        result["input_changes"] = sorted(part.decode("utf-8", errors="replace") for part in status.split(b"\0") if part)
        result["input_state"] = "modified" if status else "clean"
    except (OSError, subprocess.SubprocessError) as error:
        problems.append(f"source changes unavailable: {error}")
    result["complete"] = not problems
    return result


def write_json(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, ensure_ascii=True) + "\n", encoding="utf-8")


def capture(args):
    command = args.command
    if command and command[0] == "--":
        command = command[1:]
    build_flags = selected_build_flags(command)
    if "-json" not in command[2:]:
        raise ValueError("go test capture requires -json; existing selections must be supplied explicitly")
    output = args.output.resolve()
    if output.exists():
        raise ValueError("invocation output already exists; use a distinct command output directory")
    output.mkdir(parents=True)
    invocation = {
        "schema": 1,
        "source": source_context(args.root, args.cwd, args.harness, build_flags),
        "run_id": os.environ.get("GITHUB_RUN_ID", str(uuid.uuid4())),
        "attempt": os.environ.get("GITHUB_RUN_ATTEMPT", "1"),
        "job": os.environ.get("GITHUB_JOB", "manual"),
        "matrix": args.matrix,
        "command_id": args.command_id,
        "args": command,
        "environment": [{"name": name, "value": os.environ[name]} for name in sorted(set(args.env)) if name in os.environ],
        "exit_code": None,
        "capture_complete": False,
    }
    write_json(output / "invocation.json", invocation)
    print(f"Capturing {args.command_id}: {shlex.join(command)}", flush=True)
    process = None
    old_handlers = {}

    def forward(signum, frame):
        if process is not None:
            process.send_signal(signum)

    try:
        for signum in (signal.SIGTERM, signal.SIGINT):
            old_handlers[signum] = signal.signal(signum, forward)
        with (output / "stdout.jsonl").open("wb") as stdout, (output / "stderr.log").open("wb") as stderr:
            process = subprocess.Popen(command, cwd=args.cwd, stdout=stdout, stderr=stderr)
            invocation["exit_code"] = process.wait()
        invocation["capture_complete"] = True
    finally:
        for signum, handler in old_handlers.items():
            signal.signal(signum, handler)
        write_json(output / "invocation.json", invocation)
    code = invocation["exit_code"]
    print(f"{args.command_id} exited {code}; raw observations: {output}", flush=True)
    if code:
        return code if code > 0 else 128 - code
    return 0 if invocation["source"]["complete"] else 1


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    subparsers = parser.add_subparsers(dest="action", required=True)
    for name in ("context", "run"):
        operation = subparsers.add_parser(name)
        operation.add_argument("--root", type=Path, required=True)
        operation.add_argument("--cwd", type=Path)
        operation.add_argument("--harness", default="repository")
        operation.add_argument("--output", type=Path, required=True)
        if name == "context":
            operation.add_argument("--build-flag", action="append", default=[])
        else:
            operation.add_argument("--command-id", required=True)
            operation.add_argument("--matrix", default="")
            operation.add_argument("--env", action="append", default=[])
            operation.add_argument("command", nargs=argparse.REMAINDER)
    args = parser.parse_args()
    args.root = args.root.resolve()
    args.cwd = (args.cwd or args.root).resolve()
    if args.action == "context":
        context = source_context(args.root, args.cwd, args.harness, args.build_flag)
        write_json(args.output, context)
        return 0 if context["complete"] else 1
    return capture(args)


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (OSError, ValueError, RuntimeError) as error:
        print(f"Observation capture failed: {error}", file=sys.stderr)
        raise SystemExit(1)
