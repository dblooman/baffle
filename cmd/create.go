package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dblooman/baffle/server/backends"
	awsauth "github.com/smartystreets/go-aws-auth"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return create(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func create(ctx context.Context) error {
	var secrets []backends.CreateSecret

	secretsFile, err := ioutil.ReadFile("secret.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(secretsFile), &secrets)
	if err != nil {
		return err
	}

	for _, secret := range secrets {

		ctx := context.Background()

		url := "http://localhost:8080/put"

		payload := backends.CreateSecret{
			Secret:    secret.Secret,
			Backends:  secret.Backends,
			Fragement: secret.Fragement,
			Name:      secret.Name,
			Regex:     secret.Regex,
			Path:      secret.Path,
		}

		encoded, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		r := bytes.NewReader(encoded)

		req, err := http.NewRequest("PUT", url, r)
		if err != nil {
			return err
		}

		req.WithContext(ctx)
		awsauth.Sign(req)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		fmt.Println(string(body))

	}

	return nil
}
