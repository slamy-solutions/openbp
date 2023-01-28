package native

import (
	"context"
	"regexp"

	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"
	auth "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools/auth"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/tools/libs/sdk/golang/native/namespace"
)

var /* const */ nameValidator = regexp.MustCompile(`^[A-Za-z0-9]+$`)

type namespaceService struct {
	namespaceGRPC.UnimplementedNativeNamespaceServiceServer

	log *log.Logger

	modules *tools.ModulesStub
}

func RegisterNamespaceService(modules *tools.ModulesStub, grpcServer *grpc.Server) {
	namespaceGRPC.RegisterNativeNamespaceServiceServer(grpcServer, &namespaceService{
		log:     log.StandardLogger(),
		modules: modules,
	})
}

func (s *namespaceService) Ensure(ctx context.Context, in *namespaceGRPC.EnsureNamespaceRequest) (*namespaceGRPC.EnsureNamespaceResponse, error) {
	err := auth.AuthorizeRPC(ctx, s.modules, auth.Scope{
		Namespace: "",
		Resources: []string{"native.namespace"},
		Actions:   []string{"native.namespace.create"},
	})
	if err != nil {
		return nil, err
	}

	response, err := s.modules.Native.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
		Name:        in.Name,
		FullName:    in.FullName,
		Description: in.Description,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				return nil, status.Error(codes.InvalidArgument, st.Message())
			}
		}

		s.log.WithError(err).WithField("module", "native").WithField("service", "namespace").Error("Failed to get namespace.")
		return nil, status.Error(codes.Internal, "Failed to get identity information: "+err.Error())
	}

	return &namespaceGRPC.EnsureNamespaceResponse{
		Namespace: &namespaceGRPC.Namespace{
			Name:        response.Namespace.Name,
			FullName:    response.Namespace.FullName,
			Description: response.Namespace.Description,
			Updated:     response.Namespace.Updated,
			Created:     response.Namespace.Created,
			Version:     response.Namespace.Version,
		},
	}, status.Error(codes.OK, "")
}

func (s *namespaceService) Create(ctx context.Context, in *namespaceGRPC.CreateNamespaceRequest) (*namespaceGRPC.CreateNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func (s *namespaceService) Update(ctx context.Context, in *namespaceGRPC.UpdateNamespaceRequest) (*namespaceGRPC.UpdateNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *namespaceService) Get(ctx context.Context, in *namespaceGRPC.GetNamespaceRequest) (*namespaceGRPC.GetNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (s *namespaceService) GetAll(in *namespaceGRPC.GetAllNamespacesRequest, out namespaceGRPC.NativeNamespaceService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}

func (s *namespaceService) Delete(ctx context.Context, in *namespaceGRPC.DeleteNamespaceRequest) (*namespaceGRPC.DeleteNamespaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

func (s *namespaceService) Exists(ctx context.Context, in *namespaceGRPC.IsNamespaceExistRequest) (*namespaceGRPC.IsNamespaceExistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}

func (s *namespaceService) Stat(ctx context.Context, in *namespaceGRPC.GetNamespaceStatisticsRequest) (*namespaceGRPC.GetNamespaceStatisticsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stat not implemented")
}
