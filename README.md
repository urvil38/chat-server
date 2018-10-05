# Message Server

![chat-server](https://storage.googleapis.com/chat-bin/chat-server.png)

## Downloads

### Linux and Macos:
1. For Linux and macos user you download chat-server with `curl`:

    #### Linux:

    ```
    $ curl -LO https://storage.googleapis.com/chat-bin/v1.0/linux/server 
    $ curl -LO https://storage.googleapis.com/chat-bin/v1.0/linux/clinet
    ```

    #### MacOs:

    ```
    $ curl -LO https://storage.googleapis.com/chat-bin/v1.0/darwin/server
    $ curl -LO https://storage.googleapis.com/chat-bin/v1.0/darwin/clinet
    ```

2. Make binary executable :
    ```
    $ chmod +x ./server
    $ chmod +x ./client
    ```

In order to gernerate X509 self signed certificates run the `cert.sh` script as following:
```
$ ./cert.sh admin@gmail.com  # which creates certs folder in current directory
```

Now start server and client as following:
```
$ ./server -addr 127.0.0.1 -p 8080
$ ./client -addr 127.0.0.1 -p 8080
```

> You can configure address and port using `-addr` and `-p` command line flags.