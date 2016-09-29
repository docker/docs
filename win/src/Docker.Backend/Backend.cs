using System;
using Docker.Core;
using Docker.Core.Features;
using Docker.Backend.ContainerEngine;
using Docker.Backend.Features;
using Docker.Core.backend;

namespace Docker.Backend
{
    public class Backend : IBackend
    {
        private readonly IVersion _version;
        private readonly ISambaShare _sambaShare;
        private readonly IDockerMachineImport _dockerMachineImport;
        private readonly ICmd _cmd;
        private readonly IFiles _files;
        private readonly IContainerEngineHelper _containerEngineHelper;
        private readonly ILinux _linuxContainerEngine;
        private readonly IWindows _windowsContainerEngine;
        private readonly IHyperV _hyperV;
        private readonly IInstaller _featuresInstaller;

        public Backend(
            IVersion version, ISambaShare sambaShare,
            IDockerMachineImport dockerMachineImport, ICmd cmd, IFiles files,
            IContainerEngineHelper containerEngineHelper, IWindows windowsContainerEngine,
            ILinux linuxContainerEngine, IHyperV hyperV, IInstaller featuresInstaller)
        {
            _version = version;
            _sambaShare = sambaShare;
            _dockerMachineImport = dockerMachineImport;
            _cmd = cmd;
            _files = files;
            _containerEngineHelper = containerEngineHelper;
            _linuxContainerEngine = linuxContainerEngine;
            _windowsContainerEngine = windowsContainerEngine;
            _hyperV = hyperV;
            _featuresInstaller = featuresInstaller;
        }

        public string Version()
        {
            return _version.ToHumanString();
        }

        public void Start(Settings settings)
        {
            if (_containerEngineHelper.UseLinuxContainerEngine)
            {
                _linuxContainerEngine.Start(settings);
            }
            else
            {
                try
                {
                    _windowsContainerEngine.Start(settings);
                }
                catch (Exception)
                {
                    _containerEngineHelper.Switch();
                    throw;
                }
            }
        }

        public void Stop()
        {
            _linuxContainerEngine.Stop();
            _windowsContainerEngine.Stop();
        }

        public void Destroy(bool keepVolume)
        {
            _linuxContainerEngine.Destroy(keepVolume);
            _windowsContainerEngine.Destroy(keepVolume);
        }

        public string[] SharedDrives()
        {
            return _sambaShare.SharedDrives();
        }

        public bool Mount(string drive, Credential credential, Settings settings)
        {
            if (!_containerEngineHelper.UseLinuxContainerEngine) return true;

            return _sambaShare.Mount(drive, credential, settings);
        }

        public void Unmount(string drive)
        {
            if (!_containerEngineHelper.UseLinuxContainerEngine) return;

            _sambaShare.Unmount(drive);
        }

        public void RemoveShare(string drive)
        {
            _sambaShare.DeleteShare(drive);
        }

        public void MigrateVolume(string vmdkPath)
        {
            _dockerMachineImport.MigrateVolume(vmdkPath, Paths.TmpVolumeMigrationPath);
            _linuxContainerEngine.Stop();
            _linuxContainerEngine.Destroy(false);
            _files.ForceMove(Paths.TmpVolumeMigrationPath, Paths.MobyDiskPath);
        }

        public string GetDebugInfo()
        {
            var result = _cmd.RunAsAdministrator("powershell.exe",
                $"-ExecutionPolicy UNRESTRICTED -NoProfile -NonInteractive -command \"& {{& '{Paths.InResourcesDir("DockerDebugInfo.ps1")}'}}\"", 0);
            return result.CombinedOutput;
        }

        public string DownloadVmLogs()
        {
            return _hyperV.DownloadLogs();
        }

        public Feature[] InstallFeatures(Feature[] features)
        {
            return _featuresInstaller.Install(features);
        }

        public void SwitchDaemon(Settings settings)
        {
            _containerEngineHelper.Switch();
            Start(settings);
        }
    }
}
