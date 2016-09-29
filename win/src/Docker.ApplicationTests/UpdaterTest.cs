using System;
using System.Threading.Tasks;
using Docker.Core;
using Docker.Core.Update;
using Moq;
using NUnit.Framework;
using static Moq.MockBehavior;
using Version = Docker.Core.Version;

namespace Docker.Tests
{
    [TestFixture]
    public class UpdaterTest
    {
        private readonly Mock<IFeedDownloader> _feedDownloader = new Mock<IFeedDownloader>(Strict);
        private readonly Mock<IInstallUpdateWindow> _updateDispatcher = new Mock<IInstallUpdateWindow>(Strict);
        private readonly Mock<Action> _performUpdate = new Mock<Action>(Strict);
        private readonly Mock<Action> _upToDate = new Mock<Action>(Strict);

        [Test]
        public async Task NoUpdateWhenVersionIsUnknown()
        {
            SetRemoteVersionTo(null);
            var updater = SetLocalVersionTo(200);

            _upToDate.Setup(action => action());

            await updater.CheckForUpdates(_performUpdate.Object, _upToDate.Object);
        }

        [Test]
        public async Task NoUpdateWhenVersionIsNotANumver()
        {
            SetRemoteVersionTo("foobar");
            var updater = SetLocalVersionTo(200);

            _upToDate.Setup(action => action());

            await updater.CheckForUpdates(_performUpdate.Object, _upToDate.Object);
        }

        [Test]
        public async Task NoUpdateWhenVersionAreLower()
        {
            SetRemoteVersionTo("200");
            var updater = SetLocalVersionTo(201);

            _upToDate.Setup(action => action());

            await updater.CheckForUpdates(_performUpdate.Object, _upToDate.Object);
        }

        [Test]
        public async Task NoUpdateWhenVersionAreEqual()
        {
            SetRemoteVersionTo("200");
            var updater = SetLocalVersionTo(200);

            _upToDate.Setup(action => action());

            await updater.CheckForUpdates(_performUpdate.Object, _upToDate.Object);
        }

        [Test]
        public async Task UpdateWhenVersionAreHigher()
        {
            var update = SetRemoteVersionTo("200");
            var updater = SetLocalVersionTo(199);

            _updateDispatcher.Setup(action => action.Open(update, _performUpdate.Object));

            await updater.CheckForUpdates(_performUpdate.Object, _upToDate.Object);
        }

        private AvailableUpdate SetRemoteVersionTo(string version)
        {
            var update = new AvailableUpdate(version, $"1.1.1.{version}", null, null);
            _feedDownloader.Setup(dl => dl.DownloadLatestUpdateInfo()).Returns(Task.FromResult(update));
            return update;
        }

        private IUpdater SetLocalVersionTo(int version)
        {
            return new Updater(new Version(1, 1, 1, version, ""), Channel.Master, _feedDownloader.Object, _updateDispatcher.Object);
        }
    }
}