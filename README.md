# watcherino
A simple selfcontained tool to execute commands on filesystem changes

example:
   ./watcherino . ./example.sh 10 &
   
   watches the current folder and executes example.sh after 10 seconds since the last change event
   the script receives three arguments: folder, file and CREATE|WRITE
