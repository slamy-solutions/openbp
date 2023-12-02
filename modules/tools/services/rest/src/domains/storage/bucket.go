package storage

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bucket "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/bucket"
)

type bucketRouter struct {
	nativeStub *native.NativeStub

	logger *logrus.Entry
}

type listBucketsRequest struct {
	Namespace string `form:"namespace"`
	Skip      int64  `form:"skip" binding:"gte=0"`
	Limit     int64  `form:"limit" binding:"gte=0,lte=100"`
}
type listBucketsResponse struct {
	Buckets    []formatedBucket `json:"buckets"`
	TotalCount int64            `json:"totalCount"`
}

func (r *bucketRouter) List(ctx *gin.Context) {
	var requestData listBucketsRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": requestData.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket"},
			Actions:              []string{"native.storage.bucket.list"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	bucketClient, err := r.nativeStub.Services.Storage.Bucket.List(ctx.Request.Context(), &bucket.ListBucketsRequest{
		Namespace: requestData.Namespace,
		Skip:      uint32(requestData.Skip),
		Limit:     uint32(requestData.Limit),
	})
	if err != nil {
		err := errors.New("failed to list buckets: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer bucketClient.CloseSend()

	foundedBuckets := make([]formatedBucket, 0)
	for {
		data, err := bucketClient.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err := errors.New("failed to list buckets: " + err.Error())
			logger.Error(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		foundedBuckets = append(foundedBuckets, FormatedBuceketFromGRPC(data.Bucket))
	}

	countResponse, err := r.nativeStub.Services.Storage.Bucket.Count(ctx.Request.Context(), &bucket.CountBucketsRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		err := errors.New("failed to count buckets: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, listBucketsResponse{
		Buckets:    foundedBuckets,
		TotalCount: int64(countResponse.Count),
	})
}

type createBucketRequest struct {
	Namespace string `json:"namespace" binding:"lte=32,gte=0"`
	Name      string `json:"name" binding:"required,lte=128,gt=0,regexp=^[a-zA-Z0-9-_]+$"`
	Hidden    bool   `json:"hidden"`
}
type createBucketResponse struct {
	Bucket formatedBucket `json:"bucket"`
}

func (r *bucketRouter) Create(ctx *gin.Context) {
	var requestData createBucketRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace":   requestData.Namespace,
		"bucket.name": requestData.Name,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket." + requestData.Name},
			Actions:              []string{"native.storage.bucket.create"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	bucketCreateResponse, err := r.nativeStub.Services.Storage.Bucket.Create(ctx.Request.Context(), &bucket.CreateBucketRequest{
		Namespace: requestData.Namespace,
		Name:      requestData.Name,
		Hidden:    requestData.Hidden,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid bucket name"})
				return
			}

			if st.Code() == codes.AlreadyExists {
				ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "bucket with same name already exists"})
				return
			}
		}

		err := errors.New("failed to create bucket: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Bucket created")
	ctx.JSON(http.StatusOK, createBucketResponse{
		Bucket: FormatedBuceketFromGRPC(bucketCreateResponse.Bucket),
	})
}

type deleteBucketRequest struct {
	Namespace string `form:"namespace"`
	Name      string `form:"name" binding:"lte=128,gt=0,regexp=^[a-zA-Z0-9-_]+$"`
	UUID      string `form:"uuid" binding:"lte=32,gt=0"`
}
type deleteBucketResponse struct{}

func (r *bucketRouter) Delete(ctx *gin.Context) {
	var requestData deleteBucketRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace":   requestData.Namespace,
		"bucket.name": requestData.Name,
		"bucket.uuid": requestData.UUID,
	})

	if requestData.UUID == "" && requestData.Name == "" {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "uuid or name is required"})
		return
	}

	scopes := make([]*auth.Scope, 0, 1)
	if requestData.UUID != "" {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket"},
			Actions:              []string{"native.storage.bucket.delete"},
			NamespaceIndependent: false,
		})
	} else {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket." + requestData.Name},
			Actions:              []string{"native.storage.bucket.delete"},
			NamespaceIndependent: false,
		})
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, scopes)
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	if requestData.UUID != "" {
		_, err = r.nativeStub.Services.Storage.Bucket.DeleteByUUID(ctx.Request.Context(), &bucket.DeleteBucketByUUIDRequest{
			Namespace: requestData.Namespace,
			Uuid:      requestData.UUID,
		})
	} else {
		_, err = r.nativeStub.Services.Storage.Bucket.Delete(ctx.Request.Context(), &bucket.DeleteBucketRequest{
			Namespace: requestData.Namespace,
			Name:      requestData.Name,
		})
	}
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "bucket not found"})
				return
			}
		}

		err := errors.New("failed to delete bucket: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Bucket deleted")
	ctx.JSON(http.StatusOK, deleteBucketResponse{})
}

