# public-repo-check

We'd like our public repositories to follow some basic standards:

- A licence (preferably MIT)
- A README.md file
- A CONTRIBUTING.md file

This tool makes sure each public repository has those (aside from archived projects and forks).

## Install

```
% go get github.com/dxw/public-repo-check
```

## Use

```
% public-repo-check dxw
✅ dxw/Edit-Flow: Fork. No further checks
❌ dxw/public-repo-check: License missing!
❌ dxw/public-repo-check: No README found
❌ dxw/public-repo-check: No CONTRIBUTING found
✅ dxw/wpc: License OK
✅ dxw/wpc: Has README
❌ dxw/wpc: No CONTRIBUTING found
```
