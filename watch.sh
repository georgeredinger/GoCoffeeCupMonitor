export GOROOT="./src"
inotifywait -m -r -e close_write src/  | while read line; do go test;done
