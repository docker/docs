using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class CredentialTest
    {
        [Test]
        public void Equal()
        {
            Check.That(new Credential("user1", "pwd1")).IsEqualTo(new Credential("user1", "pwd1"));
            Check.That(new Credential("user2", "pwd2")).IsEqualTo(new Credential("user2", "pwd2"));
            Check.That(new Credential(null, null)).IsEqualTo(new Credential(null, null));
        }

        [Test]
        public void NotEqual()
        {
            Check.That(new Credential("user", "pwd")).IsNotEqualTo(new Credential("user", "pwd2"));
            Check.That(new Credential("user", "pwd")).IsNotEqualTo(new Credential("user2", "pwd"));
            Check.That(new Credential("user", "pwd")).IsNotEqualTo(new Credential("user2", null));
            Check.That(new Credential("user", null)).IsNotEqualTo(new Credential("user2", "pwd"));
        }
    }
}