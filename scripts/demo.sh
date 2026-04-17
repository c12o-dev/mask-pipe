#!/usr/bin/env bash
# asciinema demo script for mask-pipe
# Usage: asciinema rec --command "bash scripts/demo.sh" demo.cast
#
# After recording:
#   asciinema upload demo.cast
#   # or convert to GIF:
#   agg demo.cast demo.gif

set -e

# Typing simulation
type_cmd() {
    echo ""
    echo -n "$ "
    for ((i=0; i<${#1}; i++)); do
        echo -n "${1:$i:1}"
        sleep 0.04
    done
    echo ""
    sleep 0.3
}

pause() { sleep "${1:-1.5}"; }

clear

echo "  mask-pipe — filter secrets from terminal output"
echo ""
pause 2

# Demo 1: Basic masking
type_cmd "echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe"
echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe
pause 2

# Demo 2: Multiple secrets
type_cmd "echo 'DB: postgres://admin:s3cretP4ss@db.example.com:5432/app' | mask-pipe"
echo 'DB: postgres://admin:s3cretP4ss@db.example.com:5432/app' | mask-pipe
pause 2

# Demo 3: GitHub token
type_cmd "echo 'GITHUB_TOKEN=ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234' | mask-pipe"
echo 'GITHUB_TOKEN=ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234' | mask-pipe
pause 2

# Demo 4: Clean text passes through
type_cmd "echo 'no secrets here, just normal output' | mask-pipe"
echo 'no secrets here, just normal output' | mask-pipe
pause 2

# Demo 5: Dry-run mode
type_cmd "echo 'AKIAIOSFODNN7EXAMPLE' | mask-pipe --dry-run --no-color"
echo 'AKIAIOSFODNN7EXAMPLE' | mask-pipe --dry-run --no-color
pause 2

# Demo 6: Doctor
type_cmd "mask-pipe doctor"
mask-pipe doctor
pause 2

echo ""
echo "  Install: brew install c12o-dev/tap/mask-pipe"
echo "  GitHub:  https://github.com/c12o-dev/mask-pipe"
echo ""
pause 3
