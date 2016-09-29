using System.Security.Principal;
using Docker.Backend;
using Microsoft.Win32;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class HyperVGuidsTest
    {
        private readonly HyperVGuids _hyperVGuids = new HyperVGuids();

        [SetUp]
        public void SkipTestForNonAdministrators()
        {
            Assume.That(() => new WindowsPrincipal(WindowsIdentity.GetCurrent()).IsInRole(WindowsBuiltInRole.Administrator));
        }

        private static object ElementName(string path)
        {
            var key = @"SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\GuestCommunicationServices\" + path;
            using (var registryKey = Registry.LocalMachine.OpenSubKey(key, false))
            {
                return registryKey?.GetValue("ElementName");
            }
        }

        [Test]
        public void Install()
        {
            _hyperVGuids.Remove();
            _hyperVGuids.Install();

            Check.That(ElementName("C378280D-DA14-42C8-A24E-0DE92A1028E2")).IsEqualTo("Docker configuration database");
            Check.That(ElementName("30D48B34-7D27-4B0B-AAAF-BBBED334DD59")).IsEqualTo("Docker VPN proxy");
            Check.That(ElementName("0B95756A-9985-48AD-9470-78E060895BE7")).IsEqualTo("Docker port forwarding");
            Check.That(ElementName("23A432C2-537A-4291-BCB5-D62504644739")).IsEqualTo("Docker API");
            Check.That(ElementName("445BA2CB-E69B-4912-8B42-D7F494D007EA")).IsEqualTo("Docker diagnostics server");
        }

        [Test]
        public void InstallTwiceWithoutFailing()
        {
            _hyperVGuids.Install();
            _hyperVGuids.Install();
        }

        [Test]
        public void Remove()
        {
            _hyperVGuids.Install();
            _hyperVGuids.Remove();

            Check.That(ElementName("C378280D-DA14-42C8-A24E-0DE92A1028E2")).IsNull();
            Check.That(ElementName("30D48B34-7D27-4B0B-AAAF-BBBED334DD59")).IsNull();
            Check.That(ElementName("0B95756A-9985-48AD-9470-78E060895BE7")).IsNull();
            Check.That(ElementName("23A432C2-537A-4291-BCB5-D62504644739")).IsNull();
            Check.That(ElementName("445BA2CB-E69B-4912-8B42-D7F494D007EA")).IsNull();
        }

        [Test]
        public void RemoveTwiceWithoutFailing()
        {
            _hyperVGuids.Remove();
            _hyperVGuids.Remove();
        }
    }
}