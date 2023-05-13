package registry

import (
	"errors"
	"html"
	"net/http"
	"os"

	"github.com/morty-faas/morty/controller/telemetry"
	"github.com/morty-faas/morty/registry/builder"
	log "github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"
)

var (
	ErrRequiredFunctionName    = errors.New("the function name is required")
	ErrRequiredFunctionRuntime = errors.New("the function runtime is required")
	ErrRequiredFunctionVersion = errors.New("the function version is required")
	ErrFunctionVersionInvalid  = errors.New("the function version must be a valid semantic version and must start with 'v'")
	ErrInvalidFunctionArchive  = errors.New("the function code archive must be a valid zip file")
)

func (s *Server) BuildHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the file, the headers of the file from the multipart request form
	f, _, err := r.FormFile("archive")
	if err != nil {
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
		return
	}
	defer f.Close()

	// Validate DTO
	functionName, functionRuntime, functionVersion := html.EscapeString(r.Form.Get("name")), html.EscapeString(r.Form.Get("runtime")), html.EscapeString(r.Form.Get("version"))
	if functionName == "" {
		s.APIErrorResponse(w, makeAPIError(http.StatusBadRequest, ErrRequiredFunctionName))
		return
	}

	if functionRuntime == "" {
		s.APIErrorResponse(w, makeAPIError(http.StatusBadRequest, ErrRequiredFunctionRuntime))
		return
	}

	if functionVersion == "" {
		s.APIErrorResponse(w, makeAPIError(http.StatusBadRequest, ErrRequiredFunctionVersion))
		return
	}

	if !semver.IsValid(functionVersion) {
		s.APIErrorResponse(w, makeAPIError(http.StatusBadRequest, ErrFunctionVersionInvalid))
		return
	}

	id := functionName + "-" + functionVersion

	// Build the function image with the given options
	opts := &builder.BuildOptions{
		Id:      id,
		Runtime: functionRuntime,
		Archive: f,
	}

	image, err := s.builder.ImageBuild(r.Context(), opts)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, builder.ErrInvalidRuntime) {
			status = http.StatusBadRequest
		}

		s.APIErrorResponse(w, makeAPIError(status, err))
		return
	}

	// Upload the file to the remote storage before returning a response to the user
	f, _ = os.Open(image)
	if err := s.storage.PutFile(id, f); err != nil {
		s.APIErrorResponse(w, makeAPIError(http.StatusInternalServerError, err))
		return
	}

	log.Infof("build/%s: function build successful", opts.Id)

	telemetry.FunctionBuildCounter.WithLabelValues(opts.Id, opts.Runtime).Inc()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("/v1/functions/" + functionName + "/" + functionVersion))
}

func makeAPIError(status int, err error) *APIError {
	return &APIError{
		StatusCode: status,
		Message:    err.Error(),
	}
}
