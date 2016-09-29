using System;
using System.Collections.Generic;
using System.Management.Automation;
using System.Threading;
using System.Threading.Tasks;
using Docker.Backend;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.Pipe;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class PipeServerTest
    {
        private static readonly NamedPipeServer Server = new NamedPipeServer("testpipe");
        private static readonly NamedPipeClient Client = new NamedPipeClient("testpipe");

        [OneTimeSetUp]
        public void Setup()
        {
            Task.Run(() => Server.Run());
        }

        [OneTimeTearDown]
        public void Teardown()
        {
            Server?.Stop();
        }

        [Test]
        public void Nominal()
        {
            var results = new List<object[]>();
            Server.Register("test2", parameters =>
            {
                results.Add(parameters);
                return "foo";
            });

            var output = Client.Send("test2", "bArbapapa", "second");

            Check.That(output).Equals("foo");
            Check.That(results[0]).ContainsExactly("bArbapapa", "second");
        }

        [Test]
        public void MultipleParameters()
        {
            var results = new List<object[]>();
            Server.Register("test3", parameters =>
            {
                results.Add(parameters);
                return "";
            });

            var output = Client.Send("test3", "foo", "bar", "qix", "bazs");

            Check.That(output).IsEqualTo("");
            Check.That(results[0]).ContainsExactly("foo", "bar", "qix", "bazs");
        }

        [Test]
        public void Error()
        {
            Server.Register("Create", request => { throw new RuntimeException("BUG"); });

            Check.ThatCode(() => Client.Send("Create")).Throws<BackendException>().WithMessage("BUG");
        }

        [Test]
        public void MultiLineDockerException()
        {
            Server.Register("Error", request => { throw new DockerException("BUG\nFOO\nBAR\n"); });

            Check.ThatCode(() => Client.Send("Error", "foo")).Throws<BackendException>().WithMessage("BUG\nFOO\nBAR\n");
        }

        [Test]
        public void HyperVException()
        {
            Server.Register("HyperVError", request => { throw new HyperVException("BUG\nFOO\nBAR\n"); });

            Check.ThatCode(() => Client.Send("HyperVError")).Throws<BackendException>().WithMessage("BUG\nFOO\nBAR\n");
        }

        [Test]
        public void NonDeserializableException()
        {
            Server.Register("CustomExtension", request => { throw new CustomExtension("BUG"); });

            Check.ThatCode(() => Client.Send("CustomExtension")).Throws<DockerException>().WithMessage("BUG");
        }

        [Test]
        public void MultiLineLinux()
        {
            var results = new List<object[]>();
            Server.Register("delete", parameters =>
            {
                results.Add(parameters);
                return "foo\nbar\nqix\n";
            });

            var output = Client.Send("delete", "bArbapapa");

            Check.That(output).IsEqualTo("foo\nbar\nqix\n");
            Check.That(results[0]).ContainsExactly("bArbapapa");
        }

        [Test]
        public void MultiLineWindows()
        {
            var results = new List<object[]>();
            Server.Register("test", parameters =>
            {
                results.Add(parameters);
                return "foo\n\rbar\n\rqix\n\r";
            });

            var output = Client.Send("test", "bArbapapa");

            Check.That(output).IsEqualTo("foo\n\rbar\n\rqix\n\r");
            Check.That(results[0]).ContainsExactly("bArbapapa");
        }

        [Test]
        public void UnknownAction()
        {
            Check.ThatCode(() => Client.Send("unknown", "bArbapapa")).Throws<DockerException>().WithMessage("Unknown action unknown");
        }

        [Test]
        public void SupportLongLastingConnections()
        {
            Server.Register("FIRST", parameters => { Thread.Sleep(100*1000); });
            Server.Register("SECOND", parameters => "OK");

            var task = Task.Run(() => Client.Send("FIRST"));
            Thread.Sleep(100);

            var anotherClient = new NamedPipeClient("testpipe");
            var output = anotherClient.Send("SECOND");

            Check.That(output).IsEqualTo("OK");
            Check.That(task.IsCompleted).IsFalse();
        }

        [Serializable]
        public class CustomExtension : Exception
        {
            public CustomExtension(string message) : base(message)
            {
            }
        }
    }
}