package sonata

import (
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
		Metadata:    map[string]model.Object{"name": model.FromString(template.Metadata.Name)},
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

	//fmt.Printf("sonata flow definition\n%+v\n", flow)
	// set name & description

	// set the extension of loading theopenapi spec of bs

	for i, templateStep := range template.Spec.Steps {
		if i == 0 {
			flow.Start.StateName = templateStep.Name
		}
		// todo don't append twice the same action
		flow.Functions = append(flow.Functions, model.Function{
			Common:    model.Common{},
			Name:      fmt.Sprintf("runAction%s", functionNameFromStep(templateStep.Action)),
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
						RefName: functionNameFromStep(templateStep.Action),
					},
				}},
			},
		})

	}
	return flow
}

func functionNameFromStep(templateStep string) string {
	return strings.Title(strings.Join(strings.Split(templateStep, ":"), ""))
}
