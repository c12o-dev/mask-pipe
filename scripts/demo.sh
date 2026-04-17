#!/usr/bin/env bash
# asciinema / VHS demo script for mask-pipe
# Each command clears the screen and shows the header, so the recording
# stays compact. The final frame shows doctor + install info.
#
# Usage: asciinema rec --command "bash scripts/demo.sh" demo.cast
#   or:  just run in Windows Terminal and record with Win+G

set -e

HEADER="  mask-pipe — filter secrets from terminal output"
FOOTER_INSTALL="  Install: brew install c12o-dev/tap/mask-pipe"
FOOTER_GITHUB="  GitHub:  https://github.com/c12o-dev/mask-pipe"

type_cmd() {
    echo -n "$ "
    for ((i=0; i<${#1}; i++)); do
        echo -n "${1:$i:1}"
        sleep 0.04
    done
    echo ""
    sleep 0.3
}

frame() {
    clear
    echo "$HEADER"
    echo ""
}

pause() { sleep "${1:-2}"; }

# Frame 1: AWS key
frame
type_cmd "echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe"
echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe
pause

# Frame 2: DB URL
frame
type_cmd "echo 'DB: postgres://admin:s3cretP4ss@db.example.com:5432/app' | mask-pipe"
echo 'DB: postgres://admin:s3cretP4ss@db.example.com:5432/app' | mask-pipe
pause

# Frame 3: GitHub token
frame
type_cmd "echo 'GITHUB_TOKEN=ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234' | mask-pipe"
echo 'GITHUB_TOKEN=ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234' | mask-pipe
pause

# Frame 4: Clean passthrough
frame
type_cmd "echo 'no secrets here, just normal output' | mask-pipe"
echo 'no secrets here, just normal output' | mask-pipe
pause

# Frame 5: Dry-run
frame
type_cmd "echo 'AKIAIOSFODNN7EXAMPLE' | mask-pipe --dry-run --no-color"
echo 'AKIAIOSFODNN7EXAMPLE' | mask-pipe --dry-run --no-color
pause

# Frame 6 (final): Doctor + install info
frame
type_cmd "mask-pipe doctor"
mask-pipe doctor
echo ""
echo "$FOOTER_INSTALL"
echo "$FOOTER_GITHUB"
echo ""
pause 4
