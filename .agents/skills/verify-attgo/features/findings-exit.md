# Findings and the exit code

Findings from a loaded linter exit 1. A stock `golangci-lint` that has no plugin exits 3 and prints no findings. That exit is not a pass.

## Sub-features

- `findings-exit-1` reports a finding and exits 1.
- `stock-exit-3` runs the binary on `PATH` and exits 3 without a finding.

## How to get to it (user POV)

- Run `./custom-gcl run` in the consumer clone after Launch.
- Run `golangci-lint run` with the stock binary, which is what a shell does when the custom binary is not the one on `PATH`.

## Driving it with custom-gcl

Preconditions:

- Doctor passed, so `./custom-gcl` is the local plugin build. The findings half uses `./custom-gcl-pername` after that binary's own `=>` check in the per-name feature.
- `golangci-lint version` has no `-custom-gcl-` marker, and `command -v golangci-lint` is not a path under this clone or `$CHECKOUT`. Record that version line in `$EVIDENCE/vouch-stock-version.txt`.

- **Findings.** Scaffold and run the vouch per-name recipe from the per-name feature. Require exit code 1, a JSON report with at least one issue, and `FromLinter` values that include `attgo`. Save `$EVIDENCE/vouch-findings.exit` as `1` and keep the JSON report.
- **Stock binary.** From the vouch clone, run `golangci-lint run --config=.golangci.yml`. That config names `attgo` and does not name the sample plugins, so the failure is about `attgo` rather than whichever plugin the binary loads first. Capture stdout and stderr to `$EVIDENCE/vouch-stock.txt`. Require exit code 3. Require the output to say the attgo plugin was not found. Require no issue list and no `0 issues.` line. Record the exit as a failed load.
- **Proof.** The pair is the proof: `./custom-gcl-pername` exits 1 on the fixture, and the stock binary exits 3 on the consumer config. Either result alone is incomplete.

## Gotchas

- Exit 1 is the findings result. A wrapper that treats every non-zero exit as a crash will hide a real finding.
- Exit 3 prints an Error line and no summary. Exit 7 can print `0 issues.`. See SKILL.md.
- `only-new-issues`, `--new-from-rev`, and `--new-from-patch` hide findings outside the diff. This recipe does not pass them. A clean exit under those flags is not this proof.
- The consumer lint workflows set `only-new-issues: true`. A green workflow does not mean the tree has zero findings. The finding-counts feature measures the tree.
