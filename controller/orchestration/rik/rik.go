package rik

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/types"
	rik "github.com/rik-org/rik-go-client"
	log "github.com/sirupsen/logrus"
)

type adapter struct {
	cfg    *Config
	client *rik.APIClient
}

type Config struct {
	// Cluster is the address of the RIK controller
	Cluster string `yaml:"cluster"`
}

var _ orchestration.Orchestrator = (*adapter)(nil)

// NewOrchestrator initializes the RIK orchestrator adapter.
func NewOrchestrator(cfg *Config) (orchestration.Orchestrator, error) {
	log.Info("Orchestrator engine 'rik' successfully initialized")

	client := rik.NewAPIClient(&rik.Configuration{
		Servers: rik.ServerConfigurations{
			rik.ServerConfiguration{
				URL: cfg.Cluster,
			},
		},
	})

	return &adapter{cfg, client}, nil
}

func (a *adapter) GetFunctions(ctx context.Context) ([]*types.Function, error) {
	workloads, err := a.getWorkloads(ctx)
	if err != nil {
		return nil, err
	}

	var functions []*types.Function
	for _, meta := range workloads {
		workload := meta.GetValue()
		// Filter on function elements only
		if workload.GetKind() == rik.KIND_FUNCTION {
			name, version := workload.GetName(), types.DefaultFunctionVersion

			// Try to split the workload name as we have defined that the name of a RIK workload is fnName:fnVersion
			tokens := strings.Split(workload.GetName(), ":")
			if len(tokens) == 2 {
				name, version = tokens[0], tokens[1]
			}

			functions = append(functions, &types.Function{
				Name:           name,
				Version:        version,
				ImageURL:       *workload.GetSpec().Function.Execution.Rootfs,
				OrchestratorId: meta.GetId(),
			})
		}
	}

	return functions, nil
}

func (a *adapter) CreateFunction(ctx context.Context, fn *types.Function) (*types.Function, error) {
	r := a.client.WorkloadsApi.CreateWorkload(ctx).Body(*mapFnToWorkload(fn))
	wk, _, err := a.client.WorkloadsApi.CreateWorkloadExecute(r)
	if err != nil {
		return nil, err
	}

	fn.OrchestratorId = wk.CreateWorkloadResponse.GetId()
	return fn, nil
}

func (a *adapter) GetFunctionInstance(ctx context.Context, fn *types.Function) (*types.FnInstance, bool, error) {
	// isColdStart used to control wether it is a cold start or not for the instance
	isColdStart := false

	instances, err := a.getWorkloadInstances(ctx, fn.OrchestratorId)
	if err != nil {
		return nil, false, err
	}

	if len(instances) == 0 {
		log.Debugf("Deploying new instance for function: %+v", fn)
		isColdStart = true

		if err := a.createWorkloadInstance(ctx, fn.OrchestratorId, fn.Name); err != nil {
			err := fmt.Errorf("Failed to create instance: %v", err)
			log.Error(err)
			return nil, false, err
		}

		time.Sleep(500 * time.Millisecond)

		instances, err = a.getWorkloadInstances(ctx, fn.OrchestratorId)
		if err != nil {
			return nil, false, err
		}
	}

	log.Debugf("%d instance(s)", len(instances))

	rikIn := instances[rand.Intn(len(instances))]

	url, _ := url.Parse(a.cfg.Cluster)
	url, _ = url.Parse(fmt.Sprintf("%s://%s:%d", url.Scheme, url.Hostname(), rikIn.Spec.Function.Exposure.GetPort()))

	instance := &types.FnInstance{
		Id:       rikIn.GetId(),
		Function: fn,
		Endpoint: url,
	}

	return instance, isColdStart, nil
}

// getWorkloads is a helper function to retrieve all the workloads from the RIK cluster
func (a *adapter) getWorkloads(ctx context.Context) ([]rik.GetWorkloadsResponseInner, error) {
	r := a.client.WorkloadsApi.GetWorkloads(ctx)
	workloads, _, err := a.client.WorkloadsApi.GetWorkloadsExecute(r)
	if err != nil {
		return nil, err
	}
	return workloads, nil
}

// getWorkloadInstances is a helper function to retrieve all the instances of the given workload.
func (a *adapter) getWorkloadInstances(ctx context.Context, id string) ([]rik.Instance, error) {
	r := a.client.InstancesApi.GetWorkloadInstances(ctx, id)
	data, _, err := a.client.InstancesApi.GetWorkloadInstancesExecute(r)
	if err != nil {
		return nil, err
	}

	return data.GetInstances(), nil
}

// createWorkloadInstance is a helper function to create a workload instance
func (a *adapter) createWorkloadInstance(ctx context.Context, workloadId, name string) error {
	in := rik.CreateInstanceRequest{
		WorkloadId: workloadId,
	}

	r := a.client.InstancesApi.CreateWorkloadInstance(ctx).CreateInstanceRequest(in)
	_, res, err := a.client.InstancesApi.CreateWorkloadInstanceExecute(r)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New("RIK returned non 201 HTTP Status Code for create instance")
	}

	return nil
}

func (a *adapter) DeleteFunctionInstance(ctx context.Context, fn *types.FnInstance) error {
	input := rik.DeleteInstanceRequest{
		Id: &fn.Id,
	}

	_, err := a.client.InstancesApi.DeleteInstance(ctx).DeleteInstanceRequest(input).Execute()
	return err
}
