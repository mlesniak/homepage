# Automatic git synchronization in Visual Studio Code

I use [Visual Studio Code](https://code.visualstudio.com) in combination with different Markdown plugins to implement a basic version of a [Zettelkasten](https://en.wikipedia.org/wiki/Zettelkasten)-system. Hence, backing up (and synchronizing) my data becomes more and more important for me. For different reasons I am not able to use Dropbox, but a private git repository on my own server suffices as well for this particular use-case. 

Visual Studio Code allows to define [build tasks](https://code.visualstudio.com/docs/editor/tasks) which can also run in the background. While these are normally used for watcher tasks in frontend build tools like npm, arbitrary scripts can be executed. The following script will periodically commit all new files and changes to existing ones using a [timestamp](https://en.wikipedia.org/wiki/Unix_time) and push them to the git origin.

    #!/bin/bash
    # File named .update.sh
    
    echo "Starting automatic git push every 60 seconds"
    while true
    do
        clear
        echo "--- $(date) --------------------------------------------------------------"
        git add . && git commit -m $(date +%s) && git push -u origin

        sleep 60
    done

By configuring VSCode's `task.json` to hide the command output and run in the background, continous synchronisation is achieved:

    {
        "version": "2.0.0",
        "tasks": [
            {
                "label": "push",
                "type": "shell",
                "isBackground": true,
                "command": ".update.sh",
                "problemMatcher": [],
                "group": {
                    "kind": "build",
                    "isDefault": true
                },
                "presentation": {
                    "echo": true,
                    "reveal": "never",
                    "focus": false,
                    "panel": "shared",
                    "showReuseMessage": false,
                    "clear": true
                },
                "runOptions": {
                    "runOn": "folderOpen"
                }
            }
        ]
    }

Usually, you have to remember to start the build task the first time you open the specific folder. By setting the `runOn` option *and* allowing automatic tasks for this specific folder (`CMD-P`, then _Tasks: Manage automatic tasks in folder_), this problem is also solved. 