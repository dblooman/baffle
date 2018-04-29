package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/dblooman/baffle/server/backends"
	"github.com/dblooman/baffle/server/logger"
	"github.com/dblooman/baffle/server/storage"
	"github.com/dblooman/baffle/server/validate"
	"github.com/gin-gonic/gin"

	vault "github.com/hashicorp/vault/api"
)

type API struct {
	Dynamodb *dynamodb.DynamoDB
}

func main() {

	api := new(API)

	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	svc := dynamodb.New(session.New(config))

	api.Dynamodb = svc

	r := gin.Default()

	r.GET("/status", Status)
	r.PUT("/put", api.PutSecret)
	r.Run()
}

var Vault *vault.Client

func init() {
	registerVault()
}

func (a *API) PutSecret(c *gin.Context) {
	var data backends.CreateSecret

	c.BindJSON(&data)

	ok, err := validate.Ensure(data)
	if !ok {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, provider := range data.Backends {
		if provider == "vault" {
			v := backends.VaultClient{
				Data:   data,
				Client: Vault,
			}
			resp, err := backends.Put(v)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			version, err := resp.Message["version"].(json.Number).Int64()
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			data.Version = version

			err = storage.Put(a.Dynamodb, data)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, gin.H{
				"status":  "added",
				"message": resp.Message,
			})
		}
	}

}

func Status(c *gin.Context) {

	c.JSON(200, gin.H{
		"status": "ok",
	})

}

func registerVault() {
	client, err := vault.NewClient(nil)
	if err != nil {
		panic(err)
	}

	logger.Info("vault initialised", nil)

	Vault = client
}
