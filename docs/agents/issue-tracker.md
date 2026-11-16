# Issue tracker: GitHub

Issues live in GitHub Issues for `ottotur/weather-station`.
Use the `gh` CLI.

OpenSpec planning artifacts remain in `openspec/`.
Link related issues and OpenSpec changes rather than duplicating their specs.

## Operations

- Create: `gh issue create --repo ottotur/weather-station --title "..." --body-file <file>`
- Read: `gh issue view <number> --repo ottotur/weather-station --comments`
- Inspect metadata: `gh issue view <number> --repo ottotur/weather-station --json number,title,body,labels,comments`
- List: `gh issue list --repo ottotur/weather-station --state open --json number,title,body,labels`
- Comment: `gh issue comment <number> --repo ottotur/weather-station --body-file <file>`
- Label: `gh issue edit <number> --repo ottotur/weather-station --add-label "..." --remove-label "..."`
- Close: `gh issue close <number> --repo ottotur/weather-station --comment "..."`

When a skill says “publish to the issue tracker”, create a GitHub issue.
When it says “fetch the relevant ticket”, read the issue and its comments.

## Pull requests as a triage surface

**PRs as a request surface: no.**

GitHub issues and PRs share a number space. If an issue reference resolves
to a PR, inspect it using `gh pr view` and `gh pr diff`.
