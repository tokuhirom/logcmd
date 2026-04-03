# logcmd

Run a command and append its output to a log file, with ANSI escape sequences stripped.

## Usage

```
logcmd [--no-header] [--no-exit] <logfile> <command> [args...]
```

## Options

| Flag | Description |
|------|-------------|
| `--no-header` | Do not write the timestamp and command line before output |
| `--no-exit` | Do not write the exit code after output |

## Example

```sh
logcmd out.log ls -la
logcmd out.log go test ./...
```

Output appended to `out.log`:

```
=== [2026-04-03 10:00:00] $ ls -la
total 8
drwxr-xr-x  3 user staff   96 Apr  3 10:00 .
-rw-r--r--  1 user staff 1234 Apr  3 10:00 main.go
=== exit code: 0

=== [2026-04-03 10:00:05] $ go test ./...
ok  	github.com/you/yourpkg	0.123s
=== exit code: 0
```

A blank line is inserted between entries when appending to an existing log, making it easy to identify the boundary between runs.

## Install

```sh
go install github.com/you/logcmd@latest
```

Or build locally:

```sh
go build -o logcmd .
```

## Notes

- Output is always appended (`O_APPEND`); the log file is never truncated.
- Both stdout and stderr of the child process are captured and stripped of ANSI escape sequences before being written to the log.
- The child process exit code is propagated, so `logcmd` can be used transparently in scripts.
