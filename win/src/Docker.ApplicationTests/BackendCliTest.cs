using System;
using Docker.Cli;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.di;
using Docker.Core.sharing;
using Docker.WPF;
using Docker.WPF.Crash;
using Moq;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class BackendCliTest
    {
        private readonly Mock<IActions> _mockActions = new Mock<IActions>(MockBehavior.Strict);
        private readonly Mock<ITaskQueue> _mockTaskQueue = new Mock<ITaskQueue>(MockBehavior.Strict);
        private readonly Mock<ICrashReport> _mockCrashReport = new Mock<ICrashReport>(MockBehavior.Strict);
        private readonly Mock<IConsoleWriter> _mockConsole = new Mock<IConsoleWriter>();
        private readonly Mock<IShareHelper> _mockShareHelper = new Mock<IShareHelper>();
        private readonly Mock<IBackend> _mockBackend = new Mock<IBackend>();

        private void RunBackendCli(params string[] args)
        {
            var module = Module.Override(new BackendCliModule(args)).With(binder =>
            {
                binder.Bind(_mockActions.Object);
                binder.Bind(_mockTaskQueue.Object);
                binder.Bind(_mockCrashReport.Object);
                binder.Bind(_mockConsole.Object);
                binder.Bind(_mockShareHelper.Object);
                binder.Bind(_mockBackend.Object);
            });

            var singletons = new Singletons(module);
            var backendCli = singletons.Get<BackendCli>();

            backendCli.Run(args);
        }

        [TearDown]
        public void VerifyAllMocks()
        {
            Mock.VerifyAll(_mockActions, _mockTaskQueue, _mockCrashReport, _mockConsole, _mockShareHelper, _mockBackend);
        }

        [Test]
        public void SwitchDaemon()
        {
            _mockActions.Setup(_ => _.SwitchDaemon());

            RunBackendCli("-SwitchDaemon");
        }

        [Test]
        public void MigrateVolume()
        {
            _mockActions.Setup(_ => _.MigrateVolume("default", "/disk"));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-MigrateVolume=/disk");
        }

        [Test]
        public void ResetToDefault()
        {
            _mockActions.Setup(_ => _.ResetToDefault());
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-ResetToDefault");
        }

        [Test]
        public void Mount()
        {
            _mockShareHelper.Setup(_ => _.Mount("C", It.IsAny<Settings>())).Returns(true);
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-Mount=C");
        }

        [Test]
        public void Unmount()
        {
            _mockShareHelper.Setup(_ => _.Unmount("D"));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-Unmount=D");
        }

        [Test]
        public void SharedDrives()
        {
            _mockBackend.Setup(_ => _.SharedDrives()).Returns(new[] { "C", "D" });
            _mockConsole.Setup(_ => _.WriteLine("C,D"));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SharedDrives");
        }

        [Test]
        public void DontStartIfSecretIsMissing()
        {
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli("-Start");
        }

        [Test]
        public void Start()
        {
            _mockActions.Setup(_ => _.Start());
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-Start");
        }

        [Test]
        public void Stop()
        {
            _mockActions.Setup(_ => _.StopVm());
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-Stop");
        }

        [Test]
        public void ResetCredential()
        {
            _mockShareHelper.Setup(_ => _.ResetCredential());
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-ResetCredential");
        }

        [Test]
        public void SendDiagnostic()
        {
            _mockCrashReport.Setup(_ => _.SendDiagnostic()).Returns("ID");
            _mockConsole.Setup(_ => _.WriteLine("ID"));

            RunBackendCli(BackendCli.Secret, "-SendDiagnostic");
        }

        [Test]
        public void SetMemory()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetMemory=4096");

            Check.That(settings.VmMemory).IsEqualTo(4096);
        }

        [Test]
        public void SetCpus()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetCpus=4");

            Check.That(settings.VmCpus).IsEqualTo(4);
        }

        [Test]
        public void SetAutomaticDns()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetDNS=automatic");

            Check.That(settings.UseDnsForwarder).IsTrue();
        }

        [Test]
        public void SetFixedDns()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetDNS=8.8.8.8");

            Check.That(settings.UseDnsForwarder).IsFalse();
            Check.That(settings.NameServer).IsEqualTo("8.8.8.8");
        }

        [Test]
        public void SetIp()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetIP=10.0.76.0/255.255.255.248");

            Check.That(settings.SubnetAddress).IsEqualTo("10.0.76.0");
            Check.That(settings.SubnetMaskSize).IsEqualTo(29);
        }

        [Test]
        public void SetDefaultIp()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetIP=10.0.75.0/255.255.255.0");

            Check.That(settings.SubnetAddress).IsEqualTo("10.0.75.0");
            Check.That(settings.SubnetMaskSize).IsEqualTo(24);
        }

        [Test]
        public void SetDaemonJson()
        {
            var settings = new Settings();

            _mockActions.Setup(_ => _.RestartVm(It.IsAny<Action<Settings>>())).Callback<Action<Settings>>(action => action(settings));
            _mockTaskQueue.Setup(_ => _.Shutdown());

            RunBackendCli(BackendCli.Secret, "-SetDaemonJson={\"debug\":false}");

            Check.That(settings.DaemonOptions).IsEqualTo(@"{""debug"":false}");
        }
    }
}