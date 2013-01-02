inotifywait -m -r -e close_write .  | while read line
 do 
   echo $line	
   echo $line	| cut -d " " -f 3| grep "^.*\.go" 
	 if [ "$?" -eq "0" ]; then
		 go test
	 else
     echo "no go"
	 fi

 done
