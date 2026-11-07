---
name: verify-attgo
description: "Drive the custom golangci-lint binary built from this checkout against throwaway clones of vouch and go-eth2-client. Use after a change to analyzer behavior, plugin registration, nolint handling, or the custom-binary build."
---

# Verify attgo

Prove the custom binary the way a consumer runs it. analysistest fixtures do not count. Pins, baselines, and the feature index are in [features/README.md](features/README.md). Read that index before a drive.

A pass is a loaded linter, the exit code the feature expects, no `level=error` line in the saved output, and an empty or absent `Report.Error` in the JSON.

Exit 1 means findings. Exit 0 means a loaded linter and an empty report. Exit 3 means the run failed before a summary: an Error line, and it does not print `0 issues.`. Exit 5 means no Go files to analyze: a failed run, no summary. Exit 7 means an error was logged on a run that found nothing. Only exit 7 can print `0 issues.`. Any other non-zero exit is a failed run. Do not record exit 3, 5, or 7 as a count.

## Launch

Done when both clones are at their pins, each has a `custom-gcl` built in that clone, the build log shows a replace to this checkout, and `$EVIDENCE/plugin-revision.txt` exists.

1. Require `git`, `go`, and `golangci-lint` on `PATH`.
2. Set `CHECKOUT` to `git rev-parse --show-toplevel` of this repository. Set `RUN_ID` to `date -u +%Y%m%dT%H%M%SZ`. If `TMPDIR` is unset, set it to the default temp directory. Set `TMP` to `cd "$TMPDIR" && pwd -P` so the path has no symlink. Set `RUN=$TMP/verify-attgo-run/$RUN_ID` and `EVIDENCE=$TMP/verify-attgo-evidence/$RUN_ID`. Create the parent directories if needed. Create `$RUN` and `$EVIDENCE` with `mkdir`, not `mkdir -p`. If either exists, stop and pick a new `RUN_ID`. Write `$RUN_ID` to `$EVIDENCE/run-id.txt` and `$CHECKOUT` to `$EVIDENCE/checkout.txt`. Two runs use different `RUN_ID` values. Leave any other checkout alone, including an existing vouch or go-eth2-client clone.
3. Read the `module` line in this repository's `go.mod`. The organization is the path segment after the host. Clone only the repositories named `vouch` and `go-eth2-client` in that organization, into `$RUN/vouch` and `$RUN/go-eth2-client`. A local clone of the same repository may be passed to `git clone --reference`. Check out the pin from the index with `git checkout --detach <pin>`. If that fails, run `git fetch origin <pin>` and check out again. If that fails, stop. The pin is dead: re-pin, remeasure, and do not drive a different commit. Do not pass `--depth`.
4. In each clone, confirm `git rev-parse HEAD` equals the pin and `git rev-parse --is-shallow-repository` prints `false`.
5. Rewrite that clone's `.custom-gcl.yml`. Keep the top `version:` line (it selects the golangci-lint release). Under `plugins:`, replace the plugin `version:` line with `path: <CHECKOUT>`. Do not add another plugin entry. Save the rewritten file into `$EVIDENCE`. The per-name feature builds a second binary when it needs the sample plugin.
6. Export `GOTOOLCHAIN=go1.26.10` and `GOLANGCI_LINT_CACHE=$RUN/lint-cache` before the build and before every later `custom-gcl` command. Create the cache directory. In each clone, run `golangci-lint custom -v` and save the log. Done when the log contains `go mod edit -replace` and the path is `$CHECKOUT`, and `./custom-gcl` exists.
7. Save `git -C "$CHECKOUT" rev-parse HEAD` and `git -C "$CHECKOUT" status --porcelain` to `$EVIDENCE/plugin-revision.txt`.

## Fresh shell

Each later command runs in a new shell. If `TMPDIR` is unset, set it to the default temp directory. Paste this block, with the literal `RUN_ID` and `CHECKOUT` from Launch, at the top of every command. Do not run `date` again. Do not recover `CHECKOUT` with `git rev-parse` from inside a clone.

```
RUN_ID=<launch value>
CHECKOUT=<launch value>
TMP=$(cd "$TMPDIR" && pwd -P)
RUN=$TMP/verify-attgo-run/$RUN_ID
EVIDENCE=$TMP/verify-attgo-evidence/$RUN_ID
export GOTOOLCHAIN=go1.26.10
export GOLANGCI_LINT_CACHE=$RUN/lint-cache
```

## Doctor

Run this in each clone before any drive. Done when every check below holds. One miss and the clone is not worth driving.

