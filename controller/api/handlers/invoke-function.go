package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morty-faas/morty/controller/orchestration"
	"github.com/morty-faas/morty/controller/state"
	"github.com/morty-faas/morty/controller/telemetry"
	"github.com/morty-faas/morty/controller/types"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFunctionNotFound              = errors.New("function not found")
	ErrFunctionCantBeMarkedAsHealthy = errors.New("one or more instances of the function can't be marked as healthy")
)

func InvokeFunctionHandler(s state.State, orch orchestration.Orchestrator) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx, fnName := c.Request.Context(), c.Param("name")

		log.Debugf("Invoke function '%s'", fnName)

		fn, err := s.Get(ctx, fnName)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		if fn == nil {
			c.JSON(http.StatusNotFound, makeApiError(ErrFunctionNotFound))
			return
		}

		instance, isColdStart, err := orch.GetFunctionInstance(ctx, fn)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(instance.Endpoint)

		// Healthcheck the instance
		// Perform healthcheck against the Alpha agent
		// If alpha doesn't anwser to our requests, it probably that
		// the VM isn't ready yet to receive our requests.
		const maxHealthcheckRetries = 10
		healthcheck := instance.Endpoint.String() + "/_/health"
		for i := 0; i < maxHealthcheckRetries; i++ {
			log.Debugf("Performing healthcheck request on Alpha: %s", healthcheck)
			if _, err := http.Get(healthcheck); err != nil {
				if i == maxHealthcheckRetries-1 {
					log.Errorf("failed to perform healthcheck on Alpha: %v", err)
					c.JSON(http.StatusServiceUnavailable, makeApiError(ErrFunctionCantBeMarkedAsHealthy))
					return
				}
				time.Sleep(1 * time.Second)
				continue
			}
			log.Infof("Function '%s' is healthy and ready to receive requests", fnName)
			break
		}

		// Each invocation warn up function for 15 minutes
		if err := s.SetWithExpiry(ctx, instance.Id, 15*time.Minute); err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, makeApiError(err))
			return
		}

		// Increment the number of invocations for this function
		labels := []string{fnName, strconv.FormatBool(isColdStart)}
		telemetry.FunctionInvocationCounter.WithLabelValues(labels...).Inc()

		proxy.ModifyResponse = func(r *http.Response) error {
			// Record the execution time of the invocation
			elapsed := time.Since(start).Seconds()
			log.Debugf("Function '%s' was invoked in %vs (is cold start: %v)", fnName, elapsed, isColdStart)
			telemetry.FunctionInvocationDurationHistogram.WithLabelValues(labels...).Observe(elapsed)

			fnResponse := &types.FnInvocationResponse{}
			by, err := io.ReadAll(r.Body)
			if err != nil {
				log.Errorf("Could not read response body: %v", err)
				return err
			}
			defer r.Body.Close()

			if err := json.Unmarshal(by, &fnResponse); err != nil {
				log.Errorf("Could not unmarshal function response: %v", err)
				return err
			}

			var responseBytes []byte
			// if the function payload is a string, return it as text
			if value, ok := fnResponse.Payload.(string); ok {
				responseBytes = []byte(value)
			} else {
				responseBytes, err = json.Marshal(fnResponse.Payload)
				if err != nil {
					return err
				}
			}

			contentLength := len(responseBytes)

			r.Body = io.NopCloser(bytes.NewReader(responseBytes))
			r.ContentLength = int64(contentLength)
			r.Header.Set("Content-Length", strconv.Itoa(contentLength))

			return nil
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
