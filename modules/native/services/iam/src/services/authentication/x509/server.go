package x509

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	x509GRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type X509IdentificationServer struct {
	x509GRPC.UnimplementedIAMAuthenticationX509ServiceServer

	nativeStub *native.NativeStub
	systemStub *system.SystemStub

	signer *x509Signer
}

func NewX509IdentificationService(ctx context.Context, systemStub *system.SystemStub, nativeStub *native.NativeStub) (*X509IdentificationServer, error) {
	err := EnsureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("failed to ensure indexes for global namespace: " + err.Error())
	}

	return &X509IdentificationServer{
		systemStub: systemStub,
		nativeStub: nativeStub,
		signer:     newX509Signer(systemStub),
	}, nil
}

func x509CollectionByNamespace(stub *system.SystemStub, namespace string) *mongo.Collection {
	if namespace == "" {
		return stub.DB.Database("openbp_global").Collection("native_iam_authentication_x509")
	} else {
		dbName := fmt.Sprintf("openbp_namespace_%s", namespace)
		return stub.DB.Database(dbName).Collection("native_iam_authentication_x509")
	}
}

func (s *X509IdentificationServer) GetRootCAInfo(ctx context.Context, in *x509GRPC.GetRootCAInfoRequest) (*x509GRPC.GetRootCAInfoResponse, error) {
	err := s.signer.EnsureReady(ctx)
	if err != nil {
		if err == ErrX509SignerVaultSealed {
			return nil, status.Error(codes.FailedPrecondition, "The vault is sealed")
		}

		return nil, status.Error(codes.Internal, "cant prepare X509 certificate signer: "+err.Error())
	}

	return &x509GRPC.GetRootCAInfoResponse{
		Certificate: s.signer.ca.Raw,
	}, status.Error(codes.OK, "")
}

