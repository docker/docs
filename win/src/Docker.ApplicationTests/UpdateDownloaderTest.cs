using System.IO;
using System.Reflection;
using System.Threading.Tasks;
using Docker.Core;
using Docker.Core.Update;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class UpdateDownloaderTest
    {
        private string _currentDirectory;

        [SetUp]
        public void Init()
        {
            _currentDirectory = Path.GetDirectoryName(Assembly.GetExecutingAssembly().Location);
        }

        [Test]
        public void IsTestingFileExist()
        {
            var filePath = Path.Combine(_currentDirectory, "Resources\\docker.xml");

            Check.That(File.Exists(filePath)).IsTrue();
        }

        [Test]
        public async Task ParseHockeyAppFeed()
        {
            var url = $"file://{_currentDirectory}/Resources/hockeyapp.xml";

            var update = await new UpdateFeedDownloader(new Channel("custom", url, "mixpaneltoken", "bugsnapApiKey")).DownloadLatestUpdateInfo();

            Check.That(update.Version).IsEqualTo("1.11.1.3668");
            Check.That(update.ShortVersion).IsEqualTo("1.11.1.3668");
            Check.That(update.DownloadUrl).StartsWith("https://rink.hockeyapp.net/api/2/ap");
            Check.That(update.Notes).StartsWith("<div style=\'font-size: 110%");
        }

        [Test]
        public async Task ParseDockerFeed()
        {
            var url = $"file://{_currentDirectory}/Resources/docker.xml";

            var update = await new UpdateFeedDownloader(new Channel("custom", url, "mixpaneltoken", "bugsnapApiKey")).DownloadLatestUpdateInfo();

            Check.That(update.Version).IsEqualTo("2789");
            Check.That(update.ShortVersion).IsEqualTo("1.11.1");
            Check.That(update.DownloadUrl).StartsWith("https://app-server-update.s3-");
            Check.That(update.Notes).StartsWith("<ul><li><p>New</p><ul>");
        }
    }
}
