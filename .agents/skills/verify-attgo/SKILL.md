---
name: verify-attgo
description: "Drive the custom golangci-lint binary built from this checkout against throwaway clones of vouch and go-eth2-client. Use after a change to analyzer behavior, plugin registration, nolint handling, or the custom-binary build."
---

# Verify attgo

Prove the custom binary the way a consumer runs it. analysistest fixtures do not count. Pins, baselines, and the feature index are in [features/README.md](features/README.md). Read that index before a drive.

A pass is a loaded linter and the exit code the feature expects. Exit 3 means the linter did not load. Exit 7 means the run logged an error. Either one can print `0 issues`. Record that as a failed run, not as a pass and not as a count.

## Launch

Done when both clones are at their pins, each has a `custom-gcl` built in that clone, and the build log shows a replace to this checkout.

1. Require `git`, `go`, and `golangci-lint` on `PATH`.
2. Set `CHECKOUT` to `git rev-parse --show-toplevel` of this repository. Set `RUN_ID` to `date -u +%Y%m%dT%H%M%SZ`. If `TMPDIR` is unset, set it to the default temp directory. Set `TMP` to `cd "$TMPDIR" && pwd -P` so the path has no symlink. Set `RUN=$TMP/verify-attgo-run/$RUN_ID` and `EVIDENCE=$TMP/verify-attgo-evidence/$RUN_ID`. Create both directories. Two runs use different `RUN_ID` values. Leave any other checkout alone, including an existing vouch or go-eth2-client clone.
3. Read the `module` line in this repository's `go.mod`. The organization is the path segment after the host. Clone only the repositories named `vouch` and `go-eth2-client` in that organization, into `$RUN/vouch` and `$RUN/go-eth2-client`. A local clone of the same repository may be passed to `git clone --reference`. Check out the pin from the index with `git checkout --detach <pin>`. Do not pass `--depth`.
4. In each clone, confirm `git rev-parse HEAD` equals the pin and `git rev-parse --is-shallow-repository` prints `false`.
5. Rewrite that clone's `.custom-gcl.yml`. Keep the top `version:` line (it selects the golangci-lint release). Under `plugins:`, replace the plugin `version:` line with `path: <CHECKOUT>`. Do not add another plugin entry. Save the rewritten file into `$EVIDENCE`. The per-name feature builds a second binary when it needs the sample plugin.
6. Export `GOTOOLCHAIN=go1.26.10` before the build and before every later `custom-gcl` command. In each clone, run `golangci-lint custom -v` and save the log. Done when the log contains `go mod edit -replace` and the path is `$CHECKOUT`, and `./custom-gcl` exists.

## Doctor

Run this in each clone before any drive. Done when every check below holds. One miss and the clone is not worth driving.

Export `GOTOOLCHAIN=go1.26.10` for every `custom-gcl` command and for the `go version` check. A newer toolchain (Go 1.27) makes the binary report `export data version` typecheck errors and exit 1 with no attgo findings. That is not a result.

Save each output under `$EVIDENCE`.

- `git rev-parse --verify origin/master` succeeds, and `git diff --name-status origin/master...HEAD` exits 0. The current-year analyzer uses that diff. A failed diff makes it check every file, and a finding count then is not comparable to the baseline.
- `GOTOOLCHAIN=go1.26.10 go version` prints `go1.26.10`.
- `go version -m ./custom-gcl` has a `=>` line whose path is `$CHECKOUT` and whose version is `(devel)`. A released module version, or any other path, means this is not the local plugin.
- `./custom-gcl version` contains the top `version:` from this clone's `.custom-gcl.yml` and the marker `-custom-gcl-`.
- `./custom-gcl linters --json --color never` lists `attgo` under `Enabled`. That is the expected custom linter on the production config. Do not require the rest of the `default: all` list: that list changes with the golangci-lint release and is not this plugin.
- In the go-eth2-client clone, `./custom-gcl linters --config=.golangci.tests.yml --json --color never` lists exactly `testifylint`, `testpackage`, `thelper`, and `usetesting` as enabled. `attgo` is absent. Those four names are the expected linters for that config. The vouch pin has no `.golangci.tests.yml`. Record that. It is not a doctor failure. The test-lint feature owns that gap.

Run every later command with the clone as the working directory. A pattern that is an absolute path from another directory exits 7, prints `0 issues`, and did not lint the module.

## Drive

Open the feature file from the index. Run that file's recipe. A change that touches more than one feature gets each of those recipes. A proof of this skill drives one feature end to end, after launch and doctor on both clones.

The harness is `./custom-gcl` in the clone, with `GOTOOLCHAIN=go1.26.10`. There is no dry-run. Capture the command, the exit code, the combined stdout and stderr, and the JSON report. The summary line `0 issues.` is on stdout. `--output.text.path` stays empty when there are no issues, so do not use it as the proof.

## Evidence

Write proof under `$EVIDENCE`. Name the feature and the consumer in the filename.

A proof contains the command, the exit code, and the JSON report (`--output.json.path`). Text output alone is not proof: a failed load can print `0 issues`. Keep the build log and the doctor outputs next to the report.

Exercise the binary in the consumer clone. An analysistest run, a unit test, or a config that drops the consumer's attgo settings is not this proof. The sample plugin import is scaffolding inside the real binary, not a substitute for it.

## Cleanup

Delete `$RUN` and nothing else. That removes the clones, the binaries, and the scaffolding files. Leave `$EVIDENCE` in place. After cleanup, the evidence directory still contains the report you will attach.

If a `custom-gcl` you started is still running, kill that PID. Do not kill by process name.

## Gotchas

Build in the throwaway clone for this run. A consumer workflow may restore `./custom-gcl` from a cache whose key hashes only `.custom-gcl.yml`, so a change in this checkout is invisible to that cache. Never copy a binary from a consumer checkout or from that cache.

The current-year analyzer compares file paths with `git rev-parse --show-toplevel`. A clone under a symlinked temp directory makes those paths disagree, and the analyzer then treats changed files as unchanged. `$TMP` is the physical path so they agree.

## Helpers

None. The commands in this skill and in the feature files are the harness.
