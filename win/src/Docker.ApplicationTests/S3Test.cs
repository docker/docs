using System.Collections.Generic;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.Tracking;
using Docker.WPF.Crash;
using Moq;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class S3Test
    {
        private static readonly Mock<IBackend> Backend = new Mock<IBackend>();
        private static readonly Mock<ICmd> Cmd = new Mock<ICmd>();
        private static readonly Tracking Tracking = new Tracking(Channel.Master, "Id", true);

        [Test]
        public void Upload()
        {
            var s3 = new S3(new Mock<DebugInfo>(Backend.Object, Cmd.Object).Object, Backend.Object, Tracking);

            s3.Upload("timestamp", new Dictionary<string, object>
            {
                {"Log", "Unit test"}
            });

            Backend.Verify(_ => _.DownloadVmLogs());
        }
    }
}