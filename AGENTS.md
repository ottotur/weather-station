## Git Workflow

- The `main` branch must always contain stable, working code.
- Never make code changes or commits directly on `main`.
- Before starting implementation, check the current Git status and branch.
- Create a new branch from `main` for each OpenSpec change:
  - `feature/<name>` for features
  - `fix/<name>` for bug fixes
  - `docs/<name>` for documentation-only changes
- Commit completed logical units of work to the current feature branch.
- Use conventional commit messages (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`).
- Include relevant OpenSpec artifacts in version control.
- Run applicable tests, linting and build checks before proposing a merge.
- Never commit secrets, credentials, `.env` files or generated dependencies.
- Never force-push or discard existing uncommitted work.
- Deliver changes to `main` through a GitHub pull request, not a local merge.
- When asked to create a pull request, push the feature branch, open or update
  the PR, and report the required CI checks and any failures.
- Merge the PR only after explicit user approval and successful required checks.
  After merging, update the local `main` from the remote.
- After completing a change, report:
  - Current branch
  - Commits created
  - Tests performed and their results
  - Whether the change is ready to merge into `main`

## Agent skills

### Issue tracker

For reading, creating or updating issues, use GitHub Issues.
Read `docs/agents/issue-tracker.md`.

### Triage labels

For issue classification and readiness labels, use the mappings in
`docs/agents/triage-labels.md`.

### Domain docs

Before codebase exploration or domain decisions, read the single-context
consumer rules in `docs/agents/domain.md`.
