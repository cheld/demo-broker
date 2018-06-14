package broker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/pmorie/osb-broker-lib/pkg/broker"

	"reflect"

	osb "github.com/pmorie/go-open-service-broker-client/v2"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.
	return &BusinessLogic{
		async:       o.Async,
		catalogPath: o.CatalogPath,
		instances:   make(map[string]*exampleInstance, 10),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// url to remote broker
	catalogPath string
	// Synchronize go routines.
	sync.RWMutex
	// Add fields here! These fields are provided purely as an example
	instances map[string]*exampleInstance
}

var _ broker.Interface = &BusinessLogic{}

func truePtr() *bool {
	b := true
	return &b
}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*broker.CatalogResponse, error) {
	// Your catalog business logic goes here
	response := &broker.CatalogResponse{}
	//osbResponse := &osb.CatalogResponse{
	//	Services: []osb.Service{
	//		{
	//			Name:          "example-ceiser-service",
	//			ID:            "4f6e6cf6-ffdd-425f-a2c7-3c9258ad246a",
	//			Description:   "New example service from the osb starter pack!",
	//			Bindable:      true,
	//			PlanUpdatable: truePtr(),
	//			Metadata: map[string]interface{}{
	//				"displayName": "Example ceiser backend service v2",
	//				"imageUrl":    "https://cdn.iconverticons.com/files/png/0174d53a1a739691_256x256.png",
	//			},
	//			Plans: []osb.Plan{
	//				{
	//					Name:        "default",
	//					ID:          "86064792-7ea2-467b-af93-ac9694d96d5b",
	//					Description: "The default plan for the starter pack example service",
	//					Free:        truePtr(),
	//					Schemas: &osb.Schemas{
	//						ServiceInstance: &osb.ServiceInstanceSchema{
	//							Create: &osb.InputParametersSchema{
	//								Parameters: map[string]interface{}{
	//									"type": "object",
	//									"properties": map[string]interface{}{
	//										"color": map[string]interface{}{
	//											"type":    "string",
	//											"default": "Clear",
	//											"enum": []string{
	//												"Clear",
	//												"Magenta",
	//												"Grey",
	//											},
	//										},
	//									},
	//								},
	//							},
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}
	resp, err := http.Get(b.catalogPath)
	if err != nil {
		glog.Infof("Error reading json %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Infof("Error reading body %s", err)
		return nil, err
	}
	var osbResponse osb.CatalogResponse
	json.Unmarshal(body, &osbResponse)

	fmt.Println("Path %s", b.catalogPath)
	fmt.Println(osbResponse)
	glog.Infof("get catalog %c", c.Request)
	glog.Infof("catalog response: %#+v", osbResponse)

	response.CatalogResponse = osbResponse
	return response, nil
}

func (b *BusinessLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*broker.ProvisionResponse, error) {
	// Your provision business logic goes here
	fmt.Println("provision")
	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.ProvisionResponse{}

	exampleInstance := &exampleInstance{
		ID:        request.InstanceID,
		ServiceID: request.ServiceID,
		PlanID:    "my plan request.PlanID",
		Params:    request.Parameters,
	}

	glog.Infof("Provisioning: %#+v", exampleInstance)

	// Check to see if this is the same instance
	if i := b.instances[request.InstanceID]; i != nil {
		if i.Match(exampleInstance) {
			response.Exists = true
			return &response, nil
		} else {
			// Instance ID in use, this is a conflict.
			description := "InstanceID in use"
			return nil, osb.HTTPStatusCodeError{
				StatusCode:  http.StatusConflict,
				Description: &description,
			}
		}
	}
	b.instances[request.InstanceID] = exampleInstance

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*broker.DeprovisionResponse, error) {
	// Your deprovision business logic goes here
	fmt.Println("deprovision")
	// example implementation:
	b.Lock()
	defer b.Unlock()

	response := broker.DeprovisionResponse{}

	delete(b.instances, request.InstanceID)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*broker.LastOperationResponse, error) {
	// Your last-operation business logic goes here
	fmt.Println("last operation")
	glog.Infof("LastOperation")
	return nil, nil
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*broker.BindResponse, error) {
	// Your bind business logic goes here
	fmt.Println("bind")
	// example implementation:
	b.Lock()
	defer b.Unlock()

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	response := broker.BindResponse{
		BindResponse: osb.BindResponse{
			Credentials: instance.Params,
		},
	}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*broker.UnbindResponse, error) {
	// Your unbind business logic goes here
	fmt.Println("unbind")
	return &broker.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*broker.UpdateInstanceResponse, error) {
	fmt.Println("update")
	// Your logic for updating a service goes here.
	response := broker.UpdateInstanceResponse{}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	fmt.Printf("validate broker %s", version)
	glog.Infof("ValidateBrokerAPIVersion %s", version)
	return nil
}

// example types

// exampleInstance is intended as an example of a type that holds information about a service instance
type exampleInstance struct {
	ID        string
	ServiceID string
	PlanID    string
	Params    map[string]interface{}
}

func (i *exampleInstance) Match(other *exampleInstance) bool {
	return reflect.DeepEqual(i, other)
}
