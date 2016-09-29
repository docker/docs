using System.Diagnostics;
using System.Threading.Tasks;
using Docker.Backend.Processes;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.di;
using Docker.Core.Features;
using Docker.Core.Pipe;

namespace Docker.Backend
{
    public class BackendServerModule : Module
    {
        protected override void Configure()
        {
            Bind<IBackend, Backend>();
        }
    }

    public class BackendServer
    {
        private readonly Logger _logger;
        private readonly IBackend _backend;
        private readonly NamedPipeServer _namedPipeServer;

        public BackendServer(IBackend backend, IClientIdentityCallback clientIdentityCallback)
        {
            _logger = new Logger(GetType());
            _backend = backend;
            _namedPipeServer = new NamedPipeServer("dockerBackend", clientIdentityCallback);
        }

        public void Run()
        {
            new Job().AddProcess(Process.GetCurrentProcess().Handle);

            _namedPipeServer.Register("Version", args => _backend.Version());
            _namedPipeServer.Register("Start", args => _backend.Start((Settings)args[0]));
            _namedPipeServer.Register("Stop", args => _backend.Stop());
            _namedPipeServer.Register("Destroy", args => _backend.Destroy((bool)args[0]));
            _namedPipeServer.Register("SharedDrives", args => _backend.SharedDrives());
            _namedPipeServer.Register("Unmount", args => _backend.Unmount((string)args[0]));
            _namedPipeServer.Register("Mount", args => _backend.Mount((string)args[0], (Credential)args[1], (Settings)args[2]));
            _namedPipeServer.Register("RemoveShare", args => _backend.RemoveShare((string)args[0]));
            _namedPipeServer.Register("MigrateVolume", args => _backend.MigrateVolume((string)args[0]));
            _namedPipeServer.Register("GetDebugInfo", args => _backend.GetDebugInfo());
            _namedPipeServer.Register("DownloadVmLogs", args => _backend.DownloadVmLogs());
            _namedPipeServer.Register("InstallFeatures", args => _backend.InstallFeatures((Feature[])args[0]));
            _namedPipeServer.Register("SwitchDaemon", args => _backend.SwitchDaemon((Settings)args[0]));

            Task.Run(() => _namedPipeServer.Run());
            _logger.Info("Started");
        }

        public void Stop()
        {
            _logger.Info("Stopping...");
            _namedPipeServer.Stop();
            _logger.Info("Stopped");
        }
    }
}