### Steps for starting project
1) git clone https://github.com/smirnoffV/482.solution-storage.git
2) install package manager "dep" to machine
3) set GOPATH
4) in the project folder run a command "make install && make build"
5) set environment variables SERVICE_HOST and SERVICE_PORT
6) set environment variables PARENT_NODE_SERVICE_HOST and PARENT_NODE_SERVICE_PORT (if you start first instance of node, leave the values blank.)
7) to start server run a command "./bin/run"