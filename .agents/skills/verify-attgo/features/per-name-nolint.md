# Per-name registration and nolint

A `//nolint` directive silences the linter it names and leaves every other linter reporting on that line. A name that is a rule inside the attgo plugin is not a linter name.

## Sub-features

- `nolint-one` silences only the named per-name linter.
- `nolint-attgo` silences only the attgo plugin.
- `nolint-rule-name` names a rule and silences nothing.

## How to get to it (user POV)

- Build `custom-gcl` from this checkout inside a consumer clone (Launch).
- Add a `linters.settings.custom` entry named after the analyzer, with no settings.
- Put `//nolint:<name>` on the line to silence.

## Driving it with custom-gcl

Preconditions:

- Doctor passed for this consumer on the launch binary `./custom-gcl`.
- The attgo `enable_*` values in `testdata/project/vouch.yml` or `testdata/project/go-eth2-client.yml` match the consumer's `.golangci.yml` at the pin. One mismatch and the testdata file is stale: stop, and do not drive with it.
- The clone has no `verifyattgo/` directory yet.

- **Second binary.** Append a plugin entry to `.custom-gcl.yml`, on its own lines, after a newline: the same module, `import` set to that module plus `/testdata/sampleplugin`, and the same `path` as the first entry. Run `golangci-lint custom -v --name custom-gcl-pername` with the toolchain Launch exports. Restore `.custom-gcl.yml` from the copy saved at launch. `./custom-gcl` stays the launch binary. `./custom-gcl-pername` is the binary for this feature and for the absent-linter drive. The sample entry is scaffolding. Consumer CI does not have it.
- **Scaffold.** Copy `testdata/project/project.go` from this checkout to `verifyattgo/project.go` in the clone. Copy the matching testdata config to `.golangci.pername.yml` in the clone. The config enables `attgo`, `sample_first`, and `sample_second` only, and it carries that consumer's attgo settings.
- **Confirm the names loaded.** Run `./custom-gcl-pername linters --config=.golangci.pername.yml --json --color never`. `attgo`, `sample_first`, and `sample_second` are all under `Enabled`. Any other result is a failed load. Stop.
- **Run.** From the clone root, with the toolchain exported in Launch, save combined stdout and stderr to `$EVIDENCE/<consumer>-per-name.txt`:

  ```
  ./custom-gcl-pername run --config=.golangci.pername.yml \
    --output.json.path="$EVIDENCE/<consumer>-per-name.json" \
    --uniq-by-line=false --max-issues-per-linter=0 --max-same-issues=0 \
    --color never ./verifyattgo/
  ```

  Exit code is 1. Map each `name = ` line in `project.go` to its line number. Group `FromLinter` by line.

- **nolint-one.** The `byFirst` line reports `attgo` and `sample_second`, not `sample_first`. The `bySecond` line reports `attgo` and `sample_first`, not `sample_second`.
- **nolint-attgo.** The `byAttgo` line reports `sample_first` and `sample_second`, not `attgo`. The `byBoth` line reports `attgo` only.
- **nolint-rule-name.** The `byRuleName` line reports `attgo`, `sample_first`, and `sample_second`. The text log warns that `attgo_raw_string` is an unknown nolint linter. That warning is the proof, not a failure.
- **Control.** The `plain` line reports all three. Drive both consumers. The vouch run also reports `attgo_capital_comment` on the `lowercase on purpose` comment. The go-eth2-client run does not. That difference belongs to the disabled-linter feature; record it if both runs are in hand, and still run that feature's own recipe when it is the target.

## Gotchas

- `//nolint:attgo` silences every rule inside the attgo plugin. It does not silence `sample_first` or `sample_second`.
- A rule name in a directive (`attgo_raw_string`) matches nothing. The warning names the unknown linter. Exit code is still 1 because other findings remain.
- Line numbers in this file are not stable. Match the variable name.
- The sample import in `.custom-gcl.yml` is scaffolding. The consumer configs do not enable those linters. Doctor checks they stay absent there.
- The eight production analyzers still report under the linter name `attgo`. A per-name production analyzer is one registered from this module's root init. None is, today. This recipe is the registration path a consumer build actually runs.
