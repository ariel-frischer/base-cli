# VHS tape file for base-cli demo GIF
Output assets/demo.gif

Require base-cli
Require tree

# Terminal styling
Set Shell "bash"
Set FontFamily "Noto Sans Mono"
Set FontSize 16
Set Width 1100
Set Height 1050
Set Padding 20
Set Framerate 30

# Window decoration
Set WindowBar Colorful
Set WindowBarSize 40
Set Margin 20
Set MarginFill "#282A36"
Set BorderRadius 10

# Theme
Set Theme "Dracula"

# Timing
Set TypingSpeed 40ms
Set PlaybackSpeed 1.5
Set CursorBlink false

# --- Hidden setup ---
Hide
Type "cd /tmp && rm -rf base-cli-demo && mkdir base-cli-demo && cd base-cli-demo"
Enter
Sleep 500ms
# Pre-configure so init doesn't prompt for module path
Type "base-cli config set host github.com && base-cli config set git_user demo"
Enter
Sleep 500ms
Type "clear"
Enter
Sleep 300ms
Show

# --- Demo starts ---

# 1. Init a new project (module auto-derived from config)
Type "base-cli init go-cli -d 'My Go CLI app'"
Sleep 2s
Enter
Sleep 5s

Ctrl+L

# 2. Show the generated structure
Type "tree go-cli"
Sleep 1s
Enter
Sleep 7s
