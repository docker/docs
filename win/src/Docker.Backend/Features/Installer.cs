using Docker.Core;
using Docker.Core.Features;
using System.Linq;
using Docker.Core.backend;

namespace Docker.Backend.Features
{
    public interface IInstaller
    {
        void CheckInstalledFeatures(Feature[] features);
        Feature[] Install(Feature[] features);
    }

    public class Installer : IInstaller
    {
        private readonly ICmd _cmd;

        public Installer(ICmd cmd)
        {
            _cmd = cmd;
        }

        public void CheckInstalledFeatures(Feature[] features)
        {
            var featuresToEnable = features.Where(feature => !IsFeatureEnabled(feature.Name)).ToArray();

            if (featuresToEnable.Length > 0)
                throw new BackendInstallFeatureException(featuresToEnable);
        }

        private bool IsFeatureEnabled(string name)
        {
            try
            {
                var wmi = new Wmi.Wmi(@"root\CIMV2");
                var containers = wmi.GetOrFail("Win32_OptionalFeature", _ => name.Equals(_.GetPropertyValue("Name")));
                return containers.Get("InstallState") == "1";
            }
            catch
            {
                return false;
            }
        }

        public Feature[] Install(Feature[] features)
        {
            return features.Where(feature => !EnableFeature(feature.Name)).ToArray();
        }

        public bool EnableFeature(string name)
        {
            if (IsFeatureEnabled(name))
                return true;

            _cmd.Run("dism.exe", $"/Online /Enable-Feature:{name} /All /NoRestart /Quiet");

            var retryCount = 10;
            while (retryCount > 0)
            {
                if (IsFeatureEnabled(name))
                    return true;
                System.Threading.Thread.Sleep(1000);
                retryCount--;
            }

            return false;
        }
    }
}
