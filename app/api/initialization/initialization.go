package initialization

import (
	"go-klikdokter/app/api/transport"
	publictransport "go-klikdokter/app/api/transport/public"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/_struct"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/database"
	"net/http"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-kit/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func DbInit() (*gorm.DB, error) {
	// Init DB Connection
	db, err := database.NewConnectionDB(config.GetConfigString(viper.GetString("database.driver")), config.GetConfigString(viper.GetString("database.dbname")),
		config.GetConfigString(viper.GetString("database.host")), config.GetConfigString(viper.GetString("database.username")), config.GetConfigString(viper.GetString("database.password")),
		config.GetConfigInt(viper.GetString("database.port")))
	if err != nil {
		return nil, err
	}

	// Define auto migration here

	// example Seeder
	// for i := 0; i < 1000; i++ {
	// 	fmt.Println("dijalankan")
	// 	product := entity.Product{}
	// 	err := faker.FakeData(&product)
	// 	db.Create(&product)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	return db, nil
}

func InitRouting(db *mongo.Database, logger log.Logger) *http.ServeMux {
	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) // don't delete or change this !!
	globalHttp := GlobalHttpHandler(log.With(logger, "GlobalTransportLayer", "HTTP"), db)

	// Routing path
	nsm := http.NewServeMux()
	nsm.Handle("/", swagHttp) // don't delete or change this!!
	nsm.HandleFunc("/__kdhealth", func(writer http.ResponseWriter, request *http.Request) { writer.Write([]byte("OK")) })
	nsm.Handle(_struct.PrefixBase+"/", globalHttp)

	return nsm
}

func GlobalHttpHandler(logger log.Logger, db *mongo.Database) http.Handler {
	// Service registry
	ratingSvc := registry.RegisterRatingService(db, logger)
	publicRatingSvc := registry.RegisterPublicRatingService(db, logger)
	publicRatingMpSvc := registry.RegisterPublicRatingMpService(db, logger)
	daprSvc := registry.RegisterDaprService(db, logger)
	// ratingMpSvc := registry.RegisterRatingMpService(db, logger)
	updloadImgSvc := registry.RegisterUploadService(db, logger)

	pr := mux.NewRouter()

	// ratingMpHttp := transport.RatingMpHttpHandler(ratingMpSvc, log.With(logger, "RatingMpTransportLayer", "HTTP"))
	ratingHttp := transport.RatingHttpHandler(ratingSvc, log.With(logger, "RatingTransportLayer", "HTTP"), db)
	publicRatingMpHttp := publictransport.PublicRatingMpHttpHandler(publicRatingMpSvc, log.With(logger, "PublicRatingMpTransportLayer", "HTTP"))
	publicRatingHttp := publictransport.PublicRatingHttpHandler(publicRatingSvc, log.With(logger, "PublicRatingTransportLayer", "HTTP"), db)
	daprHttp := transport.DaprHttpHandler(daprSvc, log.With(logger, "DaprTransportLayer", "HTTP"))
	uploadHttp := transport.UploadHttpHandler(updloadImgSvc, log.With(logger, "UploadTransportLayer", "HTTP"))

	pr.PathPrefix(_struct.PrefixBase + "/public/rating-submissions-by-id").Handler(publicRatingMpHttp)
	pr.PathPrefix(_struct.PrefixBase + "/public/ratings-summary/store-product").Handler(publicRatingMpHttp)
	// pr.PathPrefix(_struct.PrefixBase + "/public/ratings-summary-mp").Handler(publicRatingMpHttp)
	pr.PathPrefix(_struct.PrefixBase + "/public/rating-submissions").Handler(publicRatingHttp)
	pr.PathPrefix(_struct.PrefixBase + "/public/ratings-summary").Handler(publicRatingHttp)
	pr.PathPrefix(_struct.PrefixBase + "/dapr").Handler(daprHttp)
	pr.PathPrefix(_struct.PrefixBase + "/upload/").Handler(uploadHttp) // for upload images
	// pr.PathPrefix(_struct.PrefixBase + "/rating-submissions-mp").Handler(ratingMpHttp)
	// pr.PathPrefix(_struct.PrefixBase + "/ratings-summary-mp").Handler(ratingMpHttp)
	pr.PathPrefix(_struct.PrefixBase + "/").Handler(ratingHttp)

	return pr
}
