package handler

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"hashServer/internal/generated/models"
	"hashServer/internal/generated/restapi/operations"
	"hashServer/internal/repository"
	"hashServer/pkg/hashService"
	"hashServer/pkg/logger"
	"time"
)

type Handler struct {
	repository *repository.Repository
}

var Mux *Handler

func InitHandler(r *repository.Repository) {
	Mux = &Handler{
		repository: r,
	}
}

func (h *Handler) SaveHashesFromString(params operations.PostSendParams) middleware.Responder {
	cwt, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := grpc.DialContext(cwt, fmt.Sprintf(":%s", viper.GetString("hash_service_port")),
		grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		logger.Error(context.Background(), "SaveHashesFromString", "handler", err, "grpc connection failed")
	}
	defer conn.Close()

	hash := hashService.NewHashServiceClient(conn)

	hashResult, err := hash.GetHash(cwt, &hashService.ListOfStrings{
		Strings: params.Params,
	})
	if err != nil {
		logger.Error(context.Background(), "SaveHashesFromString", "handler", err, "grpc result failed")
	}

	result := models.ArrayOfHash{}

	for _, item := range hashResult.Data {
		itemId, err := h.repository.Hash.Create(item.Hash)
		if err != nil {
			logger.Error(context.Background(), "SaveHashesFromString", "handler", err, "insert id failed")
		}

		id := int64(itemId)
		str := item.Hash
		result = append(result, &models.Hash{ID: &id, Hash: &str})
	}

	return operations.NewPostSendOK().WithPayload(result)
}

func (h *Handler) GetHashesByIds(params operations.GetCheckParams) middleware.Responder {
	result, err := h.repository.Hash.GetByIds(params.Ids)

	if err != nil {
		logger.Error(context.Background(), "GetHashesByIds", "handler", err,
			fmt.Sprintf("failed get hash by id - %s", params.Ids))
		return operations.NewGetCheckInternalServerError()
	}

	return operations.NewGetCheckOK().WithPayload(result)
}
