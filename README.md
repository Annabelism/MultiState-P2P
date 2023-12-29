# MultiState-P2P
CS 242 Fall 2023 Final Project


HOW TO USE:
We first run the main.go file to establish the first node in the network Peer1. Upon running the file, Peer1 will receive a message in the terminal saying “error: no connection is made, peerIP = node.IP”. This is because it’s the first node in this network and therefore no connection is made. 
Since we are simulating all the peers on one machine when testing, we have to hardcode the port number each peer is listening on so that there are no conflicts. In our test, we set the port of Peer1 to be 8888, and the port of Peer2 to be 8887. This won’t be needed if testing on separate machines.


Then, we establish the subsequent nodes and test the basic functionality:
1. Joining the network as a new member
2. requesting a file from the ORIGIN
3. sharing a file with the ORIGIN (receiving a download request from ORIGIN)
4. Leaving the network


We will test this by having peer A join the network, and then having another peer
B join after peer A has joined. Then, we will make ORIGIN leave the network.


Functionality and Usage Detailed
1. Connecting to the Network
On startup, the application attempts to dial a predefined network node (localhost:8888 in our test case).
A response, indicating a successful connection, is displayed (received and printed in terminal: CONNECTED TO NETWORK Access granted).


Main Commands
After connecting to the network, the application presents the user with the following options:
download(1): Request to download a file.
update(2): Update the file listing (add or delete files).
leave(3): Leave the network.
cancel request(x): Cancel the current request.


Update Action
When choosing to update (2), the user is prompted to specify the action:
add: Add a new file entry.
delete: Remove an existing file entry.
cancel request(x): Cancel the update request.


Update File Name
Upon entering either `add` or `delete` in the terminal, the user is prompted to enter the file name they want to update (eg. test.txt)
Then, the listening peer will receive an updated table (printed in terminal: Received tableH update: and then the updated tableH)


Download Action
When choosing to download(1), the application prepares the local environment to manage these file entries by adding a directory named “file_received”. Then the application will create or truncate a file with the specified input_file name within this directory. If the file creation is successful, a confirmation message (FILE CREATED) is printed to the console.


Adding or Deleting a File
The user is prompted to enter the file name for the add or delete action.
If an invalid input is provided or the entered file name does not exist, the application displays Invalid input. Please try again.


TCP Request Handling
The application receives and processes TCP requests, displaying messages such as TCP request received, indicating the type of request and the current state of the file table (tableH).


Error Handling
The application handles errors gracefully, providing feedback such as Invalid input. Please try again. in case of user input errors or issues with file operations.
Interrupt Signal Handling
The application can be stopped using an interrupt signal (e.g., Ctrl+C), which safely terminates the application.


Notes
Ensure that user input is clear. It might take some time to process the user request, so it is advised that users wait for the response message instead of trying too many times after receiving no response.
The application is interactive and requires user input for various operations.


