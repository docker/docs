using Docker.Core;

namespace Docker.Backend.Processes
{
    public interface IProxyProcess
    {
        void Start(Settings settings);
        void Start(Settings settings, string daemonNamedPipe);
        void Stop();
    }

    public class ApiProxy : IProxyProcess
    {
        private readonly IHyperV _hyperV;
        private readonly ExternalProcess _process;

        public ApiProxy(IHyperV hyperV)
        {
            _hyperV = hyperV;
            _process = new ExternalProcess(GetType(), "com.docker.proxy.exe", "com.docker.proxy");
        }

        public void Start(Settings settings)
        {
            _process.Start(() => $"-VM={_hyperV.GetId()}");
        }

        public void Start(Settings settings, string daemonNamedPipe)
        {
            _process.Start(() =>  $"-daemonNamedPipe=\\\\.\\pipe\\{daemonNamedPipe}");
        }

        public void Stop()
        {
            _process.Stop();
        }
    }
}