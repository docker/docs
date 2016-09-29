using Docker.WPF;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class KitematicHelperTest
    {
        [Test]
        public void Outdated()
        {
            Check.That(KitematicHelper.IsOutdated("0.10.0")).IsTrue();
            Check.That(KitematicHelper.IsOutdated("0.11.0")).IsTrue();
            Check.That(KitematicHelper.IsOutdated("0.11.99")).IsTrue();
        }

        [Test]
        public void Valid()
        {
            Check.That(KitematicHelper.IsOutdated("0.12.0")).IsFalse();
            Check.That(KitematicHelper.IsOutdated("0.13.1")).IsFalse();
            Check.That(KitematicHelper.IsOutdated("1.0.0")).IsFalse();
        }

        [Test]
        public void MalformedIsOutdated()
        {
            Check.That(KitematicHelper.IsOutdated("v0.13.0")).IsTrue();
            Check.That(KitematicHelper.IsOutdated("UNKNOWN")).IsTrue();
            Check.That(KitematicHelper.IsOutdated("a.b.c")).IsTrue();
        }
    }
}