using System.Threading;
using Docker.Core;
using Docker.Core.di;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture, Apartment(ApartmentState.STA)]
    public class DockerModuleTest
    {
        [Test]
        public void CreateApp()
        {
            var logger = new Logger(GetType());

            using (var singletons = new Singletons(new DockerModule(logger, new string[0])))
            {
                var app = singletons.Get<App>();
                app.Shutdown();

                Check.That(app).IsNotNull();
            }
        }

        [Test]
        public void CreateAppForE2ETests()
        {
            var logger = new Logger(GetType());

            using (var singletons = new Singletons(new DockerModule(logger, new[]
            {
                "-DisableCheckForUpdates", "-DisableWelcomeWhale", "-DisableToolboxMigration",
                "-Shared=C", "-Shared=D", "-Username=docker",
                "-Password=ftw!"
            })))
            {
                var app = singletons.Get<App>();
                app.Shutdown();

                Check.That(app).IsNotNull();
            }
        }
    }
}
