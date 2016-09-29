using Docker.Core;

namespace Docker.Backend.Processes
{
    public interface IVpnKit
    {
        void Start(Settings settings);
        void Stop();
    }

    public class VpnKit : IVpnKit
    {
        private readonly IHyperV _hyperV;
        private readonly ExternalProcess _process;

        public VpnKit(IHyperV hyperV)
        {
            _hyperV = hyperV;
            _process = new ExternalProcess(GetType(), "com.docker.slirp.exe", "com.docker.slirp");
        }

        public void Start(Settings settings)
        {
            _process.Start(() =>
            {
                var vmId = _hyperV.GetId();

                return $"--ethernet hyperv-connect://{vmId} --port hyperv-connect://{vmId} --db \\\\.\\pipe\\dockerDataBase --debug";
            });
        }

        public void Stop()
        {
            _process.Stop();
        }
    }
}