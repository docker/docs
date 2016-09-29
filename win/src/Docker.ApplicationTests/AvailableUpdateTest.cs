using Docker.Core.Update;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class AvailableUpdateTest
    {
        [Test]
        public void BuildNumberWithProperString()
        {
            var update = new AvailableUpdate("1234", null, null, null);

            Check.That(update.BuildNumber()).IsEqualTo(1234);
        }

        [Test]
        public void BuildNumberWithBigString()
        {
            var update = new AvailableUpdate("1.11.1.1234", null, null, null);

            Check.That(update.BuildNumber()).IsEqualTo(1234);
        }

        [Test]
        public void BuildWithNullVersion()
        {
            var update = new AvailableUpdate(null, null, null, null);

            Check.That(update.BuildNumber()).IsEqualTo(-1);
        }

        [Test]
        public void BuildWithBadVersion()
        {
            var update = new AvailableUpdate("foobar", null, null, null);

            Check.That(update.BuildNumber()).IsEqualTo(-1);
        }

        [Test]
        public void HumanVersion()
        {
            var update = new AvailableUpdate("1234", "1.12.0-RC1-beta16", null, null);

            Check.That(update.ToHumanStringWithBuildNumber()).IsEqualTo("1.12.0-RC1-beta16 (build: 1234)");
        }
    }
}