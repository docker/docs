using System;
using System.IO;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.Tracking;
using Docker.WPF;
using Docker.WPF.Crash;
using Moq;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class BugSnagTest
    {
        private static readonly Mock<Core.Version> Version = new Mock<Core.Version>();
        private static readonly Mock<IBackend> Backend = new Mock<IBackend>();
        private static readonly Mock<Git> Git = new Mock<Git>();
        private static readonly Mock<ISettingsLoader> SettingsLoader = new Mock<ISettingsLoader>();
        private static readonly Tracking Tracking = new Tracking(Channel.Master, "Id", true);

        private string _tempFileName;
        private BugSnag _bugSnag;

        [SetUp]
        public void Setup()
        {
            _tempFileName = Path.GetTempFileName();
            Logger.Initialize(_tempFileName); // this sucks

            _bugSnag = new BugSnag(new Mock<Logger>("").Object, Version.Object, new DebugInfo(Backend.Object, new Mock<ICmd>().Object), Channel.Master, Tracking, Git.Object, SettingsLoader.Object);
        }

        // This actually sends a real notification
        // This is to pinpoint when newtonsoft DLL is not wrapped anymore.
        [Test]
        public void TestNotify()
        {
            _bugSnag.Notify(new Exception("BugSnagUnitTest"));
        }

        [Test]
        public void TestLogFileIsRead()
        {
            File.WriteAllText(_tempFileName, @"foobar
");

            var metadataStore = _bugSnag.BuildMetadata(new Exception("error")).MetadataStore;

            Check.That(metadataStore["LOG"]["logfile"].ToString()).StartsWith("foobar");
        }

        [Test]
        public void SendChannel()
        {
            var metadataStore = _bugSnag.BuildMetadata(new Exception("error")).MetadataStore;

            Check.That(metadataStore["APP"]["channel"].ToString()).IsEqualTo("Master");
        }

        [Test]
        public void SendServiceStacktraceChannel()
        {
            try
            {
                throw new DockerException("error");
            }
            catch (Exception exception)
            {
                try
                {
                    throw new BackendException(exception);
                }
                catch (Exception backendException)
                {
                    var metadataStore = _bugSnag.BuildMetadata(backendException).MetadataStore;
                    Check.That(metadataStore["SERVICE"]["stacktrace"].ToString()).StartsWith("   at Docker.Tests.BugSnagTest.SendServiceStacktraceChannel() in");
                }
            }
        }
    }
}
