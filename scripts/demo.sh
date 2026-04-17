#!/usr/bin/env bash
# Demo script for mask-pipe (~15 seconds)
# Records well with: Win+G (Game Bar), asciinema, or VHS
#
# Usage: bash scripts/demo.sh

set -e

type_cmd() {
    echo -n "$ "
    for ((i=0; i<${#1}; i++)); do
        echo -n "${1:$i:1}"
        sleep 0.04
    done
    echo ""
    sleep 0.3
}

clear
sleep 1

# 1. AWS key — "検出して隠す"
type_cmd "echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe"
echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe
echo ""
sleep 2

# 2. DB URL — "パスワードだけ隠す"
type_cmd "echo 'postgres://admin:s3cretP4ss@db.example.com/app' | mask-pipe"
echo 'postgres://admin:s3cretP4ss@db.example.com/app' | mask-pipe
echo ""
sleep 2

# 3. Clean text — "壊さない"
type_cmd "echo 'no secrets here, just normal output' | mask-pipe"
echo 'no secrets here, just normal output' | mask-pipe
echo ""
sleep 2

# 4. Version — "軽い"
type_cmd "mask-pipe --version"
mask-pipe --version
echo ""
sleep 3
