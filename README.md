# To Run all of the necessary programs, make sure that your Linux Machines all have Go, C, Node.js, and NPM installed on them.

# To Run the Program:
First, cd into the folder "backend" inside the LatencyDetector root folder: LatencyDetector/Backend

Then run the the command "node server.js", which will start the Latency Detector Master and begin listening for nodes.

Next, you can navigate to the "index.html" file in the same folder within the editor, and right click it to "copy full path" which you can paste into your browser to see the web page.

In each Linux machine running the program, configure the following to turn off the necessity for "sudo" to prompt user entered password:
                                Open a Terminal window and type:

                                sudo visudo
                                In the bottom of the file, add the following line:

                                $USER ALL=(ALL) NOPASSWD: ALL
                                Where $USER is your username on your system. Save and close the sudoers file (if you haven't changed your default terminal editor (you'll know if you have), press Ctl + x to exit nano and it'll prompt you to save).

Now to start a basic Node that you can see automatically generated on the webpage (see comments in the server.go, and server.js files for anywhere that IP address changes are necessary) open another terminal window and cd into the eBPF folder to run the commmand "go run server.go"

This should start a node on the webpage that you can click on to input the different commands for eBPF executions.

To configure the system for multiple nodes, simply change the VM setting to Bridged adapter to give each VM a unique IP address, and make the changes mentioned above in those files.

