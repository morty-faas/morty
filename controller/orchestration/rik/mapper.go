package rik

import (
	"github.com/morty-faas/morty/controller/types"
	rik "github.com/rik-org/rik-go-client"
)

// mapFnToWorkload is a helper function that maps a Morty function to a RIK Workload
func mapFnToWorkload(fn *types.Function) *rik.Workload {
	apiVersion, kind := "v0", rik.KIND_FUNCTION
	name := fn.Id()
	return &rik.Workload{
		ApiVersion: &apiVersion,
		Kind:       &kind,
		Name:       &name,
		Spec: &rik.WorkloadSpec{
			Function: &rik.Function{
				Execution: &rik.FunctionExecution{
					Rootfs: &fn.ImageURL,
				},
			},
		},
	}
}
