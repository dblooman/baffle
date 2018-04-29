// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dblooman/baffle/server/backends"

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
