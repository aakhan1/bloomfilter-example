# bloomfilter-example

To run example, use docker-compose
````
docker-compose up
````

This example leverages the [bloom filter](https://www.krakend.io/docs/authorization/revoking-tokens/) plugin in krakend to blacklist tokens.

A. To create a token use:
````shell script
curl -X POST 'http://localhost:8080/api/v1/sessions'

# Example output
{"session_id":"eyJhbGciOiJIUzI1NiIsImtpZCI6ImRlZmF1bHQiLCJ0eXA....."}
````
The jwt payload includes a `jti` key which is used to blacklist the token.  The krakend.json must be configured with the same key.  See config/krakend.json

````json
{
  "exp": 7600309454,
  "iat": 1600309454,
  "jti": "417b884d-db17-42c0-8222-9eca67b5420a",
  "username": "joe"
}
````


B. Test out a secured api.  This will fail if no auth token is passed.
````shell script
curl -X GET 'http://localhost:3000/messages' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6ImRlZmF1.....'
````

C. Blacklist the token generated in step A.
````shell script
curl -X DELETE 'http://localhost:3000/sessions' \
-H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6ImRlZmF1.....'
````

D. Run the API in step B again, this time it should fail because token is black listed.


