using Docker.Backend;
using Docker.Backend.ContainerEngine;
using Docker.Backend.Features;
using Docker.Backend.Processes;
using Docker.Core;
using Docker.Core.Features;
using Moq;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class LinuxContainerEngineTest
    {
        private readonly Mock<IHyperV> _hyperV = new Mock<IHyperV>(MockBehavior.Strict);
        private readonly Mock<IFirewall> _firewall = new Mock<IFirewall>(MockBehavior.Strict);
        private readonly Mock<IVpnKit> _vpnKit = new Mock<IVpnKit>(MockBehavior.Strict);
        private readonly Mock<IDataKit> _dataKit = new Mock<IDataKit>(MockBehavior.Strict);
        private readonly Mock<IProxyProcess> _proxy = new Mock<IProxyProcess>(MockBehavior.Strict);
        private readonly Mock<IHyperVGuids> _hyperVGuids = new Mock<IHyperVGuids>(MockBehavior.Strict);
        private readonly Mock<IDatabase> _database = new Mock<IDatabase>(MockBehavior.Strict);
        private readonly Mock<IDnsUpdater> _dnsUpdater = new Mock<IDnsUpdater>(MockBehavior.Strict);
        private readonly Mock<IDockerDaemonChecker> _dockerDaemonChecker = new Mock<IDockerDaemonChecker>(MockBehavior.Strict);
        private readonly Mock<IInstaller> _featuresInstaller = new Mock<IInstaller>(MockBehavior.Strict);
        private readonly Mock<IContainerEngineHelper> _containerEngineHelper = new Mock<IContainerEngineHelper>(MockBehavior.Strict);

        private Linux _engine;

        [SetUp]
        public void Setup()
        {
            _engine = new Linux(
                _hyperV.Object, _firewall.Object, _vpnKit.Object,
                _dataKit.Object, _proxy.Object, _hyperVGuids.Object,
                _database.Object, _dnsUpdater.Object, _dockerDaemonChecker.Object,
                _featuresInstaller.Object, _containerEngineHelper.Object);
        }

        [Test]
        public void Create()
        {
            var settings = new Settings();

            _hyperV.Setup(_ => _.Create(settings));

            _engine.Create(settings);
        }

        [Test]
        public void Start()
        {
            var settings = new Settings();

            var sequence = new MockSequence();
            _containerEngineHelper.InSequence(sequence).Setup(_ => _.ForceKillLingeringDaemon());
            _hyperVGuids.InSequence(sequence).Setup(_ => _.Install());
            _firewall.InSequence(sequence).Setup(_ => _.OpenPorts());
            _hyperV.InSequence(sequence).Setup(_ => _.Create(settings));
            _dataKit.InSequence(sequence).Setup(_ => _.Start(settings));
            _vpnKit.InSequence(sequence).Setup(_ => _.Start(settings));
            _database.InSequence(sequence).Setup(_ => _.Write(settings));
            _dnsUpdater.InSequence(sequence).SetupSet(_ => _.StaticDnsAndSearchDomains = "");
            _dnsUpdater.InSequence(sequence).Setup(_ => _.UpdateDnsAndSearchDomains());
            _hyperV.InSequence(sequence).Setup(_ => _.Start());
            _dnsUpdater.InSequence(sequence).Setup(_ => _.RegisterWatcher());
            _proxy.InSequence(sequence).Setup(_ => _.Start(settings));
            _dockerDaemonChecker.InSequence(sequence).Setup(_ => _.Check());

            _engine.Status = ContainerEngineStatus.Stopped;
            _engine.Start(settings);
        }

        [Test]
        public void Restart()
        {
            var settings = new Settings();

            var sequence = new MockSequence();
            _proxy.InSequence(sequence).Setup(_ => _.Start(settings));
            _dockerDaemonChecker.InSequence(sequence).Setup(_ => _.Check());

            _engine.Status = ContainerEngineStatus.Started;
            _engine.Start(settings);
        }

        [Test]
        public void Stop()
        {
            var sequence = new MockSequence();
            _dnsUpdater.InSequence(sequence).Setup(_ => _.UnregisterWatcher());
            _hyperV.InSequence(sequence).Setup(_ => _.Stop());
            _proxy.InSequence(sequence).Setup(_ => _.Stop());
            _vpnKit.InSequence(sequence).Setup(_ => _.Stop());
            _dataKit.InSequence(sequence).Setup(_ => _.Stop());

            _engine.Stop();
        }

        [Test]
        public void StopProcessesOnHyperVFailure()
        {
            var sequence = new MockSequence();
            _hyperV.InSequence(sequence).Setup(_ => _.Stop()).Throws(new HyperVException("BUG"));
            _proxy.InSequence(sequence).Setup(_ => _.Stop());
            _vpnKit.InSequence(sequence).Setup(_ => _.Stop());
            _dataKit.InSequence(sequence).Setup(_ => _.Stop());

            _engine.Stop();
        }

        [Test]
        public void Destroy()
        {
            var sequence = new MockSequence();
            _hyperV.InSequence(sequence).Setup(_ => _.Destroy(false));
            _firewall.InSequence(sequence).Setup(_ => _.RemoveRules());
            _hyperVGuids.InSequence(sequence).Setup(_ => _.Remove());

            _engine.Destroy(false);
        }

        [Test]
        public void DestroyKeepVolume()
        {
            var sequence = new MockSequence();
            _hyperV.InSequence(sequence).Setup(_ => _.Destroy(true));
            _firewall.InSequence(sequence).Setup(_ => _.RemoveRules());
            _hyperVGuids.InSequence(sequence).Setup(_ => _.Remove());

            _engine.Destroy(true);
        }

        [Test]
        public void CheckInstallation()
        {
            var sequence = new MockSequence();
            _featuresInstaller.InSequence(sequence).Setup(_ => _.CheckInstalledFeatures(new[] { Feature.HyperV }));
            _hyperV.InSequence(sequence).Setup(_ => _.CheckHyperVState());

            _engine.CheckInstallation();
        }
    }
}