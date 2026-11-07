# Disabled or absent linter

A rule the config turns off, and a linter the config does not name, produce no findings. The run still has to load and report a control finding. An empty report from a failed load is not this proof.

## Sub-features

- `disabled-rule` turns a rule off with its `enable_*` setting and reports nothing for it.
- `absent-linter` omits a linter from the config and reports nothing for it.

## How to get to it (user POV)

- Set `enable_capital_comment: false` in the attgo settings (go-eth2-client does this at the pin; vouch sets it true).
- Leave a linter out of `linters.enable` and out of `linters.settings.custom`.

## Driving it with custom-gcl

Preconditions:

- Doctor passed.
- The per-name feature's `./custom-gcl-pername` exists in this clone.
- The testdata attgo settings match the consumer `.golangci.yml`, as the per-name feature requires.

- **Disabled rule.** Scaffold `verifyattgo/project.go` and `.golangci.pername.yml` as the per-name feature does, using `testdata/project/go-eth2-client.yml`. Run the same `./custom-gcl-pername run` command against `./verifyattgo/`. Exit code is 1. The JSON report has at least one issue with `FromLinter` `attgo` (the raw-string control). No issue text starts with `attgo_capital_comment`. The vouch config is the opposite control: the same comment is reported there. Run that control when the claim is that the flag is what changed the result.
- **Absent linter.** Copy the vouch per-name config to `.golangci.absent.yml`. Remove `sample_first` from `linters.enable` and remove its `settings.custom` entry. Leave `attgo` and `sample_second`. Run `./custom-gcl-pername run` with that config on `./verifyattgo/`, same flags, JSON path `$EVIDENCE/vouch-absent.json`. Exit code is 1. The report has `attgo` and `sample_second` issues and zero `sample_first` issues. The text log warns that `sample_first` is an unknown nolint linter. That warning is expected.
- **Proof.** Save both JSON reports and both exit codes. The control finding is part of the proof. A report that is empty is a failure of this feature even when the exit code is 0.

## Gotchas

- go-eth2-client also sets `enable_no_pkg_logger` and `enable_struct_field_order` false. This recipe uses the capital-comment comment because the fixture contains it. Do not read a zero count of those other rules on the real tree as proof they are off. The tree may have no violation.
- An absent name does not appear in `linters --json` at all, enabled or disabled. A disabled attgo rule still leaves the linter name `attgo` enabled.
- Exit 3 with `plugin "attgo" not found` is a failed load. It is not an absent linter. The findings-exit feature covers that case.
