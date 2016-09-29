using Docker.Backend;
using Docker.Backend.ContainerEngine;
using Docker.Backend.Features;
using Docker.Core;
using Docker.Core.Features;
using Moq;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class BackendTest
    {
        private readonly Mock<IVersion> _version = new Mock<IVersion>(MockBehavior.Strict);
        private readonly Mock<ISambaShare> _sambaShare = new Mock<ISambaShare>(MockBehavior.Strict);
        private readonly Mock<IHyperV> _hyperV = new Mock<IHyperV>(MockBehavior.Strict);
        private readonly Mock<IDockerMachineImport> _dockerMachineImport = new Mock<IDockerMachineImport>(MockBehavior.Strict);
        private readonly Mock<ICmd> _cmd = new Mock<ICmd>(MockBehavior.Strict);
        private readonly Mock<IFiles> _files = new Mock<IFiles>(MockBehavior.Strict);
        private readonly Mock<IContainerEngineHelper> _containerEngineHelper = new Mock<IContainerEngineHelper>(MockBehavior.Strict);
        private readonly Mock<ILinux> _linuxContainerEngine = new Mock<ILinux>(MockBehavior.Strict);
        private readonly Mock<IWindows> _windowsContainerEngine = new Mock<IWindows>(MockBehavior.Strict);
        private readonly Mock<IInstaller> _featuresInstaller = new Mock<IInstaller>(MockBehavior.Strict);

        private Backend.Backend _backend;

        [SetUp]
        public void Setup()
        {
            _backend = new Backend.Backend(
            _version.Object, _sambaShare.Object, _dockerMachineImport.Object,
            _cmd.Object, _files.Object, _containerEngineHelper.Object,
            _windowsContainerEngine.Object, _linuxContainerEngine.Object, _hyperV.Object,
            _featuresInstaller.Object);
        }

        [Test]
        public void Version()
        {
            _version.Setup(_ => _.ToHumanString()).Returns("VERSION");

            var version = _backend.Version();

            Check.That(version).IsEqualTo("VERSION");
        }

        [Test]
        public void Start()
        {
            var settings = new Settings();

            _containerEngineHelper.SetupGet(_ => _.UseLinuxContainerEngine).Returns(true);
            _linuxContainerEngine.Setup(_ => _.Start(settings));

            _backend.Start(settings);
        }

        [Test]
        public void Stop()
        {
            _linuxContainerEngine.Setup(_ => _.Stop());
            _windowsContainerEngine.Setup(_ => _.Stop());

            _backend.Stop();
        }

        [Test]
        public void Destroy()
        {
            _linuxContainerEngine.Setup(_ => _.Destroy(false));
            _windowsContainerEngine.Setup(_ => _.Destroy(false));

            _backend.Destroy(false);
        }

        [Test]
        public void DestroyKeepVolume()
        {
            _linuxContainerEngine.Setup(_ => _.Destroy(true));
            _windowsContainerEngine.Setup(_ => _.Destroy(true));

            _backend.Destroy(true);
        }

        [Test]
        public void Unmount()
        {
            _sambaShare.Setup(_ => _.Unmount("C"));

            _backend.Unmount("C");
        }

        [Test]
        public void MountSuccess()
        {
            var credential = new Credential("user", "pwd");
            _containerEngineHelper.Setup(_ => _.UseLinuxContainerEngine).Returns(true);
            _sambaShare.Setup(_ => _.Mount("D", credential, It.IsAny<Settings>())).Returns(true);

            var shared = _backend.Mount("D", credential, new Settings());

            Check.That(shared).IsTrue();
        }

        [Test]
        public void MountFailure()
        {
            var credential = new Credential("user", "pwd");
            _containerEngineHelper.Setup(_ => _.UseLinuxContainerEngine).Returns(true);
            _sambaShare.Setup(_ => _.Mount("F", credential, It.IsAny<Settings>())).Returns(false);

            var shared = _backend.Mount("F", credential, new Settings());

            Check.That(shared).IsFalse();
        }

        [Test]
        public void RemoveShare()
        {
            _sambaShare.Setup(_ => _.DeleteShare("C"));

            _backend.RemoveShare("C");
        }

        [Test]
        public void MigrateVolume()
        {
            var sequence = new MockSequence();

            _dockerMachineImport.InSequence(sequence).Setup(_ => _.MigrateVolume("C:\\path\\default.vmdk", Paths.TmpVolumeMigrationPath));
            _linuxContainerEngine.InSequence(sequence).Setup(_ => _.Stop());
            _linuxContainerEngine.InSequence(sequence).Setup(_ => _.Destroy(false));
            _files.InSequence(sequence).Setup(_ => _.ForceMove(Paths.TmpVolumeMigrationPath, Paths.MobyDiskPath));

            _backend.MigrateVolume("C:\\path\\default.vmdk");
        }

        [Test]
        public void GetDebugInfo()
        {
            var processExecutionInfo = new Cmd.ProcessExecutionInfo { CombinedOutput = "INFO" };
            _cmd.Setup(_ => _.RunAsAdministrator("powershell.exe", It.IsAny<string>(), 0)).Returns(processExecutionInfo);

            var info = _backend.GetDebugInfo();

            Check.That(info).IsEqualTo("INFO");
        }

        [Test]
        public void DownloadVmLogs()
        {
            _hyperV.Setup(_ => _.DownloadLogs()).Returns("/path");

            var path = _backend.DownloadVmLogs();

            Check.That(path).IsEqualTo("/path");
        }

        [Test]
        public void InstallFeatures()
        {
            var featuresToInstall = new[] {Feature.HyperV, Feature.Containers};
            var expectedFailedFeatures = new[] {Feature.HyperV};

            _featuresInstaller.Setup(_ => _.Install(featuresToInstall)).Returns(expectedFailedFeatures);

            var failedFeatures = _backend.InstallFeatures(featuresToInstall);

            Check.That(failedFeatures).ContainsExactly(expectedFailedFeatures);
        }
    }
}