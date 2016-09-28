package main

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"os"
	"path"

	apiserver "github.com/docker/dhe-deploy/adminserver/api/server"
	"github.com/docker/distribution/context"
	"github.com/emicklei/go-restful/swagger"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	file, err := os.Create("/out/docs.tar")
	handleErr(err)

	tarWriter := tar.NewWriter(file)

	apiServer := apiserver.NewAPIServer(context.Background(), nil, nil, nil, nil, nil)
	config, _ := apiServer.BuildSubroutes("/api", "apidocgen")
	sb := swagger.NewSwaggerBuilder(*config)

	listing := sb.ProduceListing()
	str, err := json.Marshal(listing)
	handleErr(err)
	err = tarWriter.WriteHeader(&tar.Header{Name: "docs.json", Mode: 0666, Size: int64(len(str))})
	handleErr(err)
	_, err = tarWriter.Write(str)
	handleErr(err)

	decls := sb.ProduceAllDeclarations()
	for declPath, decl := range decls {
		str, err := json.Marshal(decl)
		handleErr(err)
		err = tarWriter.WriteHeader(&tar.Header{Name: path.Join("docs", declPath), Mode: 0666, Size: int64(len(str))})
		handleErr(err)
		_, err = tarWriter.Write(str)
		handleErr(err)
	}
	err = tarWriter.Flush()
	handleErr(err)
}
