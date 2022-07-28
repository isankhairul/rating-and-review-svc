//  KD Boilerplate Go kit:
//   version: 1.0
//   title: Boilerplate Go Kit Api
//  Schemes: https, http
//  Host:
//  BasePath: /rating-svc/api/v1
//  Produces:
//    - application/json
//  securityDefinitions:
//   Bearer:
//    type: apiKey
//    name: Authorization
//    in: header
//
// swagger:meta
package main

import (
	"fmt"
	"go-klikdokter/app/api/initialization"
	"go-klikdokter/helper/_struct"
	"go-klikdokter/helper/config"
	"go-klikdokter/helper/consul"
	"go-klikdokter/helper/database"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	//Load configuration
	viper.SetConfigType("yaml")
	var profile string = "dev"
	if os.Getenv("KD_ENV") == "prd" || os.Getenv("KD_ENV") == "stg" {
		profile = "prd"
	}

	var configFileName []string
	configFileName = append(configFileName, "config-")
	configFileName = append(configFileName, profile)
	viper.SetConfigName(strings.Join(configFileName, ""))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}

	//override with env vars
	viper.AutomaticEnv()
	viper.SetEnvPrefix("KD")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Logging init
	logfile, err := os.OpenFile(viper.GetString("server.output-file-path"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	defer logfile.Close()
	var logger log.Logger
	{
		if config.GetConfigString(viper.GetString("server.log-output")) == "file" {
			w := log.NewSyncWriter(logfile)
			logger = log.NewLogfmtLogger(w)
		} else {
			logger = log.NewLogfmtLogger(os.Stderr)
		}

		logger = level.NewFilter(logger, level.AllowAll())
		logger = level.NewInjector(logger, level.InfoValue())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}

	// Init MongoDB Connection
	db, err := database.NewMongo()
	if err != nil {
		_ = logger.Log("Err Db connection :", err.Error())
		panic(err.Error())
	}
	// Init DB Connection
	//db, err := initialization.DbInit()
	if err != nil {
		_ = logger.Log("Err Db connection :", err.Error())
		panic(err.Error())
	}
	_ = logger.Log("message", "Connection Db Success")

	// Consul initialization
	registar := consul.ConsulRegisterService(config.GetConfigString(viper.GetString("server.service-name")), config.GetConfigInt(viper.GetString("server.port")), logger)
	registar.Register()
	defer registar.Deregister()

	// Routing initialization
	_struct.PrefixBase = viper.GetString("route.site") + viper.GetString("route.apiprefix")

	mux := initialization.InitRouting(db, logger)
	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)

	// // configure hystrix
	// var prescriptionEndpoint endpoint.Endpoint
	// hystrix.ConfigureCommand("prescription Request", hystrix.CommandConfig{Timeout: 1000})
	// prescriptionEndpoint = Hystrix("Prescription Request", "Service currently unavailable", logger)(prescriptionEndpoint)
	// hystrixStreamHandler := hystrix.NewStreamHandler()
	// hystrixStreamHandler.Start()
	// go func() {
	// 	errs <- http.ListenAndServe(net.JoinHostPort("", "9000"), hystrixStreamHandler)
	// }()

	go func() {
		_ = logger.Log("transport", "HTTP", "addr", config.GetConfigInt(viper.GetString("server.port")))
		errs <- http.ListenAndServe(fmt.Sprintf(":%d", config.GetConfigInt(viper.GetString("server.port"))), nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("exit", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
