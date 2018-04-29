package backends

type Backender interface {
	Write() (Response, error)
}

type Response struct {
	Message map[string]interface{}
}

func Put(b Backender) (Response, error) {
	return b.Write()
}

type CreateSecret struct {
	Path      string   `json:"path"`
	Secret    string   `json:"secret"`
	Backends  []string `json:"backends"`
	Regex     string   `json:"regex"`
	Fragement string   `json:"fragment"`
	Version   int64    `json:"version"`
	Name      string   `json:"name"`
}
