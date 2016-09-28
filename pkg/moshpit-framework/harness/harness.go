package harness

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/commands"
)

type moshpitServer struct {
	log          logrus.FieldLogger
	clientConfig string
}

func (ms *moshpitServer) Register(stream commands.Mosher_RegisterServer) error {
	jobSent := false
	// this function handles a single client
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// print out jobs from client
		ms.log.WithField("jobs", in.Jobs).WithField("node", in.Client).Debug("jobs")
		if len(in.Jobs) > 0 {
			job := in.Jobs[0]
			if job.State != commands.JobState_RUNNING {
				// kill the client as soon as it finishes with success or failure
				if job.State == commands.JobState_FAILURE {
					ms.log.WithFields(map[string]interface{}{
						"state": commands.JobState_name[int32(job.State)],
						"error": job.Error,
						"node":  in.Client,
					}).Error("job failed on node")
				} else if job.State == commands.JobState_SUCCESS {
					ms.log.WithFields(map[string]interface{}{
						"state": commands.JobState_name[int32(job.State)],
						"node":  in.Client,
					}).Info("job succeeded on node")
				} else {
					ms.log.WithFields(map[string]interface{}{
						"state": commands.JobState_name[int32(job.State)],
						"node":  in.Client,
					}).Info("job ended on node with unknown state")
				}
				return nil
			}
		} else if !jobSent {
			ms.log.WithField("client", in.Client).Info("sending job to client")
			// create a job on this client
			err = stream.Send(&commands.JobConfig{
				// TODO: autogenerate it and track it
				Uuid:       "abc",
				MoshConfig: ms.clientConfig,
			})
			if err != nil {
				return err
			}
			jobSent = true
		}
	}
}

func Server(ctx context.Context, config moshpit.ServerConfig, setupFunc moshpit.SetupFunc, clientRunFunc moshpit.ClientRunFunc) error {
	log := moshpit.LoggerFromCtx(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.ListenIP, config.ListenPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	// TODO: use tls
	//if *tls {
	//	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	//	if err != nil {
	//		grpclog.Fatalf("Failed to generate credentials %v", err)
	//	}
	//	opts = []grpc.ServerOption{grpc.Creds(creds)}
	//}
	grpcServer := grpc.NewServer(opts...)

	clientConfig, err := yaml.Marshal(config.Client)
	if err != nil {
		return err
	}

	commands.RegisterMosherServer(grpcServer, &moshpitServer{
		clientConfig: string(clientConfig),
		log:          log,
	})

	// Before starting the server we execute the setup actions:
	// TODO: maybe do this after starting the sever?
	setupConfig, err := yaml.Marshal(config.Setup)
	if err != nil {
		return err
	}
	err = setupFunc(ctx, string(setupConfig))
	if err != nil {
		return err
	}

	return grpcServer.Serve(lis)
}

type jobInfo struct {
	started time.Time
	uuid    string
	job     moshpit.Job
}

func Client(ctx context.Context, server, name string, clientRunFunc moshpit.ClientRunFunc) error {
	log := moshpit.LoggerFromCtx(ctx)
	// TODO: remove insecure
	log.WithField("server", server).Debug("connecting to server")
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := commands.NewMosherClient(conn)
	log.Debug("register")
	registerClient, err := client.Register(context.Background())
	if err != nil {
		return err
	}

	jobsLock := sync.Mutex{}
	jobs := &[]jobInfo{}
	// wait for new jobs in a goroutine
	go func() {
		err := func() error {
			for {
				log.Debug("receiving job")
				config, err := registerClient.Recv()
				if err == io.EOF {
					return nil
				}
				if err != nil {
					if err.Error() == `rpc error: code = 1 desc = "context canceled"` {
						return nil
					}
					return err
				}
				log.WithField("jobs", jobs).Debug("received job")
				clientJob, err := clientRunFunc(ctx, name, config.MoshConfig)
				if err != nil {
					return err
				}
				jobsLock.Lock()
				*jobs = append(*jobs, jobInfo{
					uuid:    config.Uuid,
					started: time.Now(),
					job:     clientJob,
				})
				jobsLock.Unlock()
			}
		}()
		if err != nil {
			log.WithField("error", err).Warn("error receiving")
		}
	}()

	for {
		status := &commands.Status{
			Client: name,
		}

		jobsLock.Lock()
		for _, job := range *jobs {
			state, err := job.job.State()
			errString := ""
			if err != nil {
				errString = err.Error()
			}
			status.Jobs = append(status.Jobs, &commands.JobInstance{
				Uuid:   job.uuid,
				State:  state,
				Error:  errString,
				Length: uint64(time.Now().Sub(job.started).Seconds()),
			})
		}
		jobsLock.Unlock()

		log.WithField("jobs", *jobs).WithField("status", status).Debug("sending status")

		err := registerClient.Send(status)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			if err.Error() == `rpc error: code = 1 desc = "context canceled"` {
				return nil
			}
			return err
		}

		// don't report status too often
		time.Sleep(time.Second)
	}
	return nil
}
