using System;
using Docker.Backend.Wmi;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class WmiTest
    {
        private readonly Wmi _wmi = new Wmi();

        [Test]
        public void Get()
        {
            var vm = _wmi.GetOrNull("Win32_OperatingSystem");
            var manufacturer = vm.Get("Manufacturer");

            Check.That(manufacturer).IsEqualTo("Microsoft Corporation");
        }

        [Test]
        public void Find()
        {
            var vm = _wmi.GetOrNull("Win32_OperatingSystem", _ => "Microsoft Corporation".Equals(_.Properties["Manufacturer"].Value));
            var manufacturer = vm.Get("Manufacturer");

            Check.That(manufacturer).IsEqualTo("Microsoft Corporation");
        }

        [Test]
        public void NotFound()
        {
            var vm = _wmi.GetOrNull("Win32_OperatingSystem", _ => "INVALID".Equals(_.Properties["Manufacturer"].Value));

            Check.That(vm).IsNull();
        }

        [Test]
        public void FailToFind()
        {
            Check.ThatCode(() => _wmi.GetOrFail("Win32_OperatingSystem", _ => "INVALID".Equals(_.Properties["Manufacturer"].Value)))
                .Throws<ArgumentException>()
                .WithMessage("Not found System.Management.ManagementScope/Win32_OperatingSystem");
        }
    }
}
