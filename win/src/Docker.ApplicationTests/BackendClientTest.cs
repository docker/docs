using System.Management.Automation;
using Docker.Backend;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.Backend;
using Moq;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class BackendClientTest
    {
        private static readonly Mock<IBackend> Backend = new Mock<IBackend>();
        private static readonly BackendServer Server = new BackendServer(Backend.Object, null);
        private static readonly BackendClient Client = new BackendClient("dockerBackend");
        private static readonly Settings Settings = new Settings();

        [OneTimeSetUp]
        public void Setup()
        {
            Server.Run();
        }

        [SetUp]
        public void ResetMocks()
        {
            Backend.Reset();
        }

        [OneTimeTearDown]
        public void Teardown()
        {
            Server?.Stop();
        }

        [Test]
        public void Error()
        {
            Backend.Setup(_ => _.Start(It.IsAny<Settings>())).Throws(new RuntimeException("BUG"));

            Check.ThatCode(() => Client.Start(Settings)).Throws<BackendException>().WithMessage("BUG");
        }

        [Test]
        public void WorksAfterAnError()
        {
            Backend.Setup(_ => _.Stop()).Throws(new RuntimeException("BUG"));

            Check.ThatCode(() => Client.Stop()).Throws<BackendException>().WithMessage("BUG");

            Client.Start(Settings);

            Backend.Verify(_ => _.Start(It.IsAny<Settings>()));
        }

        [Test]
        public void Version()
        {
            Backend.Setup(_ => _.Version()).Returns("Version 1.0");

            var version = Client.Version();

            Check.That(version).IsEqualTo("Version 1.0");
        }

        [Test]
        public void Start()
        {
            Client.Start(Settings);

            Backend.Verify(_ => _.Start(It.IsAny<Settings>()));
        }

        [Test]
        public void Stop()
        {
            Client.Stop();

            Backend.Verify(_ => _.Stop());
        }

        [Test]
        public void DestroyKeepVolume()
        {
            Client.Destroy(true);

            Backend.Verify(_ => _.Destroy(true));
        }

        [Test]
        public void Destroy()
        {
            Client.Destroy(false);

            Backend.Verify(_ => _.Destroy(false));
        }

        [Test]
        public void SharedDrives()
        {
            Backend.Setup(_ => _.SharedDrives()).Returns(new[] {"C", "D"});

            var drives = Client.SharedDrives();

            Check.That(drives).ContainsExactly("C", "D");
        }

        [Test]
        public void Unmount()
        {
            Client.Unmount("D");

            Backend.Verify(_ => _.Unmount("D"));
        }

        [Test]
        public void MountSuccess()
        {
            var credential = new Credential("user", "pwd");
            Backend.Setup(_ => _.Mount("D", credential, It.IsAny<Settings>())).Returns(true);

            var shared = Client.Mount("D", credential, Settings);

            Check.That(shared).IsTrue();
        }

        [Test]
        public void MountFailure()
        {
            var credential = new Credential("user", "pwd");
            Backend.Setup(_ => _.Mount("F", credential, It.IsAny<Settings>())).Returns(false);

            var shared = Client.Mount("F", credential, Settings);

            Check.That(shared).IsFalse();
        }

        [Test]
        public void RemoveShare()
        {
            Backend.Setup(_ => _.RemoveShare("C"));

            Client.RemoveShare("C");
        }

        [Test]
        public void MigrateVolume()
        {
            Client.MigrateVolume("C:\\Users\\user\\.docker\\machines\\machine\\default\\disk.vmdk");

            Backend.Verify(_ => _.MigrateVolume("C:\\Users\\user\\.docker\\machines\\machine\\default\\disk.vmdk"));
        }

        [Test]
        public void DownloadVmLogs()
        {
            Client.DownloadVmLogs();

            Backend.Verify(_ => _.DownloadVmLogs());
        }
    }
}