Export `GOTOOLCHAIN=go1.26.10` and `GOLANGCI_LINT_CACHE=$RUN/lint-cache` for every `custom-gcl` command and for the `go version` check. A newer toolchain (Go 1.27) makes the binary report `export data version` typecheck errors and exit 1 with no attgo findings. That is not a result.

Save each output under `$EVIDENCE`.

- `git symbolic-ref refs/remotes/origin/HEAD` prints `refs/remotes/origin/master`, and `git diff --name-status origin/master...HEAD` exits 0. The current-year analyzer tries that ref first. Save both. A different default branch, or a failed diff, makes a finding count not comparable to the baseline. The zero rationale is in the index.
- `GOTOOLCHAIN=go1.26.10 go version` prints `go1.26.10`.
- `go version -m ./custom-gcl` has a `=>` line whose path is `$CHECKOUT` and whose version is `(devel)`. A released module version, or any other path, means this is not the local plugin.
- `./custom-gcl version` contains the top `version:` from this clone's `.custom-gcl.yml` and the marker `-custom-gcl-`.
- `./custom-gcl linters --json --color never` lists `attgo` under `Enabled`. That is the expected custom linter on the production config. Do not require the rest of the `default: all` list: that list changes with the golangci-lint release and is not this plugin.
- In the go-eth2-client clone, `./custom-gcl linters --config=.golangci.tests.yml --json --color never` lists exactly `testifylint`, `testpackage`, `thelper`, and `usetesting` as enabled. `attgo` is absent. Those four names are the expected linters for that config. The vouch pin has no `.golangci.tests.yml`. Record that. It is not a doctor failure. The test-lint feature does not pass for vouch.

Run every later command with the clone as the working directory. An absolute pattern from another directory exits 7 and prints `0 issues.`. A plain absolute directory from outside a module exits 5. Neither linted the module.

## Drive

Open the feature file from the index. Run that file's recipe. A change that touches more than one feature gets each of those recipes. A proof of this skill drives one feature end to end, after launch and doctor on both clones. The index states the drive order.

The harness is `./custom-gcl` in the clone, with the exports in the fresh-shell block. There is no dry-run. Capture the command, the exit code, the combined stdout and stderr, and the JSON report. The summary line `0 issues.` is on stdout. `--output.text.path` stays empty when there are no issues, so do not use it as the proof.

One checkout state per run. After an edit to `$CHECKOUT`, start a new `RUN_ID`. To confirm a replay, set `GL_DEBUG=goanalysis/issues/cache`. A line `analyzing 0/N packages` means the result was loaded from the lint cache. That is not a measurement.

## Evidence

Write proof under `$EVIDENCE`. Name the feature and the consumer in the filename.

A proof contains the command, the exit code, and the JSON report (`--output.json.path`). Text output alone is not proof. Exit 7 can print `0 issues.`. Keep the build log, `plugin-revision.txt`, and the doctor outputs next to the report. A saved drive that contains `level=error` or `attgo_current_year: warning` is not a pass.

Exercise the binary in the consumer clone. An analysistest run, a unit test, or a config that drops the consumer's attgo settings is not this proof. The sample plugin import is scaffolding inside the second binary, not a substitute for the launch binary.

## Cleanup

Paste the fresh-shell block. Refuse to delete unless `$RUN_ID` is non-empty, `$RUN` equals `$TMP/verify-attgo-run/$RUN_ID`, that directory exists, and `cat "$EVIDENCE/run-id.txt"` equals `$RUN_ID`. Then delete `$RUN` and nothing else. That removes the clones, the binaries, the scaffolding, and the lint cache. Leave `$EVIDENCE` in place. After cleanup, the evidence directory still contains the report you will attach.

If a `custom-gcl` or `golangci-lint custom` you started is still running, kill that PID. Do not kill by process name. A killed `golangci-lint custom` leaves a temp directory outside `$RUN`. The build log names it. Delete only that directory.

## Gotchas

The custom binary's version suffix hashes only the plugin root files, so an analyzer edit does not change it. The default lint cache then replays the previous run, and Doctor stays green. `GOLANGCI_LINT_CACHE=$RUN/lint-cache` is what stops that across runs. Cleanup removes that directory.

Build in the throwaway clone for this run. A consumer workflow may restore `./custom-gcl` from a cache whose key hashes only `.custom-gcl.yml`, so a change in this checkout is invisible to that cache. Never copy a binary from a consumer checkout or from that cache.

The current-year analyzer compares file paths with `git rev-parse --show-toplevel`. A clone under a symlinked temp directory makes those paths disagree, and the analyzer then treats changed files as unchanged, with no warning. `$TMP` is the physical path so this skill's clones agree. That does not fix the analyzer. A baseline of 0 does not show the guard working.
