package registry

import (
	rp "go-klikdokter/app/repository"
	publicrepository "go-klikdokter/app/repository/public"
	"go-klikdokter/app/service"
	"go-klikdokter/app/service/public"
	"go-klikdokter/pkg/util"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/log"
)

func RegisterRatingService(db *mongo.Database, logger log.Logger) service.RatingService {
	return service.NewRatingService(
		logger,
		rp.NewRatingRepository(db),
		publicrepository.NewPublicRatingRepository(db),
		util.NewMedicalFacilitySvc(util.ResponseHttp{}),
	)
}

func RegisterPublicRatingService(db *mongo.Database, logger log.Logger) publicservice.PublicRatingService {
	return publicservice.NewPublicRatingService(
		logger,
		rp.NewRatingRepository(db),
		publicrepository.NewPublicRatingRepository(db),
	)
}

func RegisterDaprService(db *mongo.Database, logger log.Logger) service.DaprService {
	return service.NewDaprService(
		logger,
	)
}

func RegisterRatingMpService(db *mongo.Database, logger log.Logger) service.RatingMpService {
	return service.NewRatingMpService(
		logger,
		rp.NewRatingMpRepository(db),
	)
}

func RegisterPublicRatingMpService(db *mongo.Database, logger log.Logger) publicservice.PublicRatingMpService {
	return publicservice.NewPublicRatingMpService(
		logger,
		rp.NewRatingMpRepository(db),
		publicrepository.NewPublicRatingMpRepository(db),
	)
}
