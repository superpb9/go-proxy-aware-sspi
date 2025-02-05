# go-proxy-aware-sspi

A simple Go module for transparent Windows proxy authentication calling SSPI, similar to curl's "-U :" feature.

## Installation
```
go get github.com/superpb9/go-proxy-aware-sspi
```
## Features

- Transparent Windows proxy authentication using SSPI
- No explicit credentials required (uses current Windows user context)
- Similar to curl's "-U :" functionality
- Clean and simple API


## Usage Example
```
package main

import (
    "log"
    "github.com/superpb9/go-proxy-aware-sspi/wsconnect"
)

func main() {
    websocketURL := "wss://example.com/websocket"
    proxyURL := "http://proxy.example.com:8080"
    
    // Connect through proxy using current Windows user credentials
    conn, err := wsconnect.Connect(websocketURL, proxyURL)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Use websocket connection...
}

```
## Requirements
- Windows operating system
- A proxy with NTLM authentication support
- Go 1.21 or later

## Dependencies
- github.com/alexbrainman/sspi/negotiate
- github.com/gorilla/websocket

## License
MIT License