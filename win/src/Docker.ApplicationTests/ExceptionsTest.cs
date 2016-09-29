using System.IO;
using System.Management.Automation;
using System.Runtime.Serialization.Formatters.Binary;
using Docker.Backend;
using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class ExceptionsTest
    {
        [Test]
        public void CanSerializeDockerException()
        {
            var error = new DockerException("message", new RuntimeException("cause"));

            var deserialized = (DockerException) Deserialize(Serialize(error));

            Check.That(deserialized.Message).IsEqualTo(error.Message);
            Check.That(deserialized.InnerException.Message).IsEqualTo(error.InnerException.Message);
        }

        [Test]
        public void CanSerializeHyperVException()
        {
            var error = new HyperVException("message", new RuntimeException("cause"));

            var deserialized = (HyperVException) Deserialize(Serialize(error));

            Check.That(deserialized.Message).IsEqualTo(error.Message);
            Check.That(deserialized.InnerException.Message).IsEqualTo(error.InnerException.Message);
        }

        private static byte[] Serialize(object value)
        {
            var binary = new BinaryFormatter();

            var stream = new MemoryStream();
            binary.Serialize(stream, value);

            return stream.ToArray();
        }

        private static object Deserialize(byte[] buffer)
        {
            var binary = new BinaryFormatter();

            var stream = new MemoryStream(buffer);

            return binary.Deserialize(stream);
        }
    }
}