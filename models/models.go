package models

type Config struct {
	Enviornments []Enviornment `json:"enviornments"`
	Stacks       []Stack       `json:"includeStacks"`
	UseStacks    int16         `json:"useStacks"`
	Exclude      []Stack       `json:"exclude"`
}

type Enviornment struct {
	Url    string `json:"url"`
	ApiKey string `json:"apiKey"`
	Name   string `json:"name"`
}

type Stack struct {
	Name string `json:"name"`
}

type PortList struct {
	Objects []PortObject
}

type PortObject struct {
	Spec Spec `json:"Spec"`
}

type Spec struct {
	Labels Labels `json:"Labels"`
	Name   string `json:"Name"`
}

type Labels struct {
	Image     string `json:"com.docker.stack.image"`
	Namespace string `json:"com.docker.stack.namespace"`
}

type EnvVersion struct {
	Enviornment string `json:"enviornment"`
	Stack       string `json:"stack"`
	Docker      string `json:"docker"`
	DockerPath  string `json:"dockerPath"`
}
