using Docker.Core;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class SystemChangesNotifierTest
    {
        [Test]
        public void SystemChangesNotifierInitTest()
        {
            var notifier = new SystemChangesNotifier();

            notifier.OnDiskMount += Notifier_DiskChanged;
            notifier.OnDiskUnmount += Notifier_DiskChanged;
        }

        private static void Notifier_DiskChanged()
        {
        }
    }
}
