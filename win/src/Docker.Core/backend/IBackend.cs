using Docker.Core.Features;

namespace Docker.Core.backend
{
    public interface IBackend
    {
        string Version();
        void Start(Settings settings);
        void Stop();
        void Destroy(bool keepVolume);
        string[] SharedDrives();
        void Unmount(string drive);
        bool Mount(string drive, Credential credential, Settings settings);
        void RemoveShare(string drive);
        void MigrateVolume(string vmdkPath);
        string GetDebugInfo();
        string DownloadVmLogs();
        Feature[] InstallFeatures(Feature[] features);
        void SwitchDaemon(Settings settings);
    }
}