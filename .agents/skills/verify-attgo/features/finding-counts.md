# Per-analyzer finding counts

One run of the attgo plugin reports every enabled rule under `FromLinter` `attgo`. The analyzer name is the prefix of `Text`, up to the first `:`. Compare those counts to the baseline table in the index.

## Sub-features

- `count-attgo` counts each attgo analyzer on a clean pin.
- `count-test` counts the four test-config linters on the go-eth2-client pin.

## How to get to it (user POV)

- Run `./custom-gcl run --enable-only=attgo` with the consumer's `.golangci.yml`, caps disabled.
- Run `./custom-gcl run --config=.golangci.tests.yml` for the test config.

## Driving it with custom-gcl

Preconditions:

- Doctor passed, including `git symbolic-ref refs/remotes/origin/HEAD` printing `refs/remotes/origin/master` and a successful `git diff origin/master...HEAD`.
- The clone is a clean pin. `verifyattgo/` and `foo/` are absent. A count that includes scaffolding is invalid.
- `GOTOOLCHAIN=go1.26.10`.

- **Production count.** From the clone root:

  ```
  ./custom-gcl run --enable-only=attgo \
    --output.json.path="$EVIDENCE/<consumer>-attgo.json" \
    --uniq-by-line=false --max-issues-per-linter=0 --max-same-issues=0 \
    --color never ./...
  ```

  Exit code is 0 when every baseline cell for that consumer is 0, and 1 when any cell is non-zero. Group issues by the analyzer prefix. Each prefix equals the baseline in the index. An issue whose `FromLinter` is `typecheck` fails the drive. Do not subtract it and continue.
- **Test-config count.** On go-eth2-client, use the clean test-lint command. Group by `FromLinter`. testpackage, testifylint, thelper, and usetesting are each 0. On vouch, record the missing config and stop, as the test-lint feature does.
- **Proof.** Save the JSON and a one-line count per analyzer next to it in `$EVIDENCE/<consumer>-counts.txt`. The text summary `N issues` is not the proof. The grouped JSON is.

## Gotchas

- The default cap of 3 identical texts would report 3 `attgo_capital_comment` findings on the vouch pin, not 9. The flags above are part of the command.
- `attgo_current_year` is 0 when the diff has no production Go file, including an empty diff. See the index. A shallow clone, a failed diff, or a saved line `attgo_current_year: warning` is not this baseline.
- `--enable-only=attgo` keeps the consumer's attgo settings. It does not turn on rules that consumer sets false. The index marks those rules `off`. Do not record them as 0.
- A cached `custom-gcl` from a consumer workflow is the wrong binary. See SKILL.md Gotchas. Doctor's `(devel)` replace path is the check.
- A run whose working directory is not the clone can exit 5 or 7 and match no baseline. See SKILL.md.
- `GOLANGCI_LINT_CACHE` is the run cache from Launch. A replay is not a count. See SKILL.md Gotchas.
