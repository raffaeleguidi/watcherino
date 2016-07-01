# watcherino
A simple, selfcontained, tool to execute commands on filesystem changes


usage: watcherino [<flags>] <command> [<folder>]

    Flags:
        --help         Show context-sensitive help (also try --help-long and --help-man).
        --pattern="*"  Pattern filter
        --delay=5      Number seconds to wait from last change
        --version      Show application version.
    
    Args:
        <command>   Command to execute
        [<folder>]  Folder to watch for changes

example:
    
    ./watcherino ./example.sh --pattern="*.txt" &
      
Watches the current folder for new and modified .txt files and executes example.sh after 5 seconds since the last change event. The script receives three arguments: folder, file and CREATE|WRITE
