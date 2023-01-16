# go-jwt-demo
authentication of go server


# launch project
```golang
go run .
```

- **step1: request home page. verify that the middleware is valid.**
```shell
curl --location --request GET 'http://127.0.0.1:8888/other/home'
```

- **step2: create user**
```shell
curl --location --request GET 'http://127.0.0.1:8888/auth/signup' \
--header 'Email: 717750878@qq.com' \
--header 'Username: lrf' \
--header 'Passwordhash: abc123' \
--header 'Fullname: luoruofeng'
```

- **step3: sign in**
```shell
curl --location --request GET 'http://127.0.0.1:8888/auth/signin' \
--header 'Email: 717750878@qq.com' \
--header 'Passwordhash: abc123'
- **the token is obtained as follows**
SFMyNTY=.eyJhdWQiOiJmcm9udGVuZC5rbm93c2VhcmNoLm1sIiwiZXhwIjoiMTY3Mjc5NDExMyIsImlzcyI6Imtub3dzZWFyY2gubWwifQ==.m9KGQhJBrd/lwLTsqcyeoO+/FLekYyaCEsEYaSBnEMM=
```

- **step4: request home page with token**
```shell
curl --location --request GET 'http://127.0.0.1:8888/other/home' \
--header 'Token: SFMyNTY=.eyJhdWQiOiJmcm9udGVuZC5rbm93c2VhcmNoLm1sIiwiZXhwIjoiMTY3Mjc5NDExMyIsImlzcyI6Imtub3dzZWFyY2gubWwifQ==.m9KGQhJBrd/lwLTsqcyeoO+/FLekYyaCEsEYaSBnEMM='
```

# build docker image
```shell
docker build . -t luoruofeng/auth-demo
docker run -d -p 8888:8888 luoruofeng/auth-demo
```

# use third lib
```shell
go get -u github.com/golang-jwt/jwt/v4
```