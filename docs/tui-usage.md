# TUI Usage

## Main View

```
+------------------------------------------------------------------+
| Peachy - Theme Creator                    wallpaper.png [Normal] |
+--------------------------------+---------------------------------+
|  COLORS                        |  PREVIEW                        |
| ---------------------------    | ------------------------------- |
|   0 #### #1A1A2E  black        |  Source: wallpaper.png          |
|   1 #### #F54242  red          |  ################################|
| > 2 #### #42F554  green        |  ################################|
|   3 #### #F5F542  yellow       |                                 |
|   4 #### #4287F5  blue         |  Normal: ################       |
|   5 #### #F542F5  magenta      |  Bright: ################       |
|   6 #### #42F5F5  cyan         |                                 |
|   7 #### #E0E0E0  white        |  Sample:                        |
|   ...                          |  +---------------------------+  |
|                                |  | Hello, World!             |  |
|                                |  | Error Success Warning     |  |
|                                |  +---------------------------+  |
+--------------------------------+---------------------------------+
| j/k:nav  Enter:edit  o:open  e:extract  m:mode  s:save  ?:help   |
+------------------------------------------------------------------+
```

## Keyboard Shortcuts

### Main View

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Edit selected color |
| `o` | Open file picker |
| `e` | Re-extract colors from image |
| `s` | Save theme (prompts for name) |
| `l` | Browse and load themes |
| `a` | Apply current palette to system |
| `r` | Reset to extracted colors |
| `m` | Cycle extraction mode |
| `M` | Cycle extraction mode (reverse) |
| `t` | Toggle light/dark mode |
| `?` | Show help |
| `q` | Quit |

### Theme Browser

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate themes |
| `Enter` | Load theme (populate colors) |
| `a` | Apply theme (save as active) |
| `Esc` | Cancel |

### File Picker

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate |
| `Enter` | Select file / Enter directory |
| `Backspace` | Go to parent directory |
| `/` | Start search |
| `w` | Jump to ~/Wallpapers |
| `~` | Jump to home directory |
| `p` | Toggle image preview |
| `Esc` | Cancel |

### Search Mode (after pressing `/`)

| Key | Action |
|-----|--------|
| Type | Filter files by name |
| `Enter` | Select file |
| `Esc` | Cancel search |
| `Up/Down` | Navigate results |

### Color Editor

| Key | Action |
|-----|--------|
| `j` / `k` | Select field (Hue/Saturation/Lightness) |
| `h` / `l` | Adjust value (small step) |
| `H` / `L` | Adjust value (large step) |
| `#` | Enter hex mode |
| `u` | Reset to original |
| `Enter` | Confirm changes |
| `Esc` | Cancel |
