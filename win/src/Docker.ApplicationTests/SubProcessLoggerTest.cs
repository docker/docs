using Docker.Backend.Processes;
using Docker.Core;
using Moq;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class SubProcessLoggerTest
    {
        private readonly Mock<ILogger> _logger = new Mock<ILogger>();

        private SubProcessLogger _subProcessLogger;

        [SetUp]
        public void SetUp()
        {
            _subProcessLogger = new SubProcessLogger(_logger.Object);
        }

        [Test]
        public void RemoveDate()
        {
            _subProcessLogger.Log("2016/05/02 11:39:27 docker proxy (on deprecated port): ready");

            _logger.Verify(_ => _.Info("docker proxy (on deprecated port): ready"));
        }

        [Test]
        public void LogWithoutDate()
        {
            _subProcessLogger.Log("docker proxy: ready");

            _logger.Verify(_ => _.Info("docker proxy: ready"));
        }

        [Test]
        public void Debug()
        {
            _subProcessLogger.Log("com.docker.slirp.exe: [DEBUG] hvsock connect got Success: retrying in 1s");

            _logger.Verify(_ => _.Debug("com.docker.slirp.exe: hvsock connect got Success: retrying in 1s"));
        }

        [Test]
        public void Warning()
        {
            _subProcessLogger.Log("com.docker.slirp.exe: [WARNING] hvsock connect got Success: retrying in 1s");

            _logger.Verify(_ => _.Warning("com.docker.slirp.exe: hvsock connect got Success: retrying in 1s"));
        }

        [Test]
        public void Info()
        {
            _subProcessLogger.Log("com.docker.slirp.exe: [INFO] hvsock connect got Success: retrying in 1s");

            _logger.Verify(_ => _.Info("com.docker.slirp.exe: hvsock connect got Success: retrying in 1s"));
        }

        [Test]
        public void Error()
        {
            _subProcessLogger.Log("com.docker.slirp.exe: [ERROR] hvsock connect got Success: retrying in 1s");

            _logger.Verify(_ => _.Error("com.docker.slirp.exe: hvsock connect got Success: retrying in 1s"));
        }
    }
}
