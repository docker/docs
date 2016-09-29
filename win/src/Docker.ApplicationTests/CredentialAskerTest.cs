using Docker.Core;
using Docker.WPF.Credentials;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class CredentialAskerTest
    {
        private const string TestPath = "Test Docker";

        [TearDown]
        public void CleanUpCredentialManager()
        {
            new CredentialAsker(new Cmd()).Reset(TestPath);
        }

        [Test]
        public void StoreCredential()
        {
            var credential = new Credential("LOGIN", "PASSWORD");

            CredentialAsker.StoreCredential(TestPath, credential);
            var retrieved = CredentialAsker.RetrieveCredential(TestPath);

            Check.That(retrieved.User).IsEqualTo("LOGIN");
            Check.That(retrieved.Password).IsEqualTo("PASSWORD");
        }

        [Test]
        public void DeleteCredential()
        {
            new CredentialAsker(new Cmd()).Reset(TestPath);

            var retrieved = CredentialAsker.RetrieveCredential(TestPath);

            Check.That(retrieved).IsNull();
        }
    }
}