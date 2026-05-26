# Desktop OTP Gate
A lightweight TOTP gate for Nginx `auth_request`, based on [jarylc/simple-nginx-otp](https://gitlab.com/jarylc/simple-nginx-otp).

This fork keeps the original minimal reverse-auth design and adds a dark login screen, minute-based sessions, session-cookie mode, and stronger default cookie flags for the desktop workspace use case.

[**Container Image »**](https://github.com/Yakrel/desktop-otp-gate/pkgs/container/desktop-otp-gate)


## About
### Features
- Lightweight and fast
- Dark, mobile-friendly OTP login screen
- Basic TOTP support
- YubiOTP support
- Rate limiting support
- Minute-based session lifetime support

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
| SNO_COOKIE_LIFETIME     | 14               | Session cookie lifetime in days. Ignored when `SNO_COOKIE_LIFETIME_MINUTES` is greater than zero                                                        |
| SNO_COOKIE_LIFETIME_MINUTES | 0            | Session cookie lifetime in minutes. Use this for shorter sessions, for example `60` for one hour                                                        |
| SNO_SESSION_COOKIE      | false            | If true, the browser cookie has no Expires attribute and lasts until the browser clears session cookies                                                  |
| SNO_COOKIE_SECURE       | true             | Add the Secure flag to the session cookie                                                                                                                |
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
  -e SNO_COOKIE_LIFETIME_MINUTES=60 \
  -e SNO_SESSION_COOKIE=false \
  -e SNO_COOKIE_SECURE=true \
  -e SNO_COOKIE_DOMAIN="" \
  -e SNO_RATE_LIMIT_COUNT=3 \
  -e SNO_RATE_LIMIT_LIFETIME=1 \
  -p 7079:7079 \
  --restart unless-stopped \
  ghcr.io/yakrel/desktop-otp-gate:latest
```

### 1b. Docker-compose
> Please change/remove `SNO_SECRET` and `SNO_YUBIOTP` accordingly as they are examples, run without both to generate a random `SNO_SECRET` for use.
```docker-compose
desktop-otp-gate:
    image: ghcr.io/yakrel/desktop-otp-gate:latest
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
        - SNO_COOKIE_LIFETIME_MINUTES=60
        - SNO_SESSION_COOKIE=false
        - SNO_COOKIE_SECURE=true
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
export SNO_COOKIE_LIFETIME_MINUTES=60
export SNO_SESSION_COOKIE=false
export SNO_COOKIE_SECURE=true
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
