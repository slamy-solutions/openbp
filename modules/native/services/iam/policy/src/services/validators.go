package services

import (
	"regexp"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var actionAndResourceRegex = regexp.MustCompile(`^[A-Za-z0-9._]*\*?$`)

func validatePolicyData(name string, description string, actions []string, resources []string) error {
	if name == "" {
		return status.Error(codes.InvalidArgument, "Policy name cannot be empty")
	}

	for _, action := range actions {
		if len(action) == 0 {
			return status.Error(codes.InvalidArgument, "Policy action cannot be empty")
		}
		if !actionAndResourceRegex.MatchString(action) {
			return status.Error(codes.InvalidArgument, "Policy action ["+action+"] doesnt match \"^[A-Za-z0-9.]*\\*?$\" regex.")
		}
	}

	for _, resource := range resources {
		if len(resource) == 0 {
			return status.Error(codes.InvalidArgument, "Policy resource cannot be empty")
		}
		if !actionAndResourceRegex.MatchString(resource) {
			return status.Error(codes.InvalidArgument, "Policy resource ["+resource+"] doesnt match \"^[A-Za-z0-9.]*\\*?$\" regex.")
		}
	}

	return nil
}
