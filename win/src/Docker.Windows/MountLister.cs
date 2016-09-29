using System.Collections.Generic;
using Docker.WPF;

namespace Docker
{
    public interface ISharedDrives
    {
        Dictionary<string, bool> List();
        void UnShare(IEnumerable<string> sharesToDisable);
        void SetShared(string shareName, bool enabled);
    }

    public class SettingsBasedSharedDrives : ISharedDrives
    {
        private readonly ISettingsLoader _settingsLoader;

        public SettingsBasedSharedDrives(ISettingsLoader settingsLoader)
        {
            _settingsLoader = settingsLoader;
        }

        public Dictionary<string, bool> List()
        {
            return _settingsLoader.Load().SharedDrives;
        }

        public void UnShare(IEnumerable<string> sharesToDisable)
        {
            _settingsLoader.SaveChanges(_ =>
            {
                foreach (var shareName in sharesToDisable)
                {
                    _.SharedDrives[shareName] = false;
                }
            });
        }

        public void SetShared(string shareName, bool enabled)
        {
            _settingsLoader.SaveChanges(_ => _.SharedDrives[shareName] = enabled);
        }
    }
}