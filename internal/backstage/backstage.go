package backstage

type TemplateMetadata struct {
	Name        string `yaml:"name"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type ParameterProperty struct {
	Title       string `yaml:"title"`
	Type        string `yaml:"type"`
	Description string `yaml:"description"`
	UI          struct {
		Autofocus bool `yaml:"autofocus,omitempty"`
		Options   struct {
			Rows int `yaml:"rows,omitempty"`
		} `yaml:"options,omitempty"`
	} `yaml:"ui,omitempty"`
}

type Parameter struct {
	Title     string             `yaml:"title"`
	Required  []string           `yaml:"required,omitempty"`
	Properties map[string]ParameterProperty `yaml:"properties,omitempty"`
}

type StepInput struct {
	URL    string            `yaml:"url,omitempty"`
	Values map[string]string `yaml:"values,omitempty"`
}

type Step struct {
	ID     string    `yaml:"id"`
	Name   string    `yaml:"name"`
	Action string    `yaml:"action"`
	Input  StepInput `yaml:"input"`
}

type Spec struct {
	Owner      string    `yaml:"owner"`
	Type       string    `yaml:"type"`
	Parameters []Parameter `yaml:"parameters,omitempty"`
	Steps      []Step    `yaml:"steps,omitempty"`
}

type Template struct {
	APIVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   TemplateMetadata `yaml:"metadata"`
	Spec       Spec            `yaml:"spec"`
}

