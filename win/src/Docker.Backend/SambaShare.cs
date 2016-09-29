using System;
using System.Collections.Generic;
using System.Linq;
using Docker.Core;

namespace Docker.Backend
{
    public interface ISambaShare
    {
        string[] SharedDrives();
        void Unmount(string drive);
        bool Mount(string drive, Credential credential, Settings settings);
        void DeleteShare(string drive);
    }

    public class SambaShare : ISambaShare
    {
        private const string CifsOptionsFormat = "noperm,iocharset=utf8,nobrl,mfsymlinks,vers=3.02,domain={0}";

        private readonly Logger _logger;
        private readonly ICmd _cmd;
        private readonly IMobyCommand _mobyCommand;
        private readonly IpHelper _ipHelper;
        private readonly HashSet<string> _sharedDrives;

        public SambaShare(ICmd cmd, IMobyCommand mobyCommand, IpHelper ipHelper)
        {
            _logger = new Logger(GetType());
            _cmd = cmd;
            _mobyCommand = mobyCommand;
            _ipHelper = ipHelper;
            _sharedDrives = new HashSet<string>();
        }

        public string[] SharedDrives()
        {
            var drives = _sharedDrives.ToArray();
            Array.Sort(drives);
            return drives;
        }

        public void Unmount(string drive)
        {
            _logger.Info($"Unmount {drive}");

            UnmountInMoby(drive);
        }

        public bool Mount(string drive, Credential credential, Settings settings)
        {
            if (credential == null || credential.EmptyPassword)
            {
                return false;
            }

            _logger.Info($"Mount {drive}");

            if (IsShared(drive))
            {
                if (MountInMoby(drive, credential, settings))
                {
                    return true;
                }

                DeleteShare(drive);
            }

            CreateShare(credential, drive);

            if (IsShared(drive))
            {
                if (MountInMoby(drive, credential, settings))
                {
                    return true;
                }
            }

            if (!credential.User.Contains("\\"))
            {
                return false;
            }

            var splittedUserName = credential.User.Split('\\');
            var credentialWithoutDomain = new Credential(splittedUserName[splittedUserName.Length - 1], credential.Password);

            return Mount(drive, credentialWithoutDomain, settings);
        }

        public void DeleteShare(string drive)
        {
            _logger.Info($"Removing share {drive}");

            var exitCode = _cmd.Run("net", $"share {drive} /delete");
            if (0 != exitCode)
            {
                _logger.Error($"Removing share \"{drive}\" failed with code: {exitCode}");
            }
        }

        private bool IsShared(string drive)
        {
            var exitCode = _cmd.Run("net", $"share {drive}");

            _logger.Info($"\"{drive}\" is {(exitCode == 0 ? "shared" : "not shared")}");

            return exitCode == 0;
        }

        private void CreateShare(Credential credential, string drive)
        {
            if (credential == null) return;

            var directory = drive + @":\";

            _logger.Info($"Creating share \"{directory}\" as \"{drive}\" with Full Control to \"{credential.User}\"");

            var exitCode = _cmd.Run("net", $"share {drive}={directory} /grant:\"{credential.User}\",FULL /CACHE:None");
            if (0 != exitCode)
            {
                _logger.Error($"Failed to create share \"{directory}\" as \"{drive}\" with Full Control to \"{credential.User}\" with code: {exitCode}");
            }
        }

        public bool MountInMoby(string drive, Credential credential, Settings settings)
        {
            CredentialInfo credInfo;
            try
            {
                credInfo = CredentialInfo.FromCredential(credential);
            }
            catch (Exception e)
            {
                _logger.Error($"Unable to validate credential: {e.Message}");
                return false;
            }

            var userPwd = new Dictionary<string, string>
            {
                // mount.cifs recognizes / as the domain separator
                // https://www.samba.org/samba/docs/man/manpages-3/mount.cifs.8.html#variablelist
                //
                {"USER", credInfo.UserName},
                {"PASSWD", credential.Password}
            };

            var hostIp = _ipHelper.SwitchIp(settings.SubnetAddress, settings.SubnetMaskSize);
            var cifsOptions = string.Format(CifsOptionsFormat, credInfo.DomainName);
            var success = MountInMoby(drive.ToLowerInvariant(), drive, cifsOptions, hostIp, userPwd, true)
                          && MountInMoby(drive.ToUpperInvariant(), drive, cifsOptions, hostIp, userPwd);

            if (success)
            {
                _sharedDrives.Add(drive);
            }

            return success;
        }

        private bool MountInMoby(string mountPoint, string drive, string cifsOptions, string hostIp, IDictionary<string, string> envVariables = null, bool load = false)
        {
            var result = _mobyCommand.Run($"(umount /{mountPoint} 2>/dev/null || true) && rm -Rf /{mountPoint} && mkdir -p /{mountPoint} && mount -t cifs -o {cifsOptions} //{hostIp}/{drive} /{mountPoint}", envVariables, load);
            if (result.ExitCode != 0)
            {
                _logger.Error($"Unable to mount {drive} drive: {result.ErrorOutput}");
            }

            return result.ExitCode == 0;
        }

        public void UnmountInMoby(string drive)
        {
            _sharedDrives.Remove(drive);
            _mobyCommand.Run($"umount /{drive.ToLowerInvariant()}; umount /{drive.ToUpperInvariant()}", null, true);
        }
    }
}
