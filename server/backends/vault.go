package backends

import (
	"path"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	Data   CreateSecret
	Client *vault.Client
}

func (v VaultClient) Write() (Response, error) {
	l := v.Client.Logical()

	path := v.Data.Path

	mountPath, v2, err := vaultKV(path, v.Client)
	if err != nil {
		return Response{}, err
	}

	data := map[string]interface{}{
		"value": v.Data.Secret,
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		data = map[string]interface{}{
			"data":    data,
			"options": map[string]interface{}{},
		}
	}

	resp, err := l.Write(path, data)
	if err != nil {
		return Response{}, err
	}

	return Response{Message: resp.Data}, nil

}

func kvPreflightVersionRequest(client *vault.Client, path string) (string, int, error) {
	r := client.NewRequest("GET", "/v1/sys/internal/ui/mounts/"+path)
	resp, err := client.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := vault.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	var mountPath string
	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}
	version := versionRaw.(string)
	switch version {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}

func vaultKV(path string, client *vault.Client) (string, bool, error) {
	mountPath, version, err := kvPreflightVersionRequest(client, path)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}

func addPrefixToVKVPath(p, mountPath, apiPrefix string) string {
	p = strings.TrimPrefix(p, mountPath)
	return path.Join(mountPath, apiPrefix, p)
}
