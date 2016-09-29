using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class EnvTest
    {
        [Test]
        public void MaxMemoryMustBePositive()
        {
            Check.That(Env.MaxMemory).IsPositive(); 
        }
    }
}
