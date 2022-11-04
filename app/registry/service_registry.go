package registry

import (
	rp "go-klikdokter/app/repository"
	"go-klikdokter/app/service"
	"go-klikdokter/pkg/util"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/log"
)

func RegisterRatingService(db *mongo.Database, logger log.Logger) service.RatingService {
	return service.NewRatingService(
		logger,
		rp.NewRatingRepository(db),
		rp.NewPublicRatingRepository(db),
		util.NewMedicalFacilitySvc(util.ResponseHttp{}),
	)
}

func RegisterPublicRatingService(db *mongo.Database, logger log.Logger) service.PublicRatingService {
	return service.NewPublicRatingService(
		logger,
		rp.NewRatingRepository(db),
		rp.NewPublicRatingRepository(db),
	)
}

func RegisterDaprService(db *mongo.Database, logger log.Logger) service.DaprService {
	return service.NewDaprService(
		logger,
	)
}
