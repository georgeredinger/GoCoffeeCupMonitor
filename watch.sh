export GOPATH="/home/george/workspace/GoCoffeeCupMonitor"
inotifywait -m -r -e close_write src/  | while read line
do 
  go test xbee
done
