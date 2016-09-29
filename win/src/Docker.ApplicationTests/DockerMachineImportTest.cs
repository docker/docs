using Docker.Backend;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class DockerMachineImportTest
    {
        [Test]
        public void SupportUsernameWithSpace()
        {
            const string inputVolumePath = @"C:\Users\USERNAME WITH SPACE\.docker\machine\machines\default\disk.vmdk";
            const string outputVolumePath = @"C:\Users\Public\Documents\Hyper-V\Virtual hard disks\MigratedVolume.vhdx";

            var arguments = DockerMachineImport.ImportArguments(inputVolumePath, outputVolumePath);

            Check.That(arguments).IsEqualTo(@"convert ""C:\Users\USERNAME WITH SPACE\.docker\machine\machines\default\disk.vmdk"" -O vhdx -o subformat=dynamic -p ""C:\Users\Public\Documents\Hyper-V\Virtual hard disks\MigratedVolume.vhdx""");
        }
    }
}