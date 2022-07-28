package initialization

import (
	"go-klikdokter/app/api/transport"
	"go-klikdokter/app/registry"
	"go-klikdokter/helper/_struct"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/database"
	"net/http"

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

	//Define auto migration here

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
	// Service registry
	ratingSvc := registry.RegisterRatingService(db, logger)
	publicRatingSvc := registry.RegisterPublicRatingService(db, logger)

	// Transport initialization
	swagHttp := transport.SwaggerHttpHandler(log.With(logger, "SwaggerTransportLayer", "HTTP")) //don't delete or change this !!
	ratingHttp := transport.RatingHttpHandler(ratingSvc, log.With(logger, "RatingTransportLayer", "HTTP"))
	publicRatingHttp := transport.PublicRatingHttpHandler(publicRatingSvc, log.With(logger, "PublicRatingTransportLayer", "HTTP"))

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //don't delete or change this!!
	mux.Handle(_struct.PrefixBase+"/public/", publicRatingHttp)
	mux.Handle(_struct.PrefixBase+"/", ratingHttp)

	return mux
}
