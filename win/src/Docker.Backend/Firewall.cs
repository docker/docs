using System.Linq;
using WindowsFirewallHelper;
using Docker.Core;
using System.Runtime.InteropServices;

namespace Docker.Backend
{
    public interface IFirewall
    {
        void OpenPorts();
        void RemoveRules();
    }

    public class Firewall : IFirewall
    {
        private const string DockerProxy = "DockerProxy";

        private readonly Logger _logger;

        public Firewall()
        {
            _logger = new Logger(GetType());
        }

        public void OpenPorts()
        {
            var proxyExe = Paths.DockerProxyExe;

            try
            {
                _logger.Info($"Opening ports for {proxyExe}...");

                RemoveAllRules();

                var fw = FirewallManager.Instance;
                var rule = fw.CreateApplicationRule(FirewallProfiles.Private | FirewallProfiles.Public, DockerProxy, FirewallAction.Allow, proxyExe, FirewallProtocol.Any);
                fw.Rules.Add(rule);

                rule.IsEnable = true;

                _logger.Info("Ports are opened");
            }
            catch (COMException)
            {
                _logger.Warning("Firewall service is not running.");
            }
        }

        public void RemoveRules()
        {
            try
            {
                _logger.Info("Closing ports...");

                RemoveAllRules();

                _logger.Info("Ports are closed");
            }
            catch (COMException)
            {
                _logger.Warning("Firewall service is not running.");
            }
        }

        private static void RemoveAllRules()
        {
            var fw = FirewallManager.Instance;
            var rules = fw.Rules.Where(rule => rule.Name == DockerProxy).ToArray();
            foreach (var rule in rules)
            {
                fw.Rules.Remove(rule);
            }
        }
    }
}