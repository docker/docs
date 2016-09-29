using Docker.Backend;
using Docker.Core;
using Moq;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class HyperVTest
    {
        private readonly Mock<IPowerShell> _powerShell = new Mock<IPowerShell>();
        private readonly Mock<ICmd> _cmd = new Mock<ICmd>();

        private HyperV _hyperV;

        [SetUp]
        public void Setup()
        {
            _hyperV = new HyperV(_powerShell.Object, _cmd.Object);
        }

        [Test]
        public void GetId()
        {
            _powerShell.Setup(_ => _.Output("(Get-VM MobyLinuxVM).Id.Guid")).Returns("227ae5a6-b9f8-446b-9893-3dfd41443a1b");

            var id = _hyperV.GetId();

            Check.That(id).IsEqualTo("227ae5a6-b9f8-446b-9893-3dfd41443a1b");
        }
    }
}