type getBucketRequest struct {
	Namespace string `form:"namespace"`
	Name      string `form:"name" binding:"lte=128,gt=0,regexp=^[a-zA-Z0-9-_]+$"`
	UUID      string `form:"uuid" binding:"lte=32,gt=0"`
}
type getBucketResponse struct {
	Bucket formatedBucket `json:"bucket"`
}

func (r *bucketRouter) Get(ctx *gin.Context) {
	var requestData getBucketRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace":   requestData.Namespace,
		"bucket.name": requestData.Name,
		"bucket.uuid": requestData.UUID,
	})

	if requestData.UUID == "" && requestData.Name == "" {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": "uuid or name is required"})
		return
	}

	scopes := make([]*auth.Scope, 0, 1)
	if requestData.UUID != "" {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket"},
			Actions:              []string{"native.storage.bucket.get"},
			NamespaceIndependent: false,
		})
	} else {
		scopes = append(scopes, &auth.Scope{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket." + requestData.Name},
			Actions:              []string{"native.storage.bucket.get"},
			NamespaceIndependent: false,
		})
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, scopes)
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	var foundedBucket *bucket.Bucket
	if requestData.UUID != "" {
		response, rErr := r.nativeStub.Services.Storage.Bucket.GetByUUID(ctx.Request.Context(), &bucket.GetBucketByUUIDRequest{
			Namespace: requestData.Namespace,
			Uuid:      requestData.UUID,
		})
		if rErr != nil {
			err = rErr
		} else {
			foundedBucket = response.Bucket
		}
	} else {
		response, rErr := r.nativeStub.Services.Storage.Bucket.Get(ctx.Request.Context(), &bucket.GetBucketRequest{
			Namespace: requestData.Namespace,
			Name:      requestData.Name,
		})

		if rErr != nil {
			err = rErr
		} else {
			foundedBucket = response.Bucket
		}
	}
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "bucket not found"})
				return
			}
		}

		err := errors.New("failed to get bucket: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getBucketResponse{
		Bucket: FormatedBuceketFromGRPC(foundedBucket),
	})
}

type updateBucketRequest struct {
	Namespace string `json:"namespace" binding:"lte=32,gte=0"`
	UUID      string `json:"uuid" binding:"lte=32,gt=0"`
	NewName   string `json:"newName" binding:"required,lte=128,gt=0,regexp=^[a-zA-Z0-9-_]+$"`
	NewHidden bool   `json:"newHidden"`
}

type updateBucketResponse struct {
	Bucket formatedBucket `json:"bucket"`
}

func (r *bucketRouter) Update(ctx *gin.Context) {
	var requestData updateBucketRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace":        requestData.Namespace,
		"bucket.newName":   requestData.NewName,
		"bucket.newHidden": requestData.NewHidden,
		"bucket.uuid":      requestData.UUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.storage.bucket"},
			Actions:              []string{"native.storage.bucket.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	response, err := r.nativeStub.Services.Storage.Bucket.Update(ctx.Request.Context(), &bucket.UpdateBucketRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		Name:      requestData.NewName,
		Hidden:    requestData.NewHidden,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid bucket name"})
				return
			}

			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "bucket not found"})
				return
			}
		}

		err := errors.New("failed to update bucket: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Bucket updated")
	ctx.JSON(http.StatusOK, updateBucketResponse{
		Bucket: FormatedBuceketFromGRPC(response.Bucket),
	})
}
