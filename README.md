## This repository contains the server code for Paper Volleyball which can be deployed publically. 
It is currently in development and made publicly accessible for a limited time.

## Deployming instructions
* Ensure all unit tests run successfully: `go test .\...`
* Build command (run via Terminal): `go build`
* Test on local machine by simply running the built `.exe` to start the server on console.
* For servers hosted on Windows, connect to the RDP instance, copy the built `.exe` into the server, and run it.
* When hosting via cloud services, ensure that the Windows Firewall setting on the instance is set to allow TCP on the specified port number in this file, and also port 80 to allow for WebSocket connections.
* The url for a request will be <http or ws>://<serverAddress>:<port>/<command>. If running locally, the value of <serverAddress> is `localhost`. If deploying on the cloud, then it is the public IP address of the instance.

## License and Copyright
* This repository and its contents are © 2024 Terence Ma. All rights reserved.
* Unauthorized copying, distribution, or modification of any part of this repository, via any medium, is strictly prohibited without the express written permission of the authors.
* For inquiries regarding usage rights or permissions, please contact `papervolleyballdev@gmail.com`.