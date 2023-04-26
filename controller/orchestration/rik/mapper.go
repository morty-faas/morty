package rik

import (
	"github.com/morty-faas/morty/controller/types"
	rik "github.com/rik-org/rik-go-client"
)

// mapRegisteredWorkloadToFn is a helper function that maps a RIK Workload to a Morty function
func mapRegisteredWorkloadToFn(wk *rik.GetWorkloadsResponseInner) *types.Function {
	return &types.Function{
		Id:       wk.GetId(),
		Name:     wk.GetName(),
		ImageURL: *wk.GetValue().Spec.Function.Execution.Rootfs,
	}
}

// mapFnToWorkload is a helper function that maps a Morty function to a RIK Workload
func mapFnToWorkload(fn *types.Function) *rik.Workload {
	apiVersion, kind := "v0", rik.KIND_FUNCTION
	return &rik.Workload{
		ApiVersion: &apiVersion,
		Kind:       &kind,
		Name:       &fn.Name,
		Spec: &rik.WorkloadSpec{
			Function: &rik.Function{
				Execution: &rik.FunctionExecution{
					Rootfs: &fn.ImageURL,
				},
			},
		},
	}
}
