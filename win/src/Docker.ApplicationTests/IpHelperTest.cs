using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    class IpHelperTest
    {
        private readonly IpHelper _ipHelper = new IpHelper();

        [Test]
        public void GetMaskSize()
        {
            Check.That(_ipHelper.GetMaskSize("255.255.255.255")).IsEqualTo(32);
            Check.That(_ipHelper.GetMaskSize("255.255.255.254")).IsEqualTo(31);
            Check.That(_ipHelper.GetMaskSize("255.255.255.253")).IsEqualTo(-1);
            Check.That(_ipHelper.GetMaskSize("255.255.255.252")).IsEqualTo(30);
            Check.That(_ipHelper.GetMaskSize("255.255.255.251")).IsEqualTo(-1);
            Check.That(_ipHelper.GetMaskSize("255.255.255.250")).IsEqualTo(-1);
            Check.That(_ipHelper.GetMaskSize("255.255.255.249")).IsEqualTo(-1);
            Check.That(_ipHelper.GetMaskSize("255.255.255.248")).IsEqualTo(29);
            Check.That(_ipHelper.GetMaskSize("255.255.255.247")).IsEqualTo(-1);
            Check.That(_ipHelper.GetMaskSize("255.255.255.240")).IsEqualTo(28);
            Check.That(_ipHelper.GetMaskSize("255.255.255.224")).IsEqualTo(27);
            Check.That(_ipHelper.GetMaskSize("255.255.255.192")).IsEqualTo(26);
            Check.That(_ipHelper.GetMaskSize("255.255.255.128")).IsEqualTo(25);
            Check.That(_ipHelper.GetMaskSize("255.255.255.0")).IsEqualTo(24);
            Check.That(_ipHelper.GetMaskSize("255.255.254.0")).IsEqualTo(23);
            Check.That(_ipHelper.GetMaskSize("255.255.252.0")).IsEqualTo(22);
            Check.That(_ipHelper.GetMaskSize("255.255.248.0")).IsEqualTo(21);
            Check.That(_ipHelper.GetMaskSize("255.255.240.0")).IsEqualTo(20);
            Check.That(_ipHelper.GetMaskSize("255.255.224.0")).IsEqualTo(19);
            Check.That(_ipHelper.GetMaskSize("255.255.192.0")).IsEqualTo(18);
            Check.That(_ipHelper.GetMaskSize("255.255.128.0")).IsEqualTo(17);
            Check.That(_ipHelper.GetMaskSize("255.255.0.0")).IsEqualTo(16);
            Check.That(_ipHelper.GetMaskSize("255.254.0.0")).IsEqualTo(15);
            Check.That(_ipHelper.GetMaskSize("255.252.0.0")).IsEqualTo(14);
            Check.That(_ipHelper.GetMaskSize("255.248.0.0")).IsEqualTo(13);
            Check.That(_ipHelper.GetMaskSize("255.240.0.0")).IsEqualTo(12);
            Check.That(_ipHelper.GetMaskSize("255.224.0.0")).IsEqualTo(11);
            Check.That(_ipHelper.GetMaskSize("255.192.0.0")).IsEqualTo(10);
            Check.That(_ipHelper.GetMaskSize("255.128.0.0")).IsEqualTo(9);
            Check.That(_ipHelper.GetMaskSize("255.0.0.0")).IsEqualTo(8);
            Check.That(_ipHelper.GetMaskSize("254.0.0.0")).IsEqualTo(7);
            Check.That(_ipHelper.GetMaskSize("252.0.0.0")).IsEqualTo(6);
            Check.That(_ipHelper.GetMaskSize("248.0.0.0")).IsEqualTo(5);
            Check.That(_ipHelper.GetMaskSize("240.0.0.0")).IsEqualTo(4);
            Check.That(_ipHelper.GetMaskSize("224.0.0.0")).IsEqualTo(3);
            Check.That(_ipHelper.GetMaskSize("192.0.0.0")).IsEqualTo(2);
            Check.That(_ipHelper.GetMaskSize("128.0.0.0")).IsEqualTo(1);
            Check.That(_ipHelper.GetMaskSize("0.0.0.0")).IsEqualTo(0);
        }

        [Test]
        public void GetMaskFromSize()
        {
            Check.That(_ipHelper.GetMaskFromSize(32)).IsEqualTo("255.255.255.255");
            Check.That(_ipHelper.GetMaskFromSize(31)).IsEqualTo("255.255.255.254");
            Check.That(_ipHelper.GetMaskFromSize(30)).IsEqualTo("255.255.255.252");
            Check.That(_ipHelper.GetMaskFromSize(29)).IsEqualTo("255.255.255.248");
            Check.That(_ipHelper.GetMaskFromSize(28)).IsEqualTo("255.255.255.240");
            Check.That(_ipHelper.GetMaskFromSize(27)).IsEqualTo("255.255.255.224");
            Check.That(_ipHelper.GetMaskFromSize(26)).IsEqualTo("255.255.255.192");
            Check.That(_ipHelper.GetMaskFromSize(25)).IsEqualTo("255.255.255.128");
            Check.That(_ipHelper.GetMaskFromSize(24)).IsEqualTo("255.255.255.0");
            Check.That(_ipHelper.GetMaskFromSize(23)).IsEqualTo("255.255.254.0");
            Check.That(_ipHelper.GetMaskFromSize(22)).IsEqualTo("255.255.252.0");
            Check.That(_ipHelper.GetMaskFromSize(21)).IsEqualTo("255.255.248.0");
            Check.That(_ipHelper.GetMaskFromSize(20)).IsEqualTo("255.255.240.0");
            Check.That(_ipHelper.GetMaskFromSize(19)).IsEqualTo("255.255.224.0");
            Check.That(_ipHelper.GetMaskFromSize(18)).IsEqualTo("255.255.192.0");
            Check.That(_ipHelper.GetMaskFromSize(17)).IsEqualTo("255.255.128.0");
            Check.That(_ipHelper.GetMaskFromSize(16)).IsEqualTo("255.255.0.0");
            Check.That(_ipHelper.GetMaskFromSize(15)).IsEqualTo("255.254.0.0");
            Check.That(_ipHelper.GetMaskFromSize(14)).IsEqualTo("255.252.0.0");
            Check.That(_ipHelper.GetMaskFromSize(13)).IsEqualTo("255.248.0.0");
            Check.That(_ipHelper.GetMaskFromSize(12)).IsEqualTo("255.240.0.0");
            Check.That(_ipHelper.GetMaskFromSize(11)).IsEqualTo("255.224.0.0");
            Check.That(_ipHelper.GetMaskFromSize(10)).IsEqualTo("255.192.0.0");
            Check.That(_ipHelper.GetMaskFromSize(9)).IsEqualTo("255.128.0.0");
            Check.That(_ipHelper.GetMaskFromSize(8)).IsEqualTo("255.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(7)).IsEqualTo("254.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(6)).IsEqualTo("252.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(5)).IsEqualTo("248.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(4)).IsEqualTo("240.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(3)).IsEqualTo("224.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(2)).IsEqualTo("192.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(1)).IsEqualTo("128.0.0.0");
            Check.That(_ipHelper.GetMaskFromSize(0)).IsEqualTo("0.0.0.0");
        }

        [Test]
        public void IsValidAdress()
        {
            Check.That(_ipHelper.IsValidAdress("192.168.0.1")).IsTrue();
            Check.That(_ipHelper.IsValidAdress("192.168.0.a")).IsFalse();
            Check.That(_ipHelper.IsValidAdress("192.168.0")).IsFalse();
            Check.That(_ipHelper.IsValidAdress("192.168")).IsFalse();
            Check.That(_ipHelper.IsValidAdress("192")).IsFalse();
        }

        [Test]
        public void VmIpAndSwitchIp()
        {
            Check.That(_ipHelper.SwitchIp("192.168.10.0", 24)).IsEqualTo("192.168.10.1");
            Check.That(_ipHelper.VmIp("192.168.10.0", 24)).IsEqualTo("192.168.10.2");

            Check.That(_ipHelper.SwitchIp("10.0.75.0", 24)).IsEqualTo("10.0.75.1");
            Check.That(_ipHelper.VmIp("10.0.75.0", 24)).IsEqualTo("10.0.75.2");

            Check.That(_ipHelper.SwitchIp("192.168.1.128", 25)).IsEqualTo("192.168.1.129");
            Check.That(_ipHelper.VmIp("192.168.1.128", 25)).IsEqualTo("192.168.1.130");
        }

        [Test]
        public void InvalidSubNet()
        {
            Check.ThatCode(() => _ipHelper.SwitchIp("192.168.10.5", 24)).ThrowsAny().WithMessage("Subnet address [192.168.10.5] is incompatible with mask [255.255.255.0]");
            Check.ThatCode(() => _ipHelper.VmIp("192.168.10.5", 24)).ThrowsAny().WithMessage("Subnet address [192.168.10.5] is incompatible with mask [255.255.255.0]");
        }

        [Test]
        public void Mask()
        {
            Check.That(_ipHelper.SubnetMask("10.0.75.0", 24)).IsEqualTo("255.255.255.0");
        }
    }
}
