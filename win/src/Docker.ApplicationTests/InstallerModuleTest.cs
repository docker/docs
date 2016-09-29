using Docker.Core.di;
using Docker.Installer;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class InstallerModuleTest
    {
        [Test]
        public void CreateInstaller()
        {
            using (var singletons = new Singletons(new InstallerModule()))
            {
                var program = singletons.Get<Installer.Program>();

                Check.That(program).IsNotNull();
            }
        }
    }
}