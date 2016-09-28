package cmd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/server"
	"github.com/docker/orca/enzi/api/server/accounts"
	"github.com/docker/orca/enzi/api/server/admin"
	"github.com/docker/orca/enzi/api/server/config"
	"github.com/docker/orca/enzi/api/server/jobs"
	"github.com/docker/orca/enzi/api/server/openid"
	"github.com/docker/orca/enzi/api/server/workers"
	"github.com/docker/orca/enzi/schema"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

// APIServer is the command for running an API server.
var APIServer = cli.Command{
	Name:   "api",
	Usage:  "Run an API Server",
	Action: runAPI,
}

const (
	indexSourcePath  = "/ui/index.html"
	jQuerySourcePath = "/ui/jquery.min.js"
	jQueryRoute      = "static/jquery.min.js"
)

var (
	rootPrefix string
	enableDocs bool
)

func init() {
	APIServer.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "root-prefix, r",
			Usage:       "prefix all HTTP routes with this path",
			Destination: &rootPrefix,
		},
		cli.BoolFlag{
			Name:        "enable-docs",
			Usage:       "whether to enable the interactive API documentation",
			Destination: &enableDocs,
		},
	}
}

func runAPI(*cli.Context) error {
	indexServer, err := newIndexHTMLServer(rootPrefix)
	if err != nil {
		log.Fatalf("unable to create index.html server: %s", err)
	}

	jQueryServer, err := newStaticFileServer(jQuerySourcePath)
	if err != nil {
		log.Fatalf("unable to create %s server: %s", jQueryRoute, err)
	}

	tlsConfig := GetTLSConfig(tls.NoClientCert)

	log.Println("connecting to db ...")
	dbSession := GetDBSession(tlsConfig)
	defer dbSession.Close()

	schemaMgr := schema.NewRethinkDBManager(dbSession)
	workerClient := server.NewHTTPClient(tlsConfig)

	serviceContainer := restful.NewContainer()

	log.Println("generating private key ...")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("unable to generate RSA private key: %s", err)
	}

	log.Println("initializing services ...")

	idService, err := openid.NewService(context.Background(), schemaMgr, "/tls/ca.pem", privateKey, time.Hour*6, path.Join("/", rootPrefix, "v0/id"))
	if err != nil {
		log.Fatalf("unable to initialize OpenID Connect Provider service: %s", err)
	}

	accountsService := accounts.NewService(context.Background(), schemaMgr, path.Join("/", rootPrefix, "v0/accounts"))
	adminService := admin.NewService(context.Background(), schemaMgr, path.Join("/", rootPrefix, "v0/admin"))
	configService := config.NewService(context.Background(), schemaMgr, path.Join("/", rootPrefix, "v0/config"))
	workersService := workers.NewService(context.Background(), schemaMgr, workerClient, path.Join("/", rootPrefix, "v0/workers"))

	jobsService, err := jobs.NewService(context.Background(), schemaMgr, workerClient, path.Join("/", rootPrefix, "v0/jobs"))
	if err != nil {
		log.Fatalf("unable to create jobs service: %s", err)
	}

	serviceContainer.
		Add(idService.WebService).
		Add(accountsService.WebService).
		Add(adminService.WebService).
		Add(configService.WebService).
		Add(workersService.WebService).
		Add(jobsService.WebService)

	if enableDocs {
		swaggerConfig := swagger.Config{
			WebServices: []*restful.WebService{
				idService.WebService,
				configService.WebService,
				accountsService.WebService,
				workersService.WebService,
				jobsService.WebService,
			},
			ApiPath:         path.Join("/", rootPrefix, "v0/docs/docs.json"),
			SwaggerPath:     path.Join("/", rootPrefix, "v0/docs") + "/",
			SwaggerFilePath: "/ui/docs",
			Info:            swagger.Info{Title: "eNZi API v0 Documentation"},
		}

		swagger.RegisterSwaggerService(swaggerConfig, serviceContainer)
	}

	serveMux := http.NewServeMux()
	serveMux.Handle(path.Join("/", rootPrefix, "_ping"), http.HandlerFunc(pingHandler))
	serveMux.Handle(path.Join("/", rootPrefix, "v0")+"/", serviceContainer)
	serveMux.Handle(path.Join("/", rootPrefix, jQueryRoute), jQueryServer)
	serveMux.Handle("/", indexServer)

	server := &http.Server{
		Addr:      ":4443",
		Handler:   serveMux,
		TLSConfig: tlsConfig,
	}

	log.Println("listening for connections ...")
	log.Fatal(server.ListenAndServeTLS("", ""))

	return nil
}

type indexTemplateData struct {
	JQuerySource            string
	DefaultRedirectLocation string

	LoginPage     string
	LogoutPage    string
	AuthorizePage string

	LoginEndpoint     string
	LogoutEndpoint    string
	AuthorizeEndpoint string
}

type staticFileServer struct {
	modtime  time.Time
	baseName string
	buffer   *bytes.Buffer
}

func newStaticFileServer(sourcePath string) (http.Handler, error) {
	sourceBytes, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read %s: %s", sourcePath, err)
	}

	return &staticFileServer{
		modtime:  time.Now(),
		baseName: path.Base(sourcePath),
		buffer:   bytes.NewBuffer(sourceBytes),
	}, nil
}

func (s *staticFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, s.baseName, s.modtime, bytes.NewReader(s.buffer.Bytes()))
}

func newIndexHTMLServer(rootPrefix string) (http.Handler, error) {
	indexTemplate, err := template.ParseFiles(indexSourcePath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse index.html template: %s", err)
	}

	// By default, redirect to the logout page after login if 'next' page
	// is not specified.
	defaultRedirectLocation := path.Join("/", rootPrefix, "logout")
	if enableDocs {
		// If the interactive docs are enabled, make the documentaiton
		// the default redirect after login.
		defaultRedirectLocation = path.Join("/", rootPrefix, "v0/docs") + "/"
	}

	templateData := indexTemplateData{
		JQuerySource:            path.Join("/", rootPrefix, jQueryRoute),
		DefaultRedirectLocation: defaultRedirectLocation,
		LoginPage:               path.Join("/", rootPrefix, "login"),
		LogoutPage:              path.Join("/", rootPrefix, "logout"),
		AuthorizePage:           path.Join("/", rootPrefix, "authorize"),
		LoginEndpoint:           path.Join("/", rootPrefix, "v0/id/login"),
		LogoutEndpoint:          path.Join("/", rootPrefix, "v0/id/logout"),
		AuthorizeEndpoint:       path.Join("/", rootPrefix, "v0/id/authorize"),
	}

	buffer := new(bytes.Buffer)
	if err := indexTemplate.Execute(buffer, templateData); err != nil {
		return nil, fmt.Errorf("unable to execute template: %s", err)
	}

	return &staticFileServer{
		modtime:  time.Now(),
		baseName: "index.html",
		buffer:   buffer,
	}, nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	// For now just always return OK - in the future consider
	// adding logic for very rudimentary health checking
	w.Header().Set("content-type", "text/plain")
	w.Write([]byte("OK"))
}
