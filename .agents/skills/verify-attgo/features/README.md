# attgo verification map

This directory is the maintained source for verifying the custom binary a consumer runs. Read this index before driving, then use the matching feature file as the recipe.

## Baseline preconditions

- Launch both clones as `SKILL.md` describes, from a clean pin, and pass doctor before a drive.
- Pins:

  | Consumer | Pin | golangci-lint release in `.custom-gcl.yml` |
  | --- | --- | --- |
  | vouch | `83758284525e3abb5bf46a1680a944020dbf831f` | v2.12.2 |
  | go-eth2-client | `fb4a89371097364564680a54fe0ca0ebcc57d7c1` | v2.13.1 |

- The go-eth2-client pin is the commit that contains `.golangci.tests.yml`. It is not the default branch. When that file is on the default branch, move this pin to the default-branch commit that contains it, remeasure the baselines, and replace the numbers below. Until then every go-eth2-client drive uses this pin.
- Use the toolchain Launch exports (`GOTOOLCHAIN=go1.26.10`) for every `custom-gcl` run. Go 1.27 makes this golangci-lint release report export-data errors and exit 1 with no attgo findings.
- Run from the clone root. Pass `--uniq-by-line=false --max-issues-per-linter=0 --max-same-issues=0` whenever the recipe counts findings. The defaults stop at 50 per linter and 3 per identical text, which hides findings.
- The vouch pin has no `.golangci.tests.yml`. A test-lint drive on vouch reports that unmet precondition and stops. Do not add the file.

## Baselines

Measured with the plugin at `0b7bfd3e5d282aee95632c4bc674fe2e75625ed0`, on a clean pin, `GOTOOLCHAIN=go1.26.10`, `--enable-only=attgo`, caps disabled. Split each issue `Text` on the first `:` to get the analyzer name. `FromLinter` is `attgo` for all eight.

vouch `83758284525e3abb5bf46a1680a944020dbf831f`:

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

go-eth2-client `fb4a89371097364564680a54fe0ca0ebcc57d7c1`:

| Analyzer or linter | Findings |
| --- | --- |
| attgo_enum_iota, attgo_current_year, attgo_func_opts, attgo_raw_string, attgo_interface_check | 0, rule on |
| attgo_no_pkg_logger, attgo_capital_comment, attgo_struct_field_order | off in this config, not a measured zero |
| testpackage | 0 |
| testifylint | 0 |
| thelper | 0 |
| usetesting | 0 |

`attgo_current_year` is 0 on this pin because the diff against `origin/master` has no production Go files. The changed files are the test config, the workflow, and `_test.go` files, and the production config does not lint tests.

A mismatch fails the feature. Replace a number only after a clean run at the pin produces a different count, and record the plugin commit next to the new number. A count taken while `verifyattgo/` exists in the clone is not a baseline.

`attgo_current_year` sees only files in `git diff origin/master...HEAD`. On the vouch pin that diff is empty, so 0 is the expected count. It does not mean the analyzer is off.

## Driving conventions

- Start from the clean pin unless the feature says to add scaffolding.
- Copy scaffolding from this checkout's `testdata/project/`. Remove it by deleting the clone, not by editing the consumer.
- Treat exit 1 as findings, not as a crashed command. Treat exit 0 as a loaded linter with an empty report. Treat exit 3 and exit 7 as a failed run.
- Record the feature id, the consumer, the command, the exit code, and the JSON report with every artifact.
- An unreachable path is reported with the command and the unmet precondition. A skipped consumer is not verified through the other one.

## Features

- [Per-name registration and nolint](./per-name-nolint.md) covers a directive that names one linter, a directive that names the attgo plugin, and a directive that names a rule instead of a linter.
- [Disabled or absent linter](./disabled-or-absent.md) covers a rule the config turns off, and a linter the config does not name.
- [Findings and the exit code](./findings-exit.md) covers a finding exiting 1, and a stock binary exiting 3 with no findings.
- [Test-lint config](./test-lint.md) covers `.golangci.tests.yml` on the go-eth2-client pin, including a deliberate testpackage violation.
- [Per-analyzer finding counts](./finding-counts.md) covers the baselines in this index.
