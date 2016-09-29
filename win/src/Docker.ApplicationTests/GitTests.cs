using System.IO;
using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class GitTests
    {
        private const string FullSha1 = "058c69f72a3f450a032e6979f1d325422f82e461";
        private const string ShortSha1 = "058c69f";

        [Test]
        public void ReturnsSha1IfExist()
        {
            var filename = Path.GetTempFileName();

            File.WriteAllText(filename, FullSha1);

            Check.That(new Git(filename).Sha1()).Equals(FullSha1);
        }

        [Test]
        public void ReturnsShortSha1IfExist()
        {
            var filename = Path.GetTempFileName();

            File.WriteAllText(filename, FullSha1);

            Check.That(new Git(filename).ShortSha1()).Equals(ShortSha1);
        }

        [Test]
        public void ReturnsEmptyStringIfCantLoad()
        {
            var filename = "\\nonexisting.txt";

            Check.That(new Git(filename).Sha1()).IsEmpty();
            Check.That(new Git(filename).ShortSha1()).IsEmpty();
        }
    }
}