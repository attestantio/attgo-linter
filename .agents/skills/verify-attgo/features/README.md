# attgo verification map

This directory is the maintained source for verifying the custom binary a consumer runs. Read this index before driving, then use the matching feature file as the recipe.

## Baseline preconditions

- Launch both clones as `SKILL.md` describes, from a clean pin, and pass doctor before a drive.
- Pins:

  | Consumer | Pin | golangci-lint release in `.custom-gcl.yml` |
  | --- | --- | --- |
  | vouch | `83758284525e3abb5bf46a1680a944020dbf831f` | v2.12.2 |
  | go-eth2-client | `fb4a89371097364564680a54fe0ca0ebcc57d7c1` | v2.13.1 |

- The go-eth2-client pin in the table is on branch `test-lint-config`. That commit contains `.golangci.tests.yml`. It is not the default branch. A clone of the default branch does not contain it until Launch fetches the pin. When the file is on the default branch, move this pin to that commit, remeasure, replace the numbers below, and reword the current-year note. If the pin is not on the remote, stop and re-pin. Do not drive a different commit.
- Use the toolchain Launch exports (`GOTOOLCHAIN=go1.26.10`) for every `custom-gcl` run. Go 1.27 makes this golangci-lint release report export-data errors and exit 1 with no attgo findings.
- Run from the clone root. Pass `--uniq-by-line=false --max-issues-per-linter=0 --max-same-issues=0` whenever the recipe counts findings. The defaults stop at 50 per linter and 3 per identical text, which hides findings.
- The vouch pin has no `.golangci.tests.yml`. Doctor records that and continues. The test-lint feature does not pass for vouch: it reports the unmet precondition and stops. Do not add the file.

## Baselines

Measured with the plugin at `0b7bfd3e5d282aee95632c4bc674fe2e75625ed0`, on a clean pin, `GOTOOLCHAIN=go1.26.10`, `--enable-only=attgo`, caps disabled. Split each issue `Text` on the first `:` to get the analyzer name. `FromLinter` is `attgo` for all eight.

vouch at its pin:

| Analyzer | Findings |
| --- | --- |
| attgo_capital_comment | 9 |
| attgo_no_pkg_logger | 3 |
| attgo_current_year | 0 |
| attgo_enum_iota | 0 |
| attgo_func_opts | 0 |
| attgo_raw_string | 0 |
| attgo_struct_field_order | 0 |
| attgo_interface_check | 0 |

go-eth2-client at its pin:

| Analyzer or linter | Findings |
| --- | --- |
| attgo_enum_iota, attgo_current_year, attgo_func_opts, attgo_raw_string, attgo_interface_check | 0, rule on |
| attgo_no_pkg_logger, attgo_capital_comment, attgo_struct_field_order | off in this config, not a measured zero |
| testpackage | 0 |
| testifylint | 0 |
| thelper | 0 |
| usetesting | 0 |

`attgo_current_year` is 0 whenever the diff against `origin/master` has no production Go file, including an empty diff. Test files do not count: the production config does not lint tests. On the vouch pin the diff is empty. On this go-eth2-client pin the diff is the test config, the workflow, and `_test.go` files. Reword this sentence when the pin moves. A 0 is not evidence the analyzer ran.

The same limit covers `attgo_enum_iota`, `attgo_func_opts`, `attgo_struct_field_order`, and `attgo_interface_check`. A row equal to the baseline does not show that analyzer ran. These rows catch a new finding. They do not catch a dropped analyzer.

A mismatch fails the feature. Replace a number only when the analyzer diff explains the new count and a before/after finding list is saved next to it. Record the plugin commit from `$EVIDENCE/plugin-revision.txt` beside the new number. A dirty status in that file is not a baseline for a clean commit. A count taken while `verifyattgo/` or `foo/` exists in the clone is not a baseline.

## Driving conventions

- Start from the clean pin unless the feature says to add scaffolding.
- Copy scaffolding from this checkout's `testdata/project/`. Do not edit the consumer's tracked files to remove it. `rm -rf verifyattgo/` and `rm -rf foo/` in the clone are allowed.
- Drive finding-counts, then the clean test-lint run, before any recipe that creates `verifyattgo/` or `foo/`. A count taken while either directory exists is invalid.
- Treat exit 1 as findings, not as a crashed command. Treat exit 0 as a loaded linter with an empty report. Any other exit is a failed run. See SKILL.md.
- Record the feature id, the consumer, the command, the exit code, and the JSON report with every artifact.
- An unreachable path is reported with the command and the unmet precondition. A skipped consumer is not verified through the other one.

## Features

- [Per-name registration and nolint](./per-name-nolint.md) covers a directive that names one linter, a directive that names the attgo plugin, and a directive that names a rule instead of a linter.
- [Disabled or absent linter](./disabled-or-absent.md) covers a rule the config turns off, and a linter the config does not name.
- [Findings and the exit code](./findings-exit.md) covers a finding exiting 1, and a stock binary exiting 3 with no findings.
- [Test-lint config](./test-lint.md) covers `.golangci.tests.yml` on the go-eth2-client pin, including a deliberate testpackage violation.
- [Per-analyzer finding counts](./finding-counts.md) covers the baselines in this index.
