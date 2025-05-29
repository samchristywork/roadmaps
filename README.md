![Banner](https://s-christy.com/sbs/status-banner.svg?icon=maps/route&hue=150&title=Roadmaps&description=Interactive%20tree%20editor%20for%20planning%20and%20progress%20tracking)

## Overview

Roadmaps is a browser-based interactive tree editor for building and tracking
hierarchical plans. A Go backend serves a D3.js frontend where trees are edited
entirely with keyboard shortcuts. Each roadmap is stored as a JSON file on the
server and can be exported as plain text, Markdown, or JSON.

The editor uses vim-style navigation (`hjkl`) to move between nodes and exposes
a full set of editing commands for inserting, renaming, deleting, duplicating,
and reordering nodes. Nodes can be marked done and a live progress counter
shows how many items are complete across the whole tree.

<p align="center">
  <img src="./assets/screenshot.png" />
</p>

## Features

- Interactive tree visualization powered by D3.js
- Vim-style keyboard navigation (`hjkl`)
- Node operations: insert child, insert sibling, rename, duplicate, delete, promote children
- Cut-and-paste to move nodes anywhere in the tree
- Mark nodes done with progress tracking across the whole tree
- Per-node notes with optional inline display
- Collapse and expand subtrees; collapse/expand all
- Incremental search with match cycling (`/`, `n`, `N`)
- Undo and redo (up to 100 levels)
- Sort children alphabetically
- Import from indented plain text or JSON; upload from file
- Export as plain text, Markdown checklist, or JSON
- Dark mode toggle
- Breadcrumb navigation bar and children sidebar panel
- Multiple roadmaps with a management page for create, rename, and delete

## Keybindings

## Usage

## Dependencies

## License

This work is licensed under the GNU General Public License version 3 (GPLv3).

[<img src="https://s-christy.com/status-banner-service/GPLv3_Logo.svg" width="150" />](https://www.gnu.org/licenses/gpl-3.0.en.html)
