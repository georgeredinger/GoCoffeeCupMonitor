inotifywait -m -r -e close_write .  | while read line; do go test;done
