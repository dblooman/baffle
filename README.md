Baffle
================

Baffle is a protoype for abstracting mulitple secret backends behind an create only API.  The workflow for this includes writing a regex test for the secret, specifying the path of which the secret will be written to, as well as a fragment for reference.  A fragment is a user defined number of characters for helping identify a credential, but a fragment could also be a phrase or even a contact email.  

In this POC, a file is read from disk to store secrets, but in a production scenario, the API would recieve a PUT request either from a command line client or directly to the API. The intention is that this sits behind API Gateway so that IAM authentication can be used. The CLI signs the request for you, but Postman can be used if you want to interact with the API directly.

## Getting Started

Download Go and setup
Download the latest [Vault](https://vaultproject.io)
Create a Dynamodb table called `baffles`

## Running

Start Vault in dev mode 
```sh
vault server --dev
```

Export your Vault token and server
```sh
export VAULT_ADDR=http://127.0.0.1:8200 
export VAULT_TOKEN=00-000000-0000
``` 

Run server
```sh
cd server && go run main.go
```

Create secrets via CLI
```
go install
baffle create
```

If you are using an existing vault setup, ensure it is enabled for v2 of kv.