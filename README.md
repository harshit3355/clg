# clg

`clg` is a small command-line tool for maintaining a Markdown `CHANGELOG.md`.
Instead of editing the changelog during development, record each change as a
YAML file and turn all unreleased entries into a dated release when you publish.

> **Warning:** This project is work in progress. Use at your own risk. APIs might change.
> The docs are 100% AI-generated. No releases published yet.

## How it works

`clg` uses this layout in the current working directory:

```text
.
├── CHANGELOG.md
├── .clg.yml                  # optional project configuration
└── changelogs/
    └── unreleased/
        └── added-0199321f-7b2c-7c4f-bd12-4c5f8f7c2a10.yml
```

Personal configuration can also be stored in `~/.clg.yml` in your home
directory; project configuration overrides matching settings from that file.

The files in `changelogs/unreleased/` are temporary release notes. `clg release`
groups them by type—or by group and then type when groups are configured—inserts
the resulting Markdown into `CHANGELOG.md`, and removes the source files.

## Installation

With Go installed:

```sh
go install github.com/hettiger/clg@latest
```

Or build the binary from a checkout:

```sh
go build -o clg .
```

`clg` requires Go 1.25.8 or newer.

## Quick start

Add the insertion marker to `CHANGELOG.md` once, usually near the top. The
default marker is `<!-- CLG -->`:

```md
<!-- CLG -->
```

Optionally configure groups, custom labels, or a different marker in `.clg.yml`:

```yaml
marker: "<!-- CLG -->"
groups:
  front: Frontend
  back: Backend
types:
  added: New Feature
  fixed: Bug Fix
```

When `groups` is configured, every entry must specify one of its keys and
releases are rendered with group headings containing type headings.

Optionally set your personal author information in `~/.clg.yml` (not in the
version-controlled project configuration):

```yaml
author:
  name: Jane Doe
  url: https://example.com/jane
```

