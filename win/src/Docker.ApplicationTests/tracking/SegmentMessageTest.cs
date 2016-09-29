using Docker.Core;
using Docker.Core.Tracking;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests.tracking
{
    [TestFixture]
    public class SegmentMessageTest
    {
        [Test]
        public void MessageJsonWhenTrackingEnabled()
        {
            var message = new SegmentMessage(true, "Master", "USER", AnalyticEvent.AppLaunched.Name, BuildMockedVersion());

            Check.That(message.ToJson()).Equals($@"{{""userId"":""USER"",""event"":""appLaunched"",""properties"":{{""os"":""windows"",""app major version"":1,""app minor version"":11,""app patch version"":2,""app version name"":""1.11.2-beta10"",""channel"":""master"",""os major version"":""{Env.Os.Name.ToLower()}"",""os minor version"":""{Env.Os.ReleaseId}"",""os patch version"":""{Env.Os.BuildNumber}"",""os language"":""{Env.Os.Language}""}}}}");
        }

        [Test]
        public void MessageJsonWhenTrackingDisabled()
        {
            var message = new SegmentMessage(false, "Master", "USER", AnalyticEvent.AppLaunched.Name, BuildMockedVersion());

            Check.That(message.ToJson()).Equals(@"{""userId"":""USER"",""event"":""appLaunched"",""properties"":{""os"":""windows"",""app major version"":1,""app minor version"":11,""app patch version"":2,""app version name"":""1.11.2-beta10"",""channel"":""master""}}");
        }

        private static Version BuildMockedVersion()
        {
            return new Version(1, 11, 2, 0, "1.11.2-beta10-master");
        }
    }
}
