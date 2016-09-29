using System.IO;
using Docker.WPF;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class SettingsLoaderTest
    {
        [Test]
        public void LoadBeta17Settings()
        {
            using (var reader = new StringReader(@"{""VmCpus"":2,""VmMemory"":2048,""SubnetAddress"":""10.0.75.0"",""SubnetMaskSize"":24,""UseDnsForwarder"":true,""NameServer"":""8.8.8.8"",""LocalhostPortForwarding"":false,""IsTracking"":true,""SharedDrives"":{""C"":false},""DaemonOptions"":""{  \""registry-mirrors\"": [],  \""insecure-registries\"": [],  \""debug\"": false}"",""SysCtlConf"":"""",""ExecutableDate"":""06/30/2016 16:20:18"",""AutoUpdateEnabled"":false,""StartAtLogin"":false,""ProxyExclude"":"""",""ProxyHttp"":"""",""ProxyHttps"":"""",""UseHttpProxy"":false}"))
            {
                var settings = SettingsLoader.ParseJson(reader);

                Check.That(settings.VmCpus).IsEqualTo(2);
                Check.That(settings.VmMemory).IsEqualTo(2048);
                Check.That(settings.SubnetAddress).IsEqualTo("10.0.75.0");
                Check.That(settings.SubnetMaskSize).IsEqualTo(24);
                Check.That(settings.UseDnsForwarder).IsTrue();
                Check.That(settings.NameServer).IsEqualTo("8.8.8.8");
                Check.That(settings.IsTracking).IsTrue();
                Check.That(settings.DaemonOptions).IsEqualTo(@"{  ""registry-mirrors"": [],  ""insecure-registries"": [],  ""debug"": false}");
                Check.That(settings.SysCtlConf).IsEqualTo("");
                Check.That(settings.AutoUpdateEnabled).IsFalse();
                Check.That(settings.StartAtLogin).IsFalse();
                Check.That(settings.ProxyExclude).IsEqualTo("");
                Check.That(settings.ProxyHttp).IsEqualTo("");
                Check.That(settings.ProxyHttps).IsEqualTo("");
                Check.That(settings.UseHttpProxy).IsFalse();
            }
        }
    }
}
