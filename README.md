# godolta
[![Build Status](https://travis-ci.org/mtib/godolta.svg?branch=master)](https://travis-ci.org/mtib/godolta)

go port of [delta-l](https://github.com/LFalch/delta-l) written in Rust by [LFalch](https://github.com/LFalch)

## Usage
    godolta [-p password] [-f] [-o file] [-c=false] <e[ncrypt] | d[ecrypt]> <file>

The Golang port of original [delta-l](https://github.com/LFalch/delta-l) in Rust averages about 1/2x Speed.
This implementation allows usage of deltal/Decoder and deltal/Encoder as well as other functions to use it as a library. Both types implement [io.ReadSeeker](https://golang.org/pkg/io/#ReadSeeker).

This shouldn't be used for heavy cryptographic use, but allows a lightweight and fast encryption algorithm for use in other applications.
## Algorithm
### Encryption
```go
// pseudo-code:
b = []byte(inputFile)
p = siphash(password)
c = []byte
c += "ΔL" // if "Δl", skip next line
c += siphash(inputFile)
c += b[0] + p[0]
for i=1; i < len(b); i++ {
    c += (b[i] + b[i-1] + p[i%8]) % 256
}
save(c)
```
### Decryption
Does the exact reverse of the encryption algorithm above, but makes sure the checksums match, if ```c[2]=='L'```.

## Installation
```
go get -u github.com/mtib/godolta
```
Documentation: [godoc.org](https://godoc.org/github.com/mtib/godolta/deltal)

# License
The MIT License (MIT)

Copyright (c) 2016 mtib

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
