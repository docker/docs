using Docker.Backend;
using Docker.Core;
using System;
using System.Globalization;
using System.IO;
using System.ServiceProcess;
using System.Threading;
using System.Threading.Tasks;
using Docker.Core.di;
using Docker.Core.Pipe;

namespace Docker.Service
{
    public partial class DockerService : ServiceBase
    {
        private BackendServer _backendServer;

        public DockerService()
        {
            InitializeComponent();
        }

        public void DoRun()
        {
            Thread.CurrentThread.CurrentCulture = new CultureInfo("en-US");

            Logger.Initialize(Path.Combine(Paths.CommonApplicationData, "service.txt"));

            var logger = new Logger(GetType());
            logger.Info("Version: " + new Core.Version().ToHumanString());
            logger.Info("Starting on: " + DateTime.Now);
            logger.Info("Sha1: " + new Git().Sha1());

            var logPipeServer = new LogPipeServer("dockerLogs");
            Logger.SetListener(logPipeServer);
            Task.Run(() => logPipeServer.Run());

            var singletons = new Singletons(new BackendServerModule());
            _backendServer = singletons.Get<BackendServer>();
            _backendServer.Run();
        }

        public void DoStop()
        {
            _backendServer?.Stop();
       }

        protected override void OnStart(string[] args)
        {
            DoRun();
        }

        protected override void OnStop()
        {
            DoStop();
            Environment.Exit(0);
        }
    }
}
