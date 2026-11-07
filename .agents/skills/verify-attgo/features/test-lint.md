# Test-lint config

`.golangci.yml` sets `run.tests: false`, so a production run never sees a `_test.go` file. The test config is a second run of the same binary.

## Sub-features

- `test-config-clean` runs `.golangci.tests.yml` on the pin and matches the baseline.
- `test-config-violation` adds a test file in the wrong package and gets exactly that finding.
- `test-config-missing` reports the vouch pin, which has no test config.

## How to get to it (user POV)

- Run `./custom-gcl run --config=.golangci.tests.yml` from the consumer root.
- The workflow step passes that config as `--config=<file>`. Other flag forms make the action check `.golangci.yml` instead. This recipe calls the binary directly and passes `--config=.golangci.tests.yml`.

## Driving it with custom-gcl

Preconditions:

- Doctor passed on the go-eth2-client clone, including the test-config linter list.
- `GOTOOLCHAIN=go1.26.10`.
- The clone is a clean pin. No `verifyattgo/` or `foo/` directory.

- **Missing on vouch.** Confirm `$RUN/vouch/.golangci.tests.yml` does not exist. Record the unmet precondition in `$EVIDENCE/vouch-test-lint.txt`. Stop for vouch. Do not create the file.
- **Clean run.** From the go-eth2-client clone, save combined stdout and stderr to `$EVIDENCE/go-eth2-client-test-lint.txt`:

  ```
  ./custom-gcl run --config=.golangci.tests.yml \
    --output.json.path="$EVIDENCE/go-eth2-client-test-lint.json" \
    --color never
  ```

  Exit code is 0. Stdout says `0 issues.` The JSON `Issues` list is empty or absent. That matches the baseline (testpackage 0, testifylint 0, thelper 0, usetesting 0). A typecheck issue fails the drive.
- **Deliberate violation.** Create `foo/foo_test.go` containing `package foo`, an import of `testing`, and `func TestFoo(t *testing.T) {}`. Run the same command with `./foo/` added, JSON path `$EVIDENCE/go-eth2-client-test-violation.json`, and combined output in `$EVIDENCE/go-eth2-client-test-violation.txt`. Exit code is 1. The report has one issue. `FromLinter` is `testpackage`, the file is `foo/foo_test.go`, and `Text` is exactly:

  ```
  package should be `foo_test` instead of `foo`
  ```

  No other linter reports. Delete `foo/` before any later count.
- **Proof.** The clean report and the violation report are both required. The clean report alone does not show the config can fail.

## Gotchas

- The production config enables testpackage because `default` is `all`, and `run.tests: false` still hides test files. A production run that exits 0 does not prove the test config.
- The test config's issue caps are already 0. Keep them. A run that adds `--max-same-issues` with a non-zero value hides repeats.
- `typecheck` may appear as enabled in a run report even when `linters --json` lists four names. Zero typecheck issues is required. A typecheck issue means the toolchain is wrong.
- Vouch adopting a test config is a later change. Until the pin has the file, this feature does not pass for vouch.
