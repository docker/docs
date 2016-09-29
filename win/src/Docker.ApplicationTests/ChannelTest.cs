using NFluent;
using Docker.Core;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class ChannelTest
    {
        [Test]
        public void IsStaging()
        {
            Check.That(Channel.Master.IsStaging()).IsTrue();
            Check.That(Channel.Test.IsStaging()).IsTrue();
            Check.That(Channel.Beta.IsStaging()).IsFalse();
            Check.That(Channel.Stable.IsStaging()).IsFalse();
        }

        [Test]
        public void IsStable()
        {
            Check.That(Channel.Master.IsStable()).IsFalse();
            Check.That(Channel.Test.IsStable()).IsFalse();
            Check.That(Channel.Beta.IsStable()).IsFalse();
            Check.That(Channel.Stable.IsStable()).IsTrue();
        }

        [Test]
        public void Parse()
        {
            Check.That(Channel.Parse(null)).IsSameReferenceThan(Channel.Master);
            Check.That(Channel.Parse("")).IsSameReferenceThan(Channel.Master);
            Check.That(Channel.Parse("whatever")).IsSameReferenceThan(Channel.Master);
            Check.That(Channel.Parse("master")).IsSameReferenceThan(Channel.Master);
            Check.That(Channel.Parse("test")).IsSameReferenceThan(Channel.Test);
            Check.That(Channel.Parse("beta")).IsSameReferenceThan(Channel.Beta);
            Check.That(Channel.Parse("stable")).IsSameReferenceThan(Channel.Stable);
        }
    }
}