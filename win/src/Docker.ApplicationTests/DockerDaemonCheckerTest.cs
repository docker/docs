using System;
using System.Collections.Generic;
using Docker.Backend;
using Docker.Core;
using Moq;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class DockerDaemonCheckerTest
    {
        private readonly Cmd.ProcessExecutionInfo _checkSucceded = new Cmd.ProcessExecutionInfo { ExitCode = 0 };
        private readonly Cmd.ProcessExecutionInfo _checkFailed = new Cmd.ProcessExecutionInfo { ExitCode = 1 };

        private Mock<ICmd> _cmd;
        private DockerDaemonChecker _dockerDaemonChecker;

        [SetUp]
        public void InitMocks()
        {
            _cmd = new Mock<ICmd>();
            _dockerDaemonChecker = new DockerDaemonChecker(_cmd.Object);
        }

        [Test]
        public void EngineIsUp()
        {
            _cmd.Setup(_ => _.Run(Paths.DockerExe, "ps", 0)).Returns(_checkSucceded);

            Check.ThatCode(() => _dockerDaemonChecker.Check()).DoesNotThrow();

            _cmd.Verify(_ => _.Run(Paths.DockerExe, "ps", 0), Times.Exactly(1));
        }

        [Test]
        public void EngineIsDown()
        {
            _cmd.Setup(_ => _.Run(Paths.DockerExe, "ps", 0)).Returns(_checkFailed);

            Check.ThatCode(() => _dockerDaemonChecker.Check()).Throws<Exception>();

            _cmd.Verify(_ => _.Run(Paths.DockerExe, "ps", 0), Times.Exactly(10));
        }

        [Test]
        public void Retry()
        {
            _cmd.Setup(_ => _.Run(Paths.DockerExe, "ps", 0)).Returns(new Queue<Cmd.ProcessExecutionInfo>(new[] { _checkFailed, _checkFailed, _checkFailed, _checkSucceded }).Dequeue);

            Check.ThatCode(() => _dockerDaemonChecker.Check()).DoesNotThrow();

            _cmd.Verify(_ => _.Run(Paths.DockerExe, "ps", 0), Times.Exactly(4));
        }
    }
}