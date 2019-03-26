#!/bin/bash

# Generate key and self-signed cert on one line
openssl req -new -newkey rsa:2048 -x509 -sha256 -days 3650 -nodes -out demo.crt -keyout demo.key