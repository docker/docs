using Docker.Core.backend;
using Docker.Core.Features;
using Docker.Core.Pipe;

namespace Docker.Core.Backend
{
    public class BackendClient : IBackend
    {
        private readonly NamedPipeClient _client;

        public BackendClient() : this("dockerBackend")
        {
        }

        internal BackendClient(string pipeName)
        {
            _client = new NamedPipeClient(pipeName);
        }

        public string Version()
        {
            return (string) SendMessage("Version");
        }

        public void Start(Settings settings)
        {
            SendMessage("Start", settings);
        }

        public void Stop()
        {
            SendMessage("Stop");
        }

        public void Destroy(bool keepVolume)
        {
            SendMessage("Destroy", keepVolume);
        }

        public string[] SharedDrives()
        {
            return (string[]) SendMessage("SharedDrives");
        }

        public void Unmount(string drive)
        {
            SendMessage("Unmount", drive);
        }

        public bool Mount(string drive, Credential credential, Settings settings)
        {
            return (bool) SendMessage("Mount", drive, credential, settings);
        }

        public void RemoveShare(string drive)
        {
            SendMessage("RemoveShare", drive);
        }

        public void MigrateVolume(string vmdkPath)
        {
            SendMessage("MigrateVolume", vmdkPath);
        }

        public string GetDebugInfo()
        {
            return (string) SendMessage("GetDebugInfo");
        }

        public string DownloadVmLogs()
        {
            return (string)SendMessage("DownloadVmLogs");
        }

        public void SaveSettingsToDb(Settings settings)
        {
            SendMessage("SaveSettingsToDb", settings);
        }

        public Feature[] InstallFeatures(Feature[] features)
        {
            return (Feature[])SendMessage("InstallFeatures", (object) features);
        }

        public void SwitchDaemon(Settings settings)
        {
            SendMessage("SwitchDaemon", settings);
        }

        private object SendMessage(string action, params object[] parameters)
        {
            return _client.Send(action, parameters);
        }
    }
}