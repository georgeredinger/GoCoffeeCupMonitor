inotifywait -m -r -e close_write .  | while read line
 do 
   echo $line	| cut -d " " -f 3| grep --silent "^.*\.go" 
	 if [ "$?" -eq "0" ]; then
		 go test
	 fi

 done