func (s *X509IdentificationServer) RegisterAndGenerate(ctx context.Context, in *x509GRPC.RegisterAndGenerateRequest) (*x509GRPC.RegisterAndGenerateResponse, error) {
	_, err := x509.ParsePKCS1PublicKey(in.PublicKey)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad public key format. DER format expected.")
	}

	err = s.signer.EnsureReady(ctx)
	if err != nil {
		if err == ErrX509SignerVaultSealed {
			return nil, status.Error(codes.FailedPrecondition, "The vault is sealed")
		}

		return nil, status.Error(codes.Internal, "cant prepare X509 certificate signer: "+err.Error())
	}

	identityUUID, err := primitive.ObjectIDFromHex(in.Identity)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Identity doesnt exist. Invalid identity UUID")
	}

	creationTime := time.Now().UTC()
	certificateInfo := &X509InMongo{
		Identity:    identityUUID,
		Disabled:    false,
		Description: in.Description,
		PublicKey:   in.PublicKey,
		Created:     creationTime,
		Updated:     creationTime,
		Version:     0,
	}
	insertResponse, err := x509CollectionByNamespace(s.systemStub, in.Namespace).InsertOne(ctx, certificateInfo)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to insert new x509 certificate to the database: "+err.Error())
	}
	certificateInfo.UUID = insertResponse.InsertedID.(primitive.ObjectID)

	signedX509, err := certificateInfo.ToSignedX509(ctx, in.Namespace, s.signer)
	if err != nil {
		if err == ErrX509SignerVaultSealed {
			return nil, status.Error(codes.FailedPrecondition, "The vault is sealed")
		}

		return nil, status.Error(codes.Internal, "error while signing x509 certificate: "+err.Error())
	}

	return &x509GRPC.RegisterAndGenerateResponse{
		Raw:  signedX509,
		Info: certificateInfo.ToGRPCCertificate(in.Namespace),
	}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) Regenerate(ctx context.Context, in *x509GRPC.RegenerateRequest) (*x509GRPC.RegenerateResponse, error) {
	err := s.signer.EnsureReady(ctx)
	if err != nil {
		if err == ErrX509SignerVaultSealed {
			return nil, status.Error(codes.FailedPrecondition, "The vault is sealed")
		}
		return nil, status.Error(codes.Internal, "cant prepare X509 certificate signer: "+err.Error())
	}

	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Certificate doesnt exist. Invalid certificate UUID format.")
	}

	var foundedCertificate X509InMongo
	err = x509CollectionByNamespace(s.systemStub, in.Namespace).FindOne(ctx, bson.M{"_id": uuid}).Decode(&foundedCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Certificate not found. Probably namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while searching for certificate in the database: "+err.Error())
	}

	if foundedCertificate.Disabled {
		return nil, status.Error(codes.PermissionDenied, "cannot sign certificate which is disabled")
	}

	signedX509, err := foundedCertificate.ToSignedX509(ctx, in.Namespace, s.signer)
	if err != nil {
		if err == ErrX509SignerVaultSealed {
			return nil, status.Error(codes.FailedPrecondition, "The vault is sealed")
		}

		return nil, status.Error(codes.Internal, "error while signing x509 certificate: "+err.Error())
	}

	return &x509GRPC.RegenerateResponse{
		Certificate: signedX509,
	}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) ValidateAndGetFromRawX509(ctx context.Context, in *x509GRPC.ValidateAndGetFromRawX509Request) (*x509GRPC.ValidateAndGetFromRawX509Response, error) {
	err := s.signer.EnsureReady(ctx)
	if err != nil {
		if err == ErrX509SignerVaultSealed {
			return nil, status.Error(codes.FailedPrecondition, "The vault is sealed")
		}
		return nil, status.Error(codes.Internal, "cant prepare X509 certificate signer: "+err.Error())
	}

	cert, err := x509.ParseCertificate(in.Raw)
	if err != nil {
		return &x509GRPC.ValidateAndGetFromRawX509Response{
			Status:      x509GRPC.ValidateAndGetFromRawX509Response_INVALID_FORMAT,
			Certificate: nil,
		}, status.Error(codes.OK, "invalid format: "+err.Error())
	}

	roots := x509.NewCertPool()
	roots.AddCert(s.signer.ca)

	_, err = cert.Verify(x509.VerifyOptions{Roots: roots})
	if err != nil {
		return &x509GRPC.ValidateAndGetFromRawX509Response{
			Status:      x509GRPC.ValidateAndGetFromRawX509Response_SIGNATURE_INVALID,
			Certificate: nil,
		}, status.Error(codes.OK, "invalid signature: "+err.Error())
	}

	if len(cert.Subject.Organization) == 0 {
		return &x509GRPC.ValidateAndGetFromRawX509Response{
			Status:      x509GRPC.ValidateAndGetFromRawX509Response_INVALID_FORMAT,
			Certificate: nil,
		}, status.Error(codes.OK, "invalid format: there is no subject organization in the certificate"+err.Error())
	}
	namespace := cert.Subject.Organization[0]

	uuid, err := primitive.ObjectIDFromHex(fmt.Sprintf("%x", cert.SerialNumber))
	if err != nil {
		return &x509GRPC.ValidateAndGetFromRawX509Response{
			Status:      x509GRPC.ValidateAndGetFromRawX509Response_INVALID_FORMAT,
			Certificate: nil,
		}, status.Error(codes.OK, "invalid format: serial number invalid: "+err.Error())
	}

	var foundedCertificate X509InMongo
	err = x509CollectionByNamespace(s.systemStub, namespace).FindOne(ctx, bson.M{"_id": uuid}).Decode(&foundedCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &x509GRPC.ValidateAndGetFromRawX509Response{
				Status:      x509GRPC.ValidateAndGetFromRawX509Response_NOT_FOUND,
				Certificate: nil,
			}, status.Error(codes.OK, "")
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &x509GRPC.ValidateAndGetFromRawX509Response{
					Status:      x509GRPC.ValidateAndGetFromRawX509Response_NOT_FOUND,
					Certificate: nil,
				}, status.Error(codes.OK, "probably namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while searching for certificate in the database: "+err.Error())
	}

	publicKey, err := x509.ParsePKCS1PublicKey(foundedCertificate.PublicKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "error parsing public key of certificate from the database: "+err.Error())
	}
	if !publicKey.Equal(cert.PublicKey) {
		return &x509GRPC.ValidateAndGetFromRawX509Response{
			Status:      x509GRPC.ValidateAndGetFromRawX509Response_INVALID_FORMAT,
			Certificate: nil,
		}, status.Error(codes.OK, "invalid signature: public key doesnt match: "+err.Error())
	}

	return &x509GRPC.ValidateAndGetFromRawX509Response{
		Status:      x509GRPC.ValidateAndGetFromRawX509Response_OK,
		Certificate: foundedCertificate.ToGRPCCertificate(namespace),
	}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) Get(ctx context.Context, in *x509GRPC.GetRequest) (*x509GRPC.GetResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Invalid certificate UUID format")
	}

	var foundedCertificate X509InMongo
	err = x509CollectionByNamespace(s.systemStub, in.Namespace).FindOne(ctx, bson.M{"_id": uuid}).Decode(&foundedCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Certificate not found. Probably namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while searching for certificate in the database: "+err.Error())
	}

	return &x509GRPC.GetResponse{
		Certificate: foundedCertificate.ToGRPCCertificate(in.Namespace),
	}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) Count(ctx context.Context, in *x509GRPC.CountRequest) (*x509GRPC.CountResponse, error) {
	countResponse, err := x509CollectionByNamespace(s.systemStub, in.Namespace).CountDocuments(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &x509GRPC.CountResponse{Count: 0}, status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while counting certificates in the database: "+err.Error())
	}

	return &x509GRPC.CountResponse{Count: uint64(countResponse)}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) List(in *x509GRPC.ListRequest, out x509GRPC.IAMAuthenticationX509Service_ListServer) error {
	ctx := out.Context()

	collection := x509CollectionByNamespace(s.systemStub, in.Namespace)
	cur, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "_id", Value: -1}}))
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		return status.Error(codes.Internal, "error while searching for identity certificates in the database: "+err.Error())
	}
	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var cert X509InMongo
		err := cur.Decode(&cert)
		if err != nil {
			return status.Error(codes.Internal, "error while searching for identity certificates in the database: error decoding certificate from the stream: "+err.Error())
		}

		err = out.Send(&x509GRPC.ListResponse{
			Certificate: cert.ToGRPCCertificate(in.Namespace),
		})
		if err != nil {
			return status.Error(codes.Internal, "error while sending certificate: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) CountForIdentity(ctx context.Context, in *x509GRPC.CountForIdentityRequest) (*x509GRPC.CountForIdentityResponse, error) {
	identityUUID, err := primitive.ObjectIDFromHex(in.Identity)
	if err != nil {
		return &x509GRPC.CountForIdentityResponse{Count: 0}, status.Error(codes.OK, "Invalid identity UUID format")
	}

	countResponse, err := x509CollectionByNamespace(s.systemStub, in.Namespace).CountDocuments(ctx, bson.M{"identity": identityUUID})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &x509GRPC.CountForIdentityResponse{Count: 0}, status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while counting certificates in the database: "+err.Error())
	}

	return &x509GRPC.CountForIdentityResponse{Count: uint64(countResponse)}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) ListForIdentity(in *x509GRPC.ListForIdentityRequest, out x509GRPC.IAMAuthenticationX509Service_ListForIdentityServer) error {
	ctx := out.Context()

	identityUUID, err := primitive.ObjectIDFromHex(in.Identity)
	if err != nil {
		return status.Error(codes.OK, "Invalid identity UUID format")
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	if in.Skip != 0 {
		findOptions.SetSkip(int64(in.Skip))
	}
	if in.Limit != 0 {
		findOptions.SetLimit(int64(in.Limit))
	}

	collection := x509CollectionByNamespace(s.systemStub, in.Namespace)
	cur, err := collection.Find(ctx, bson.M{"identity": identityUUID}, findOptions)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return status.Error(codes.OK, "Namespace doesnt exist")
			}
		}

		return status.Error(codes.Internal, "error while searching for identity certificates in the database: "+err.Error())
	}
	defer cur.Close(context.Background())

	for cur.Next(ctx) {
		var cert X509InMongo
		if err := cur.Decode(&cert); err != nil {
			return status.Error(codes.Internal, "error while searching for identity certificates in the database: error decoding certificate from the stream: "+err.Error())
		}

		err = out.Send(&x509GRPC.ListForIdentityResponse{
			Certificate: cert.ToGRPCCertificate(in.Namespace),
		})
		if err != nil {
			return status.Error(codes.Internal, "error while sending certificate: "+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) Update(ctx context.Context, in *x509GRPC.UpdateRequest) (*x509GRPC.UpdateResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Invalid certificate UUID format")
	}

	var updatedCertificate X509InMongo
	err = x509CollectionByNamespace(s.systemStub, in.Namespace).FindOneAndUpdate(
		ctx,
		bson.M{"_id": uuid},
		bson.M{
			"$set":         bson.M{"description": in.NewDescription},
			"$inc":         bson.M{"version": 1},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Certificate not found. Probably namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while updating certificate in the database: "+err.Error())
	}

	return &x509GRPC.UpdateResponse{Certificate: updatedCertificate.ToGRPCCertificate(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) Delete(ctx context.Context, in *x509GRPC.DeleteRequest) (*x509GRPC.DeleteResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &x509GRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Invalid certificate UUID format.")
	}

	deleteResponse, err := x509CollectionByNamespace(s.systemStub, in.Namespace).DeleteOne(ctx, bson.M{"_id": uuid})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &x509GRPC.DeleteResponse{Existed: false}, status.Error(codes.OK, "Certificate namespace doesnt exist.")
			}
		}

		return nil, status.Error(codes.Internal, "error while deleting certificate from the database: "+err.Error())
	}

	return &x509GRPC.DeleteResponse{Existed: deleteResponse.DeletedCount > 0}, status.Error(codes.OK, "")
}
func (s *X509IdentificationServer) Disable(ctx context.Context, in *x509GRPC.DisableRequest) (*x509GRPC.DisableResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Invalid certificate UUID format")
	}

	var updatedCertificate X509InMongo
	err = x509CollectionByNamespace(s.systemStub, in.Namespace).FindOneAndUpdate(
		ctx,
		bson.M{"_id": uuid},
		bson.M{
			"$set":         bson.M{"disabled": true},
			"$inc":         bson.M{"version": 1},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.Before),
	).Decode(&updatedCertificate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}

		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Certificate not found. Probably namespace doesnt exist")
			}
		}

		return nil, status.Error(codes.Internal, "error while updating certificate in the database: "+err.Error())
	}

	return &x509GRPC.DisableResponse{WasActive: !updatedCertificate.Disabled}, status.Error(codes.OK, "")
}
