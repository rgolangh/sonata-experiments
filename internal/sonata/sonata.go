package sonata

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kiegroup/kogito-serverless-operator/api/v1alpha08"
	"github.com/rgolangh/sonata-experiments/internal/backstage"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func NewFrom(template backstage.Template) v1alpha08.Flow {
	flow := v1alpha08.Flow{
		Start:       &model.Start{},
		Annotations: []string{"parodos-dev/infrastructure"},
		Metadata: map[string]model.Object{
			"id":   model.FromString(strings.ReplaceAll(strings.ToLower(template.Metadata.Name), " ", "-")),
			"name": model.FromString(template.Metadata.Name),
		},
		States: []model.State{
			{
				BaseState: model.BaseState{Name: "Handle Error", End: &model.End{Compensate: true}},
				OperationState: &model.OperationState{
					Actions: []model.Action{
						{
							Name: "Error action",
							FunctionRef: &model.FunctionRef{
								RefName: "sysout",
								Arguments: map[string]model.Object{
									"message": model.Object{StrVal: "Error on workflow, triggering componesation"},
								},
							},
						},
					},
				},
			},
		},
		Functions: []model.Function{},
		Errors: []model.Error{
			{
				Name: "Error on Action",
				Code: "java.lang.RuntimeException",
			},
		},
	}

	// set the extension of loading theopenapi spec of bs

	for i, templateStep := range template.Spec.Steps {
		if i == 0 {
			flow.Start.StateName = templateStep.Name
		}
		// todo don't append twice the same action
		flow.Functions = append(flow.Functions, model.Function{
			Common:    model.Common{},
			Name:      functionNameFromStep(templateStep.Action),
			Operation: fmt.Sprintf("spec/actions-openapi.json#%s", templateStep.Action),
		})

		flow.States = append(flow.States, model.State{
			BaseState: model.BaseState{
				Name: templateStep.Name,
				OnErrors: []model.OnError{
					{
						ErrorRef: "Error on Action",
						Transition: &model.Transition{
							NextState: "Handle Error",
						},
					},
				},
			},

			OperationState: &model.OperationState{

				Actions: []model.Action{model.Action{
					Name: templateStep.Name,
					FunctionRef: &model.FunctionRef{
						RefName:   functionNameFromStep(templateStep.Action),
						Arguments: stepInputToMap(templateStep),
					},
				}},
			},
		})

	}
	return flow
}

func stepInputToMap(templateStep backstage.Step) map[string]model.Object {
	m := make(map[string]model.Object)
	for key, value := range templateStep.Input {
		switch v := value.(type) {
		case int32:
			m[key] = model.FromInt(int(v))
		case string:
			m[key] = model.FromString(v)
		case map[interface{}]interface{}:
			_, err := json.Marshal(v)
			if err != nil {
				//fmt.Fprintln(os.Stderr, fmt.Errorf("failed to parse value to Raw: %v", err))
			}

			newm := make(map[string]string)
			for k, val := range v {
				newm[k.(string)] = val.(string)
			}
			m[key] = model.FromRaw(newm)
		case []interface{}:
			s := make([]string, len(v))
			for i, av := range v {
				s[i] = fmt.Sprint(av)
			}
			m[key] = model.FromRaw(s)

		case interface{}:
			nm, ok := v.(map[interface{}]interface{})
			if !ok {
				//fmt.Fprintln(os.Stderr, "failed to convert to map of strings")
			}
			m[key] = model.FromRaw(nm[key])
		default:
			//fmt.Printf("\ndefault type ##### with %+v\n")
		}
	}
	return m
}

func functionNameFromStep(templateStep string) string {
	return strings.Title(strings.Join(strings.Split(templateStep, ":"), ""))
}
