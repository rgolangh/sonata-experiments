package sonata

import (
	"fmt"
	"strings"

	"github.com/kiegroup/kogito-serverless-operator/api/v1alpha08"
	"github.com/serverlessworkflow/sdk-go/v2/model"
    "internal/backstage"
)

func newFrom(template backstage.Template) v1alpha08.Flow {
    flow := v1alpha08.Flow{
    	Start:           &model.Start{},
        Annotations:     []string{"parodos-dev/infrastructure"},
        Metadata:        map[string]model.Object{"name": model.FromString(template.Metadata.Name)},
    	Auth:            []model.Auth{},
    	States:          []model.State{},
    	Functions:       []model.Function{},
    	Retries:         []model.Retry{},
        Errors:          []model.Error{},
    }
     
    //fmt.Printf("sonata flow definition\n%+v\n", flow)
    // set name & description

    // set the extension of loading theopenapi spec of bs 

    // set Errors
    flow.Errors = append(flow.Errors, model.Error{
    	Name:        "Error on Action",
    	Code:        "java.lang.RuntimeException",
    })

    // set function for calling into backstage
    // for each action from the template, create this struct
    //  - name: runActionFetchTemplate
    // operation: specs/actions-openapi.json#fetch:template
    for i, templateStep := range template.Spec.Steps {
     //   fmt.Println("step with action " + templateStep.Action)
      //  fmt.Printf("specs/actions-openapi.json#%s\n", templateStep.Action)
        
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
                    	ErrorRef:   "Error on Action",
                    },
                },
            },

        	OperationState: &model.OperationState{
                
        		Actions:    []model.Action{model.Action{
        			Name:               templateStep.Name,
        			FunctionRef:        &model.FunctionRef{
        				RefName:      functionNameFromStep(templateStep.Action),
        			},
        		}},
        	},
        })

        if i == 0 {
            flow.Start.StateName = templateStep.Name
        }
    }
    return flow
}

func functionNameFromStep(templateStep string) string {
    return strings.Title(strings.Join(strings.Split(templateStep, ":"), ""))
}
