# Simple Nginx OTP
A simple Nginx OTP module for use with `auth_request`

[**Docker Hub Image »**](https://hub.docker.com/r/jarylc/simple-nginx-otp)

[**Explore the docs »**](https://gitlab.com/jarylc/simple-nginx-otp)

[Report Bugs](https://gitlab.com/jarylc/simple-nginx-otp/-/issues/new?issuable_template=Bug)
· [Request Features](https://gitlab.com/jarylc/simple-nginx-otp/-/issues/new?issuable_template=Feature%20Request)


## About
### Features
- Lightweight and fast
- Returns a very simple form with a text field and submit button for OTP entry
- Basic TOTP support
- YubiOTP support
- Rate limiting support

### Environment Variables
| Environment             | Default value    | Description                                                                                                                                             |
|-------------------------|------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------|
| SNO_LISTEN_IP           | 0.0.0.0          | IP which SNO will listen at                                                                                                                             |
| SNO_LISTEN_PORT         | 7079             | Port which SNO will listen at                                                                                                                           |
| SNO_SECRET              |                  | OTP secret key. Enables TOTP functionality if not empty. If both this and `SNO_YUBIOTP` are empty, application will reply a random one for use and exit |
| SNO_YUBIOTP             |                  | One example of your YubiOTP. Enables YubiOTP functionality if not empty. Only the first 12 characters are used                                          |
| SNO_TITLE               | Simple Nginx OTP | Page title on OTP entry page                                                                                                                            |
| SNO_COOKIE_NAME         | sno_session      | Session cookie name                                                                                                                                     |
| SNO_COOKIE_LENGTH       | 16               | Session cookie length (recommended >=16)                                                                                                                |
| SNO_COOKIE_LIFETIME     | 14               | Session cookie lifetime in days                                                                                                                         |
| SNO_COOKIE_DOMAIN       |                  | Session cookie domain. If empty, default to current domain                                                                                              |
| SNO_RATE_LIMIT_COUNT    | 3                | How many failures till rate limit kicks in                                                                                                              |
| SNO_RATE_LIMIT_LIFETIME | 1                | Rate limit lifetime in minutes                                                                                                                          |

### Built With
* [golang](https://golang.org/)
* [go-chi/chi](https://github.com/go-chi/chi)
* [pquerna/otp](https://github.com/pquerna/otp)


## Getting Started
To get a local copy up and running follow these simple steps.
> Make sure to only allow nginx to access the application!

> Please change/ `SNO_SECRET` and `SNO_YUBIOTP` accordingly as they are examples, run without both to generate a random `SNO_SECRET` for use.

### 1a. Docker Run
```shell
docker run -d \
  --name simple-nginx-otp \
  -e SNO_LISTEN_IP=0.0.0.0 \
  -e SNO_LISTEN_PORT=7079 \
  -e SNO_SECRET=JBSWY3DPEHPK3PXP \
  -e SNO_YUBIOTP=vvvvvvcurikvhjcvnlnbecbkubjvuittbifhndhn \
  -e SNO_TITLE="Simple Nginx OTP" \
  -e SNO_COOKIE_NAME=sno_session \
  -e SNO_COOKIE_LENGTH=16 \
  -e SNO_COOKIE_LIFETIME=14 \
  -e SNO_COOKIE_DOMAIN="" \
  -e SNO_RATE_LIMIT_COUNT=3 \
  -e SNO_RATE_LIMIT_LIFETIME=1 \
  -p 7079:7079 \
  --restart unless-stopped \
  jarylc/simple-nginx-otp
```

### 1b. Docker-compose
> Please change/remove `SNO_SECRET` and `SNO_YUBIOTP` accordingly as they are examples, run without both to generate a random `SNO_SECRET` for use.
```docker-compose
simple-nginx-otp:
    image: jarylc/simple-nginx-otp
    user: nobody
    ports:
        - "7079:7079"
    environment:
        - UID=0
        - GID=0
        - SNO_LISTEN_IP=0.0.0.0
        - SNO_LISTEN_PORT=7079
        - SNO_SECRET=JBSWY3DPEHPK3PXP
        - SNO_YUBIOTP=vvvvvvcurikvhjcvnlnbecbkubjvuittbifhndhn
        - SNO_TITLE="Simple Nginx OTP"
        - SNO_COOKIE_NAME=sno_session
        - SNO_COOKIE_LENGTH=16
        - SNO_COOKIE_LIFETIME=14
        - SNO_COOKIE_DOMAIN=""
        - SNO_RATE_LIMIT_COUNT=3
        - SNO_RATE_LIMIT_LIFETIME=1
    restart: unless-stopped
```

### 1c. Binary
[Click here for the latest binaries](https://gitlab.com/jarylc/simple-nginx-otp/-/jobs/artifacts/master/browse?job=build)
> Please change/remove `SNO_SECRET` and `SNO_YUBIOTP` accordingly as they are examples, run without both to generate a random `SNO_SECRET` for use.
```shell
export UID=0
export GID=0
export SNO_LISTEN_IP=0.0.0.0
export SNO_LISTEN_PORT=7079
export SNO_SECRET=JBSWY3DPEHPK3PXP
export SNO_YUBIOTP=vvvvvvcurikvhjcvnlnbecbkubjvuittbifhndhn
export SNO_TITLE="Simple Nginx OTP"
export SNO_COOKIE_NAME=sno_session
export SNO_COOKIE_LENGTH=16
export SNO_COOKIE_LIFETIME=14
export SNO_COOKIE_DOMAIN=""
export SNO_RATE_LIMIT_COUNT=3
export SNO_RATE_LIMIT_LIFETIME=1
./simple-nginx-otp.linux-(arch)
```

### 2. Nginx
Inside the `server` block:
```nginx
error_page 401 = @error401;
location @error401 {
    return 302 /sno;
}
location /sno {
    error_page 401 /;
    proxy_pass http://127.0.0.1:7079;
    proxy_pass_request_body off;
    proxy_set_header Content-Length "";
    proxy_set_header X-Original-URI $scheme://$http_host$request_uri;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
location / {
    auth_request /sno;
    proxy_pass http://endpoint;
}
```

## Development
### Building
```shell
cd /path/to/project/folder
go build -ldflags="-w -s"
```

### Docker build
```shell
cd /path/to/project/folder
docker build .
```


## Roadmap
See the [open issues](https://gitlab.com/jarylc/simple-nginx-otp/-/issues) for a list of proposed features (and known issues).


## Contributing
Feel free to fork the repository and submit pull requests.


## License
Distributed under the MIT License. See `LICENSE` for more information.


## Contact
Jaryl Chng - git@jarylchng.com

https://jarylchng.com

Project Link: [https://gitlab.com/jarylc/simple-nginx-otp/](https://gitlab.com/jarylc/simple-nginx-otp/)
