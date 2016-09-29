using System.IO;
using Docker.Core;
using Docker.Core.Tracking;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class TrackingTest
    {
        [Test]
        public void Create()
        {
            var filename = Path.GetTempFileName();

            var trackId = Tracking.ReadOrCreate(filename);

            Check.That(trackId).IsNotEmpty();
            Check.That(File.Exists(filename)).IsTrue();
        }

        [Test]
        public void Read()
        {
            var filename = Path.GetTempFileName();
            File.WriteAllText(filename, @"foobar");

            var trackId = Tracking.ReadOrCreate(filename);

            Check.That(trackId).Equals("foobar");
        }

        [Test]
        public void DefaultToNoTracking()
        {
            Check.That(new Tracking(Channel.Master).IsEnabled).IsFalse();
            Check.That(new Tracking(Channel.Test).IsEnabled).IsFalse();
            Check.That(new Tracking(Channel.Beta).IsEnabled).IsFalse();
            Check.That(new Tracking(Channel.Stable).IsEnabled).IsFalse();
        }

        [Test]
        public void OverrideTrackingSetInBeta()
        {
            var status = new Tracking(Channel.Beta);

            status.ChangeTo(false);
            Check.That(status.IsEnabled).IsTrue();

            status.ChangeTo(true);
            Check.That(status.IsEnabled).IsTrue();
        }

        [Test]
        public void HonorTrackingInStable()
        {
            var status = new Tracking(Channel.Stable);

            status.ChangeTo(false);
            Check.That(status.IsEnabled).IsFalse();

            status.ChangeTo(true);
            Check.That(status.IsEnabled).IsTrue();
        }

        [Test]
        public void HonorTrackingInTest()
        {
            var status = new Tracking(Channel.Test);

            status.ChangeTo(false);
            Check.That(status.IsEnabled).IsFalse();

            status.ChangeTo(true);
            Check.That(status.IsEnabled).IsTrue();
        }

        [Test]
        public void HonorTrackingInMaster()
        {
            var status = new Tracking(Channel.Master);

            status.ChangeTo(false);
            Check.That(status.IsEnabled).IsFalse();

            status.ChangeTo(true);
            Check.That(status.IsEnabled).IsTrue();
        }
    }
}
