build:
    nix build .#spyglass

run:
    go run .

try:
    ghostty --title=Spyglass -e ./result/bin/spyglass
