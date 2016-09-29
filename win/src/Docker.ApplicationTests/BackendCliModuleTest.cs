using Docker.Cli;
using Docker.Core.di;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class BackendCliModuleTest
    {
        [Test]
        public void CreateBackendCli()
        {
            using (var singletons = new Singletons(new BackendCliModule(new string[0])))
            {
                var backendCli = singletons.Get<BackendCli>();

                Check.That(backendCli).IsNotNull();
            }
        }
    }
}