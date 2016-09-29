using System;
using System.ComponentModel;
using System.Runtime.ConstrainedExecution;
using System.Runtime.InteropServices;
using System.Security.Cryptography;
using System.Text;
using System.Windows;
using Docker.Core;

namespace Docker.WPF.Credentials
{
    public class CredentialAsker : ICredentialAsker
    {
        private static string CredentialPath => "Docker Host Filesystem Access";

        private readonly Logger _logger;
        private readonly ICmd _cmd;

        public CredentialAsker(ICmd cmd)
        {
            _logger = new Logger(GetType());
            _cmd = cmd;
        }

        public void Reset()
        {
            Reset(CredentialPath);
        }

        public Credential GetStoredCredential()
        {
            return RetrieveCredential(CredentialPath);
        }

        public bool AskCredential()
        {
            var error = false;

            try
            {
                var credential = Application.Current?.Dispatcher.InvokeAsync(() =>
                {
                    Credential creds;
                    do
                    {
                        var credentialWindow = new CredentialWindow();
                        if (error)
                        {
                            credentialWindow.Error.Content = "Invalid login/password";
                        }

                        creds = credentialWindow.Ask();
                        if (creds == null)
                        {
                            return null;
                        }

                        error = true;
                    } while (!ValidateCredential(creds));

                    return creds;
                }).Result;

                if (credential == null)
                {
                    return false;
                }

                _logger.Info($"Storing credentials: {credential.User}:***********");
                StoreCredential(CredentialPath, credential);
                return true;
            }
            catch (Exception ex)
            {
                _logger.Error($"Failed to retrieve the user's credentials: {ex.Message}");
                return false;
            }
        }

        // Ensures credentials come from a valid identity
        private bool ValidateCredential(Credential credential)
        {
            try
            {
                var credInfo = CredentialInfo.FromCredential(credential);
                return credInfo != null;
            }
            catch (Exception e)
            {
                _logger.Error($"Unable to validate credentials: {e.Message}");
                return false;
            }
        }

        private static byte[] EncryptPassword(string password)
        {
            return ProtectedData.Protect(Encoding.UTF8.GetBytes(password), null, DataProtectionScope.CurrentUser);
        }

        private static string DecryptPassword(byte[] encrypted)
        {
            return Encoding.UTF8.GetString(ProtectedData.Unprotect(encrypted, null, DataProtectionScope.CurrentUser));
        }

        internal static void StoreCredential(string target, Credential credential)
        {
            var encryptedPassword = EncryptPassword(credential.Password);

            var c = new NativeMethods.Credential
            {
                UserName = credential.User,
                TargetName = target,
                Persist = NativeMethods.CredPersist.Enterprise,
                CredentialBlob = Marshal.AllocHGlobal(encryptedPassword.Length)
            };

            try
            {
                Marshal.Copy(encryptedPassword, 0, c.CredentialBlob, encryptedPassword.Length);

                c.CredentialBlobSize = (uint)encryptedPassword.Length;
                c.Type = NativeMethods.CredTypes.CredTypeGeneric;

                if (!NativeMethods.CredWrite(ref c, 0))
                {
                    throw new Win32Exception(Marshal.GetLastWin32Error());
                }
            }
            finally
            {
                Marshal.FreeCoTaskMem(c.CredentialBlob);
            }
        }

        internal static Credential RetrieveCredential(string target)
        {
            IntPtr credential;
            var result = NativeMethods.CredRead(target, NativeMethods.CredTypes.CredTypeGeneric, 0, out credential);
            var error = Marshal.GetLastWin32Error();
            if (!result)
            {
                if (error == (int)NativeMethods.CredUiReturnCodes.ErrorNotFound)
                    return null;

                throw new Win32Exception(error);
            }

            try
            {
                var c = (NativeMethods.Credential)Marshal.PtrToStructure(credential, typeof(NativeMethods.Credential));
                var encryptedPassword = new byte[c.CredentialBlobSize];
                Marshal.Copy(c.CredentialBlob, encryptedPassword, 0, encryptedPassword.Length);
                return new Credential(c.UserName, DecryptPassword(encryptedPassword));
            }
            finally
            {
                NativeMethods.CredFree(credential);
            }
        }

        internal void Reset(string target)
        {
            _cmd.Run("cmdkey.exe", $"/delete \"{target}\"");
        }

        private static class NativeMethods
        {
            internal enum CredUiReturnCodes
            {
                ErrorNotFound = 1168
            }

            internal enum CredTypes
            {
                CredTypeGeneric = 1
            }

            internal enum CredPersist
            {
                Enterprise = 3
            }

            [DllImport("advapi32.dll", CharSet = CharSet.Unicode, EntryPoint = "CredReadW", SetLastError = true)]
            [return: MarshalAs(UnmanagedType.Bool)]
            internal static extern bool CredRead(string targetName, CredTypes type, int flags, out IntPtr credential);

            [DllImport("advapi32.dll"), ReliabilityContract(Consistency.WillNotCorruptState, Cer.Success)]
            internal static extern void CredFree(IntPtr buffer);

            [DllImport("advapi32.dll", CharSet = CharSet.Unicode, EntryPoint = "CredWriteW", SetLastError = true)]
            [return: MarshalAs(UnmanagedType.Bool)]
            internal static extern bool CredWrite(ref Credential credential, int flags);

#pragma warning disable 169
#pragma warning disable 414
            public struct Credential
            {
                public int Flags;
                public CredTypes Type;
                [MarshalAs(UnmanagedType.LPWStr)]
                public string TargetName;
                [MarshalAs(UnmanagedType.LPWStr)]
                public string Comment;
                public long LastWritten;
                public uint CredentialBlobSize;
                public IntPtr CredentialBlob;
                [MarshalAs(UnmanagedType.U4)]
                public CredPersist Persist;
                public int AttributeCount;
                public IntPtr Attributes;
                [MarshalAs(UnmanagedType.LPWStr)]
                public string TargetAlias;
                [MarshalAs(UnmanagedType.LPWStr)]
                public string UserName;
            }
#pragma warning restore 414
#pragma warning restore 169
        }
    }
}
