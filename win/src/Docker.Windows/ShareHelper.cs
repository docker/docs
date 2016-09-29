using System;
using Docker.Core;
using Docker.Core.backend;
using Docker.Core.sharing;
using Docker.WPF.Credentials;

namespace Docker
{
    public class ShareHelper : IShareHelper
    {
        private readonly ILogger _logger;
        private readonly IBackend _backend;
        private readonly ICredentialAsker _credentialAsker;
        private readonly ISharedDrives _sharedDrives;

        public ShareHelper(IBackend backend, ICredentialAsker credentialAsker, ISharedDrives sharedDrives)
        {
            _logger = new Logger(GetType());
            _backend = backend;
            _credentialAsker = credentialAsker;
            _sharedDrives = sharedDrives;
        }

        public void ResetCredential()
        {
            _credentialAsker.Reset();

            foreach (var drive in _sharedDrives.List().Keys)
            {
                _backend.RemoveShare(drive);
                Unmount(drive);
            }
        }

        public void Unmount(string drive)
        {
            try
            {
                _backend.Unmount(drive);
                _sharedDrives.SetShared(drive, false);
            }
            catch (Exception e)
            {
                _logger.Error($"Failed to unmount {drive}: {e.Message}");
            }
        }

        public bool Mount(string drive, Settings settings)
        {
            try
            {
                var enabled = TryMount(drive, settings);
                _sharedDrives.SetShared(drive, enabled);
                return enabled;
            }
            catch (Exception e)
            {
                _logger.Error($"Failed to mount {drive}: {e.Message}");
                return false;
            }
        }

        public void UpdateMounts(Settings settings)
        {
            var credential = _credentialAsker.GetStoredCredential();

            var sharedDrives = _sharedDrives.List();
            foreach (var drive in sharedDrives.Keys)
            {
                if (sharedDrives[drive])
                {
                    _backend.Mount(drive, credential, settings);
                }
                else
                {
                    Unmount(drive);
                }
            }
        }

        private bool TryMount(string drive, Settings settings)
        {
            var mountWithStoreCredentials = _backend.Mount(drive, _credentialAsker.GetStoredCredential(), settings);
            if (mountWithStoreCredentials)
            {
                return true;
            }

            if (!_credentialAsker.AskCredential())
            {
                return false;
            }

            var mountWithUpdatedCredentials = _backend.Mount(drive, _credentialAsker.GetStoredCredential(), settings);
            return mountWithUpdatedCredentials;
        }
    }
}