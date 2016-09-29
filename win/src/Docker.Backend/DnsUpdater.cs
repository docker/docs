using Docker.Core;
using System;
using System.Collections.Generic;
using System.Management;

namespace Docker.Backend
{
    public interface IDnsUpdater
    {
        void RegisterWatcher();
        void UnregisterWatcher();
        void UpdateDnsAndSearchDomains();

        string StaticDnsAndSearchDomains { get; set; }
    }

    public class DnsUpdater : IDnsUpdater
    {
        private readonly Logger _logger;
        private readonly IPowerShell _powerShell;
        private readonly IDatabase _database;

        private ManagementEventWatcher _eventWatcher;
        private string _lastDnsSettings;

        public string StaticDnsAndSearchDomains { get; set; }

        public DnsUpdater(IPowerShell powerShell, IDatabase database)
        {
            _logger = new Logger(GetType());
            _powerShell = powerShell;
            _database = database;
        }

        private static ManagementEventWatcher CreateEventWatcher(string targetInstance, string eventClassName)
        {
            var scope = new ManagementScope("root\\CIMV2");
            return new ManagementEventWatcher(scope, new WqlEventQuery
            {
                EventClassName = eventClassName,
                WithinInterval = new TimeSpan(0, 0, 1),
                Condition = $"TargetInstance ISA '{targetInstance}'"
            });
        }

        public void RegisterWatcher()
        {
            try
            {
                if (_eventWatcher == null)
                {
                    _eventWatcher = CreateEventWatcher("Win32_NetworkAdapterConfiguration", "__InstanceModificationEvent");
                }
            }
            catch (Exception ex)
            {
                _logger.Error($"failed to init: {ex.Message}");
                return;
            }

            try
            {
                _eventWatcher.EventArrived += (sender, e) =>
                {
                    _logger.Info("Network configuration change detected");
                    UpdateDnsAndSearchDomains();
                };
                _eventWatcher.Start();
            }
            catch
            {
                UnregisterWatcher();
                throw;
            }
        }

        public void UnregisterWatcher()
        {
            _eventWatcher?.Stop();
            _eventWatcher = null;
            _lastDnsSettings = null;
        }

        public void UpdateDnsAndSearchDomains()
        {
            try
            {
                var dnsSettings = StaticDnsAndSearchDomains.Length != 0
                    ? StaticDnsAndSearchDomains
                    : GetDnsSettingsFromEveryUpInterfacesSortedByMetric();

                if (dnsSettings.Equals(_lastDnsSettings)) return;

                _database.WriteDnsSettings(dnsSettings);
                _lastDnsSettings = dnsSettings;
            }
            catch (Exception e)
            {
                _logger.Error($"Unable to update dns settings: {e.Message}");
            }
        }

        private string GetDnsSettingsFromEveryUpInterfacesSortedByMetric()
        {
            var mainRouteInterfaceIdx = uint.MaxValue;
            try
            {
                uint.TryParse(_powerShell.Output("$(Find-NetRoute -RemoteIPAddress 8.8.8.8).InterfaceIndex[0]"),
                    out mainRouteInterfaceIdx);
            }
            catch
            {
                // ignored
            }

            var dnsSettings = "";
            var query = new SelectQuery("Win32_NetworkAdapterConfiguration");
            using (var searcher = new ManagementObjectSearcher(query))
            {
                var sortedNetworkInterfaces = new SortedDictionary<uint, ManagementBaseObject>();

                foreach (var networkInterface in searcher.Get())
                {
                    try
                    {
                        var ipEnabled = networkInterface.Properties["IPEnabled"].Value;
                        if (ipEnabled == null) continue;
                        if (ipEnabled.ToString() == "False") continue;

                        uint interfaceIndex;
                        uint.TryParse(networkInterface.Properties["InterfaceIndex"].Value.ToString(), out interfaceIndex);
                        var priority = uint.MinValue;
                        if (interfaceIndex != mainRouteInterfaceIdx)
                        {
                            var ipConnectionMetricText = networkInterface.Properties["IPConnectionMetric"].Value.ToString();
                            uint ipConnectionMetric;
                            if (uint.TryParse(ipConnectionMetricText, out ipConnectionMetric))
                            {
                                priority = ipConnectionMetric;
                            }
                        }

                        var dnsServerSearchOrder = networkInterface.Properties["DNSServerSearchOrder"].Value;
                        if (dnsServerSearchOrder != null)
                        {
                            sortedNetworkInterfaces[priority] = networkInterface;
                        }
                    }
                    catch
                    {
                        // ignored
                    }
                }

                foreach (var networkInterface in sortedNetworkInterfaces.Values)
                {
                    var dnsList = networkInterface.Properties["DNSServerSearchOrder"].Value as string[];
                    if (dnsList != null)
                    {
                        foreach (var dns in dnsList)
                        {
                            dnsSettings += dnsSettings.Length > 0 ? $"\nnameserver {dns}" : $"nameserver {dns}";
                        }
                    }

                    var domainSuffixList = networkInterface.Properties["DNSDomainSuffixSearchOrder"].Value as string[];
                    if (domainSuffixList != null)
                    {
                        foreach (var domainSuffix in domainSuffixList)
                        {
                            dnsSettings += dnsSettings.Length > 0 ? $"\nsearch {domainSuffix}" : $"search {domainSuffix}";
                        }
                    }
                }
            }

            return dnsSettings;
        }
    }
}
