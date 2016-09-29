using Docker.Core;

namespace Docker.Backend.Processes
{
    public interface IDataKit
    {
        void Start(Settings settings);
        void Stop();
    }

    public class DataKit : IDataKit
    {
        private readonly IHyperV _hyperV;
        private readonly ExternalProcess _process;

        public DataKit(IHyperV hyperV)
        {
            _hyperV = hyperV;
            _process = new ExternalProcess(GetType(), "com.docker.db.exe", "com.docker.db");
        }

        public void Start(Settings settings)
        {
            _process.Start(() =>
            {
                var vmId = _hyperV.GetId();

                return $"--url hyperv-connect://{vmId},\\\\.\\pipe\\dockerDataBase --color never";
            });
        }

        public void Stop()
        {
            _process.Stop();
        }
    }
}