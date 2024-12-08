## Description
This directory contains the definitions for messages that are communicated between client and server.

## Compiling instructions
To generate the corresponding server (.go) files along with the corresponding client (.cs) files:
* Ensure that you have `protoc.exe` and that its location is in the `Path` system environment variable.
* Run the batch script `update_grpc.bat` located in this folder. 
* The outputted files are in the subfolders `\go` and `\cs` (the `go` files can be referenced directly within the rest of the project, whereas the `cs` files can be copied into the client codebase)