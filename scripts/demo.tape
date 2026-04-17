Output demo.gif
Output demo.webm

Set Shell "bash"
Set FontSize 16
Set Width 900
Set Height 500
Set Theme "Dracula"
Set Padding 20

Type "# mask-pipe — filter secrets from terminal output"
Enter
Sleep 1.5s

Type "echo 'AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE' | mask-pipe"
Enter
Sleep 2s

Type "echo 'postgres://admin:s3cretP4ss@db.example.com:5432/app' | mask-pipe"
Enter
Sleep 2s

Type "echo 'GITHUB_TOKEN=ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234' | mask-pipe"
Enter
Sleep 2s

Type "echo 'no secrets here' | mask-pipe"
Enter
Sleep 2s

Type "echo 'AKIAIOSFODNN7EXAMPLE' | mask-pipe --dry-run --no-color"
Enter
Sleep 2s

Type "mask-pipe doctor"
Enter
Sleep 3s

Type "# Install: brew install c12o-dev/tap/mask-pipe"
Enter
Sleep 3s
