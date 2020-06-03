# Go-Proxy-Server
This task is an implementation of reverse proxy in **golang** using **echo framework**.

The authorization is done through JWT.

The verification of user and generation of token is in `Auth service`. **Endpoint:** `/auth`

The `user service` has two end points: `/user/profile` that is a secured one and `/service/name` which is not secured 
that is it can be accessed without being authorized.

The `proxy service` acts as a **reverse proxy** for these two services.


### How does it all work

Make all the three servers running by:

```
go run main.go (Inside Auth folder): This server runs on :1324 port
go run main.go (Inside User folder): This server runs on :1325 port
go run main.go (Inside Proxy folder): This server runs on :1323 port

```
From POSTMAN, try hitting these endpoints with these instructions

```
**GET** localhost:1323/auth :- Set the header as Username:<string> where string is the name that has to be verified. You will recieve the a token upon a successful verification. 

**GET** localhost:1323/user/profile :-  Send the token recieved above as the bearer token and send the request. If the token is valid then the user details will be returned.

**GET** localhost:1323/service/name :- You will recieve the name of the microservice upon hitting this endpoint.
```
