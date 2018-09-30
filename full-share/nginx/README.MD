To generate self-signed certs for ssl:
```
openssl req -nodes -x509 -newkey rsa:4096 -keyout privkey.pem -out fullchain.pem -days 365
```