using Docker.Core;
using Docker.Core.Update;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class VersionTest
    {
        [Test]
        public void GetAppVersionName()
        {
            Check.That(Version.ParseAppVersionName("1.11.1-beta10-master")).IsEqualTo("beta10");
            Check.That(Version.ParseAppVersionName("1.11.1-beta10-test")).IsEqualTo("beta10");
            Check.That(Version.ParseAppVersionName("1.11.1-beta10")).IsEqualTo("beta10");
            
            Check.That(Version.ParseAppVersionName("1.12.0-rc1-beta16-master")).IsEqualTo("beta16");
            Check.That(Version.ParseAppVersionName("1.12.0-rc1-beta16-test")).IsEqualTo("beta16");
            Check.That(Version.ParseAppVersionName("1.12.0-rc1-beta16")).IsEqualTo("beta16");
        }

        [Test]
        public void GetHumanVersionWithString()
        {
            var version = new Version(new System.Version(1, 2, 3, 4), "1.2.3-beta16-master");

            Check.That(version.ToHumanStringWithBuildNumber()).Equals("1.2.3-beta16-master (build: 4)");
        }

        [Test]
        public void GetHumanVersionWithStringOnLocal()
        {
            var version = new Version(new System.Version(1, 2, 3, 0), "1.2.3-beta16-master");

            Check.That(version.ToHumanStringWithBuildNumber()).Equals("1.2.3-beta16-master (build: local)");
        }

        [Test]
        public void HumanVersionAreEqualBetweemUpdateAndVersion()
        {
            var update = new AvailableUpdate("4", "1.2.3-beta16-master", null, null);
            var version = new Version(new System.Version(1, 2, 3, 4), "1.2.3-beta16-master");

            Check.That(update.ToHumanStringWithBuildNumber()).IsEqualTo(version.ToHumanStringWithBuildNumber());
        }

        [Test]
        public void CrossPlatformVersion()
        {
            var version = new Version(new System.Version(1, 2, 3, 4), "1.2.3-beta16-master");
            Check.That(version.ToCrossPlatformName()).IsEqualTo("1.2.3-beta16");
        }
    }
}