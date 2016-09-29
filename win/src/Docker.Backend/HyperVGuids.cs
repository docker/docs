using System.Collections.Generic;
using Docker.Core;
using Microsoft.Win32;

namespace Docker.Backend
{
    public interface IHyperVGuids
    {
        void Install();
        void Remove();
    }

    public class HyperVGuids : IHyperVGuids
    {
        private const string GuestCommunicationServices = @"SOFTWARE\Microsoft\Windows NT\CurrentVersion\Virtualization\GuestCommunicationServices\";
        private const string ElementName = "ElementName";

        private readonly Dictionary<string, string> _guids = new Dictionary<string, string>
        {
            {"C378280D-DA14-42C8-A24E-0DE92A1028E2", "Docker configuration database"},
            {"30D48B34-7D27-4B0B-AAAF-BBBED334DD59", "Docker VPN proxy"},
            {"0B95756A-9985-48AD-9470-78E060895BE7", "Docker port forwarding"},
            {"23A432C2-537A-4291-BCB5-D62504644739", "Docker API"},
            {"445BA2CB-E69B-4912-8B42-D7F494D007EA", "Docker diagnostics server"}
        };

        private readonly Logger _logger;

        public HyperVGuids()
        {
            _logger = new Logger(GetType());
        }

        public void Install()
        {
            _logger.Info("Installing GUIDs...");

            try
            {
                CreateKeys();
            }
            catch
            {
                RemoveKeys();
                CreateKeys();
            }

            _logger.Info("GUIDs installed");
        }

        public void Remove()
        {
            _logger.Info("Removing GUIDs...");

            RemoveKeys();

            _logger.Info("GUIDs removed");
        }

        private void CreateKeys()
        {
            foreach (var guid in _guids)
            {
                using (var key = Registry.LocalMachine.CreateSubKey(GuestCommunicationServices + guid.Key))
                {
                    key?.SetValue(ElementName, guid.Value);
                }
            }
        }

        private void RemoveKeys()
        {
            foreach (var guid in _guids)
            {
                Registry.LocalMachine.DeleteSubKey(GuestCommunicationServices + guid.Key, false);
            }
        }
    }
}
