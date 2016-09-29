using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class CommandLineOptionsTest
    {
        [Test]
        public void NoArgs()
        {
            var options = new CommandLineOptions(new string[0]);

            Check.That(options.DisableCheckForUpdates).IsFalse();
            Check.That(options.DisableWelcomeWhale).IsFalse();
            Check.That(options.DisableToolboxMigration).IsFalse();
            Check.That(options.Username).IsNull();
            Check.That(options.Password).IsNull();
        }

        [Test]
        public void DisableCheckForUpdates()
        {
            var options = new CommandLineOptions(new[] { "-DisableCheckForUpdates" });

            Check.That(options.DisableCheckForUpdates).IsTrue();
            Check.That(options.DisableWelcomeWhale).IsFalse();
            Check.That(options.DisableToolboxMigration).IsFalse();
            Check.That(options.Username).IsNull();
            Check.That(options.Password).IsNull();
        }

        [Test]
        public void DisableWelcomeWhale()
        {
            var options = new CommandLineOptions(new[] { "-DisableWelcomeWhale" });

            Check.That(options.DisableWelcomeWhale).IsTrue();
            Check.That(options.DisableCheckForUpdates).IsFalse();
            Check.That(options.DisableToolboxMigration).IsFalse();
            Check.That(options.Username).IsNull();
            Check.That(options.Password).IsNull();
        }

        [Test]
        public void DisableToolboxMigration()
        {
            var options = new CommandLineOptions(new[] { "-DisableToolboxMigration" });

            Check.That(options.DisableToolboxMigration).IsTrue();
            Check.That(options.DisableWelcomeWhale).IsFalse();
            Check.That(options.DisableCheckForUpdates).IsFalse();
            Check.That(options.Username).IsNull();
            Check.That(options.Password).IsNull();
        }

        [Test]
        public void Credential()
        {
            var options = new CommandLineOptions(new[] { "-Username=docker", "-Password=ftw!" });

            Check.That(options.Username).IsEqualTo("docker");
            Check.That(options.Password).IsEqualTo("ftw!");
            Check.That(options.DisableWelcomeWhale).IsFalse();
            Check.That(options.DisableCheckForUpdates).IsFalse();
            Check.That(options.DisableToolboxMigration).IsFalse();
        }

        [Test]
        public void EmptyPassword()
        {
            var options = new CommandLineOptions(new[] { "-Username=user", "-Password=" });

            Check.That(options.Username).IsEqualTo("user");
            Check.That(options.Password).IsEmpty();
        }
    }
}
