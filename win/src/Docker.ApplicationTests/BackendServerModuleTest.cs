using Docker.Backend;
using Docker.Core.di;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class BackendServerModuleTest
    {
        [Test]
        public void CreateBackend()
        {
            using (var singletons = new Singletons(new BackendServerModule()))
            {
                var backend = singletons.Get<Backend.Backend>();

                Check.That(backend).IsNotNull();
            }
        }
    }
}