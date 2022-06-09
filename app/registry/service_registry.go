package registry

import (
	rp "go-klikdokter/app/repository"
	"go-klikdokter/app/service"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/log"
)

func RegisterRatingService(db *mongo.Database, logger log.Logger) service.RatingService {
	return service.NewRatingService(
		logger,
		rp.NewRatingRepository(db),
	)
}
