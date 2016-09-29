using System;
using System.ServiceProcess;
using Docker.Backend.Features;
using Docker.Backend.Processes;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.Features;
using Microsoft.Win32;

namespace Docker.Backend.ContainerEngine
{
    public interface IWindows : IContainerEngine
    {
    }

    public class Windows : IWindows
    {
        private readonly Feature[] _features = { Feature.HyperV, Feature.Containers };
        private readonly IWindowsDockerDaemon _windowsDockerDaemon;
        private readonly IProxyProcess _proxy;
        private readonly IDockerDaemonChecker _daemonChecker;
        private readonly IInstaller _featuresInstaller;
        private readonly IContainerEngineHelper _containerEngineHelper;

        private ContainerEngineStatus _status;

        public Windows(IWindowsDockerDaemon windowsDockerDaemon, IProxyProcess proxy, IDockerDaemonChecker daemonChecker, IInstaller featuresInstaller, IContainerEngineHelper containerEngineHelper)
        {
            _windowsDockerDaemon = windowsDockerDaemon;
            _proxy = proxy;
            _daemonChecker = daemonChecker;
            _featuresInstaller = featuresInstaller;
            _containerEngineHelper = containerEngineHelper;

            _status = ContainerEngineStatus.Stopped;
        }

        public void Start(Settings settings)
        {
            CheckDockerDaemonService();

            try
            {
                CheckInstallation();
                DoStart(settings);
            }
            catch
            {
                _status = ContainerEngineStatus.FailedToStart;
                throw;
            }
        }

        private void DoStart(Settings settings)
        {
            if (_status != ContainerEngineStatus.Started)
            {
                _status = ContainerEngineStatus.Starting;

                _containerEngineHelper.ForceKillLingeringDaemon();
                _windowsDockerDaemon.Start();
            }

            _proxy.Start(settings, Paths.DockerDaemonNamedPipe);
            _daemonChecker.Check();
            _status = ContainerEngineStatus.Started;
        }

        private void CheckInstallation()
        {
            _featuresInstaller.CheckInstalledFeatures(_features);

            using (var key = Registry.LocalMachine.OpenSubKey(@"SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\Containers", true))
            {
                key?.SetValue("VSmbDisableOplocks", 1);
            }
        }

        private static void CheckDockerDaemonService()
        {
            try
            {
                var sc = new ServiceController { ServiceName = "docker" };
                if (sc.Status == ServiceControllerStatus.Running)
                {
                    sc.Stop();
                }
            }
            catch
            {
                // Service is not installed
                return;
            }

            throw new BackendWarnException("We detected a previous installation of Windows Docker daemon\nPlease unregister it using: dockerd --unregister-service");
        }

        public void Stop()
        {
            if (_status == ContainerEngineStatus.Stopped)
                return;

            _status = ContainerEngineStatus.Stopping;
            try
            {
                _proxy.Stop();
                _windowsDockerDaemon.Stop();
            }
            catch(Exception)
            {
                _status = ContainerEngineStatus.FailedToStop;
                throw;
            }
            _status = ContainerEngineStatus.Stopped;
        }

        public void Destroy(bool keepVolume)
        {
        }
    }
}