This provides your author information across projects. Without a configured
author, `clg new` tries Git's `user.name` by default. Set `author.gitFallback` to
`false` to disable that fallback (see [Configuration](#configuration)).

Record a change. With no flags, `clg new` asks for the configured group (if any),
type, and message:

```sh
clg new
```

For scripts or a faster workflow, provide the type, message, and group (when
configured) directly:

```sh
clg new --type added --message "Support exporting reports"
clg new -g back -t fixed -m "Prevent duplicate notifications"
```

Review the unreleased entries:

```sh
clg show
```

To show only entries recorded on a specific Git branch:

```sh
clg show --branch feature/report-export
```

When you are ready to publish, pass the release tag:

```sh
clg release v1.2.0
```

This adds a section like the following immediately after `<!-- CLG -->`:

```md
## [v1.2.0] - 2026-09-06

### New Feature (1 change)

- Support exporting reports ([Jane Doe](https://example.com/jane))

### Bug Fix (1 change)

- Prevent duplicate notifications ([Jane Doe](https://example.com/jane))
```

With groups configured, the release uses one additional heading level:

```md
### Backend

#### Bug Fix (1 change)

- Prevent duplicate notifications ([Jane Doe](https://example.com/jane))
```

The release date is the current UTC date.

## Commands

### `clg new`

Create an unreleased changelog entry in `changelogs/unreleased/`.

```sh
clg new [flags]
```

| Flag | Description |
| --- | --- |
| `-g, --group` | Configured group key. If omitted, choose from an interactive list when groups are configured. |
| `-t, --type` | Configured change type. If omitted, choose from an interactive list. |
| `-m, --message` | Entry text. If omitted, enter it interactively. |
| `-a, --author` | Author name. Overrides the complete author from configuration or Git. |
| `-u, --url` | Optional author URL (HTTP or HTTPS). Requires a non-empty `--author`. |

Supplying the type, message, and group (when configured) makes the command
non-interactive. Author information does not require an interactive prompt.
`clg new` records the current Git branch in the entry, so it must be run from
a Git working tree. The generated filename contains the type and a UUIDv7, for
example `fixed-0199321f-7b2c-7c4f-bd12-4c5f8f7c2a10.yml`. When groups are
configured, the group key is prefixed to the filename, for example
`back-fixed-0199321f-7b2c-7c4f-bd12-4c5f8f7c2a10.yml`.

#### Author information

`clg new` checks author sources in this order:

1. `--author` and `--url` flags.
2. The `author` section in the merged configuration from `~/.clg.yml` and the
   project's `.clg.yml`.
3. Git's `user.name`, with no URL, unless `author.gitFallback` is `false`.

The first non-empty, valid source supplies the complete author. Fields are not
merged across sources: `--author "Jane Doe"` does not inherit a URL from the
configuration. The two configuration files are merged before author resolution,
so they form a single source; matching project settings override global settings
field by field. Keep personal author information in `~/.clg.yml` to avoid
assigning one contributor's identity to everyone using the project.

To supply both fields explicitly:

```sh
clg new -t added -m "Support exporting reports" \
  -a "Jane Doe" -u "https://example.com/jane"
```

Add `--group` when groups are configured.

Names and URLs are trimmed before validation. A source with both fields empty
is skipped, including flags explicitly set to empty strings. A non-empty URL
without a name, or an invalid URL, causes an error before any interactive
questions; the command does not fall back to another source. Sources after the
first valid author are not consulted.

Git author lookup is best-effort: lookup errors are treated as missing author
information. Set `author.gitFallback: false` in configuration to skip this
lookup. If no source supplies an author, including when Git fallback is disabled,
`clg new` prints a warning and continues without an author. It does not ask for
confirmation.

### `clg show`

Display all valid entries that have not yet been released:

```sh
clg show
```

The output includes the type, title, author name, and the Git branch associated
with each entry. The author column is empty for entries without an author.
Use `--branch` (or `-b`) to filter entries by branch. Group values are
shown when groups are configured and are used when generating a release. If
there are no matching entries, `clg` reports that there is nothing to show.

| Flag | Description |
| --- | --- |
| `-b, --branch` | Show only entries recorded on the specified Git branch. |

### `clg release [tag]`

Convert all unreleased entries into a release and insert it into
`CHANGELOG.md`:

```sh
clg release v1.2.0
```

| Flag | Default | Description |
| --- | --- | --- |
| `-m, --marker` | configured marker or `<!-- CLG -->` | Text where the new release is inserted. |

The marker must already exist in `CHANGELOG.md`. To use a different marker:

```sh
clg release v1.2.0 --marker "<!-- RELEASES -->"
```

Each release entry includes its recorded author name in parentheses when
available. If the entry also has an author URL, the name becomes a Markdown
link. Entries without an author name have no attribution suffix:

```md
- Support exporting reports ([Jane Doe](https://example.com/jane))
- Prevent duplicate notifications (Jane Doe)
- Handle empty report filters
```

This applies to both grouped and ungrouped releases. Author information comes
from the entry files, not the configuration or Git identity at release time.

If there are no unreleased entries, the command leaves the changelog unchanged.

### `clg clean`

Delete all unreleased entry files:

```sh
clg clean
```

The command asks for confirmation. Use `--force` when confirmation is not
possible or desired:

```sh
clg clean --force
```

This only removes files in `changelogs/unreleased/`; it does not modify
`CHANGELOG.md`.

## Change types

The following type keywords are supported:

| Keyword | Heading |
| --- | --- |
| `added` | New Feature |
| `fixed` | Bug Fix |
| `hotfix` | Hotfix |
| `changed` | Feature Change |
| `deprecated` | New Deprecation |
| `removed` | Feature Removal |
| `security` | Security Fix |
| `performance` | Performance Improvement |
| `other` | Other |

The heading is used when `clg release` groups entries.

## Configuration

Configuration files are loaded in this order, on top of the built-in defaults:

1. `~/.clg.yml` in your home directory: global, personal settings across projects.
2. `.clg.yml` in the current working directory: project-specific settings.

Both files are optional. Project settings override matching global settings;
settings not specified by the project retain their global or built-in values.
Nested mappings are merged field by field, rather than replaced as a whole.

Keep shared settings such as groups, change types, and the insertion marker in
the project's `.clg.yml`, which can be committed to version control. Keep personal
settings, especially your author name and URL, in `~/.clg.yml` instead.

The built-in marker and change types are:

```yaml
marker: "<!-- CLG -->"
types:
  added: New Feature
  fixed: Bug Fix
  hotfix: Hotfix
  changed: Feature Change
  deprecated: New Deprecation
  removed: Feature Removal
  security: Security Fix
  performance: Performance Improvement
  other: Other
```

Use `groups` to enable grouped releases. Group keys are used in entry files and
CLI flags; their values are the headings shown in generated Markdown:

```yaml
groups:
  front: Frontend
  back: Backend
```

In **`~/.clg.yml`**, use `author` to configure your identity for `clg new` when
no author is supplied through flags:

```yaml
author:
  name: Jane Doe
  url: https://example.com/jane
```

The URL is optional; when present, it must be a valid HTTP or HTTPS URL and have
an accompanying name. Avoid putting personal author information in the project's
`.clg.yml`: it would override contributors' global settings. For example, a
project setting only `author.name` would retain `author.url` from the global
file, potentially combining different identities.

If neither flags nor the merged configuration supplies an author name or URL,
`clg new` falls back to Git's `user.name` by default. To disable Git author lookup,
set the following in `~/.clg.yml` for a personal preference or in `.clg.yml` for a
project-wide setting:

```yaml
author:
  gitFallback: false
```

`author.gitFallback` defaults to `true` and follows the same configuration merge
rules: a project setting overrides the global setting. Disabling it does not
suppress an author supplied through flags or configuration. If no author is
available, `clg new` warns and creates the entry without one. The command still
requires a Git working tree to record the current branch.

Omitting `author` from the project file does not disable an author configured
in `~/.clg.yml`.

## Entry format

Each entry is a YAML document. `clg new` writes the `title`, `type`, and current
Git `branch`, plus `group` when groups are configured and `author` when available:

```yaml
group: back
type: changed
title: Improve report permissions
author:
  name: Jane Doe
  url: https://example.com/jane
branch: feature/report-export
```

The `author` field is an optional mapping with `name` and an optional `url`.
`clg new` omits `url` when no URL is supplied and omits the entire `author` field
when no author is available. The `branch` field is populated automatically from
the current Git branch and is used by
`clg show --branch`. Every YAML file in `changelogs/unreleased/` must have a
non-empty `title`, a configured `type`, and—when groups are configured—a
configured `group`. Invalid files prevent commands that read unreleased entries
from completing.

## Typical release workflow

```sh
# During development
clg new -g back -t added -m "Add CSV export"
clg new -g front -t fixed -m "Handle empty report filters"

# Review all entries, or only entries from one branch
clg show
clg show -b feature/report-export

# Before publishing
clg release v1.2.0
git diff -- CHANGELOG.md
git add CHANGELOG.md
git commit -m "Release v1.2.0"
```

Run `clg --help` or `clg <command> --help` for the command-line help.

## Development

Run the test suite with:

```sh
go test ./...
```

## License

See [LICENSE](LICENSE).
