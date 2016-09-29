using System;
using Docker.Backend.Features;
using Docker.Backend.Processes;
using Docker.Core;
using Docker.Core.Features;
using System.ServiceProcess;

namespace Docker.Backend.ContainerEngine
{
    public interface ILinux : IContainerEngine
    {
    }

    public class Linux : ILinux
    {
        private readonly Feature[] _features = { Feature.HyperV };
        private readonly IHyperV _hyperV;
        private readonly IFirewall _firewall;
        private readonly IVpnKit _vpnKit;
        private readonly IDataKit _dataKit;
        private readonly IProxyProcess _proxy;
        private readonly IHyperVGuids _hyperVGuids;
        private readonly IDatabase _database;
        private readonly IDnsUpdater _dnsUpdater;
        private readonly IDockerDaemonChecker _daemonChecker;
        private readonly IInstaller _featuresInstaller;
        private readonly IContainerEngineHelper _containerEngineHelper;

        internal ContainerEngineStatus Status;

        public Linux(IHyperV hyperV,
            IFirewall firewall, IVpnKit vpnKit, IDataKit dataKit, IProxyProcess proxy,
            IHyperVGuids hyperVGuids, IDatabase database, IDnsUpdater dnsUpdater, IDockerDaemonChecker daemonChecker, IInstaller featuresInstaller,
            IContainerEngineHelper containerEngineHelper)
        {
            _hyperV = hyperV;
            _firewall = firewall;
            _vpnKit = vpnKit;
            _dataKit = dataKit;
            _proxy = proxy;
            _hyperVGuids = hyperVGuids;
            _database = database;
            _dnsUpdater = dnsUpdater;
            _daemonChecker = daemonChecker;
            _featuresInstaller = featuresInstaller;
            _containerEngineHelper = containerEngineHelper;

            Status = ContainerEngineStatus.Stopped;
        }

        public void Create(Settings settings)
        {
            _hyperV.Create(settings);
        }

        public void Start(Settings settings)
        {
            try
            {
                DoStart(settings);
            }
            catch
            {
                try
                {
                    CheckInstallation();
                    Create(settings);
                    DoStart(settings);
                }
                catch
                {
                    Status = ContainerEngineStatus.FailedToStart;
                    throw;
                }
            }
            Status = ContainerEngineStatus.Started;
        }

        private void DoStart(Settings settings)
        {
            if (Status != ContainerEngineStatus.Started)
            {
                Stop();

                Status = ContainerEngineStatus.Starting;
                CheckDockerDaemonService();
                _containerEngineHelper.ForceKillLingeringDaemon();
                _hyperVGuids.Install();
                _firewall.OpenPorts();
                _hyperV.Create(settings);
                _dataKit.Start(settings);
                _vpnKit.Start(settings);
                _database.Write(settings);
                _dnsUpdater.StaticDnsAndSearchDomains = settings.UseDnsForwarder ? "" : $"nameserver {settings.NameServer}";
                _dnsUpdater.UpdateDnsAndSearchDomains();
                _hyperV.Start();
                _dnsUpdater.RegisterWatcher();
            }

            _proxy.Start(settings);
            _daemonChecker.Check();
        }

        public void Stop()
        {
            if (Status == ContainerEngineStatus.Stopped || Status == ContainerEngineStatus.Destroyed)
                return;

            Status = ContainerEngineStatus.Stopping;

            Exception hyperVException = null;
            try
            {
                _dnsUpdater.UnregisterWatcher();
                _hyperV.Stop();
            }
            catch (Exception ex)
            {
                Status = ContainerEngineStatus.FailedToStop;
                hyperVException = ex;
            }

            _proxy.Stop();
            _vpnKit.Stop();
            _dataKit.Stop();

            if (hyperVException != null)
            {
                throw hyperVException;
            }

            Status = ContainerEngineStatus.Stopped;
        }

        public void Destroy(bool keepVolume)
        {
            if (Status == ContainerEngineStatus.Destroyed)
                return;

            Status = ContainerEngineStatus.Destroying;
            try
            {
                _hyperV.Destroy(keepVolume);
                _firewall.RemoveRules();
                _hyperVGuids.Remove();
            }
            catch (Exception)
            {
                Status = ContainerEngineStatus.FailedToDestroy;
                throw;
            }
            Status = ContainerEngineStatus.Destroyed;
        }

        internal void CheckInstallation()
        {
            _featuresInstaller.CheckInstalledFeatures(_features);
            _hyperV.CheckHyperVState();
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
            }
        }
    }
}
