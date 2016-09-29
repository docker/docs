using System;
using System.Runtime.InteropServices;
using System.Security.Authentication;
using System.Security.Principal;
using System.Text;

// ReSharper disable UnusedMember.Local
// ReSharper disable InconsistentNaming
namespace Docker.Core
{
    public class CredentialInfo
    {
        private static class NativeMethods
        {
            [DllImport("advapi32.dll", SetLastError = true, BestFitMapping = false, ThrowOnUnmappableChar = true)]
            [return: MarshalAs(UnmanagedType.Bool)]
            internal static extern bool LogonUser(
                [MarshalAs(UnmanagedType.LPStr)] string pszUserName,
                [MarshalAs(UnmanagedType.LPStr)] string pszDomain,
                [MarshalAs(UnmanagedType.LPStr)] string pszPassword,
                LogonType dwLogonType,
                LogonProvider dwLogonProvider,
                ref IntPtr phToken);

            public enum LogonType
            {
                LOGON32_LOGON_INTERACTIVE = 2,
                LOGON32_LOGON_NETWORK = 3,
                LOGON32_LOGON_BATCH = 4,
                LOGON32_LOGON_SERVICE = 5,
                LOGON32_LOGON_UNLOCK = 7,
                LOGON32_LOGON_NETWORK_CLEARTEXT = 8,
                LOGON32_LOGON_NEW_CREDENTIALS = 9,
            }

            public enum LogonProvider
            {
                LOGON32_PROVIDER_DEFAULT = 0,
                LOGON32_PROVIDER_WINNT35 = 1,
                LOGON32_PROVIDER_WINNT40 = 2,
                LOGON32_PROVIDER_WINNT50 = 3
            }

            public const int NO_ERROR = 0;
            public const int ERROR_INSUFFICIENT_BUFFER = 122;

            public enum SID_NAME_USE
            {
                SidTypeUser = 1,
                SidTypeGroup,
                SidTypeDomain,
                SidTypeAlias,
                SidTypeWellKnownGroup,
                SidTypeDeletedAccount,
                SidTypeInvalid,
                SidTypeUnknown,
                SidTypeComputer
            }

            [DllImport("advapi32.dll", CharSet = CharSet.Auto, SetLastError = true)]
            public static extern bool LookupAccountSid(
                string lpSystemName,
                [MarshalAs(UnmanagedType.LPArray)] byte[] Sid,
                StringBuilder lpName,
                ref uint cchName,
                StringBuilder ReferencedDomainName,
                ref uint cchReferencedDomainName,
                out SID_NAME_USE peUse);
        }

        public string UserName { get; private set; }
        public string DomainName { get; private set; }

        public CredentialInfo(string userName, string domainName)
        {
            UserName = userName;
            DomainName = domainName;
        }

        private static bool LogonUser(string user, string domain, string password, ref IntPtr rawToken)
        {
            return NativeMethods.LogonUser(user, domain, password, NativeMethods.LogonType.LOGON32_LOGON_INTERACTIVE, NativeMethods.LogonProvider.LOGON32_PROVIDER_DEFAULT, ref rawToken);
        }

        public static CredentialInfo FromCredential(Credential cred)
        {
            var rawToken = IntPtr.Zero;

            // Try user as is was entered
            if (LogonUser(cred.User, null, cred.Password, ref rawToken))
            {
                return RetrieveCredentialInfo(rawToken);
            }

            // Try split with \
            if (cred.User.Contains("\\"))
            {
                var parts = cred.User.Split(new[] { '\\' }, 2);
                if (!LogonUser(parts[1], parts[0], cred.Password, ref rawToken))
                {
                    throw new InvalidCredentialException("Invalid username or password");
                }

                return RetrieveCredentialInfo(rawToken);
            }

            // Try split with /
            if (cred.User.Contains("/"))
            {
                var parts = cred.User.Split(new[] { '/' }, 2);
                if (!LogonUser(parts[1], parts[0], cred.Password, ref rawToken))
                {
                    throw new InvalidCredentialException("Invalid username or password");
                }

                return RetrieveCredentialInfo(rawToken);
            }

            throw new InvalidCredentialException("Invalid username or password");
        }

        private static CredentialInfo RetrieveCredentialInfo(IntPtr rawToken)
        {
            using (var handle = new WindowsIdentitySafeHandle())
            {
                handle.Handle = rawToken;
                var identity = new WindowsIdentity(handle.Handle);
                if (identity.User == null)
                {
                    throw new InvalidCredentialException("Unable to get the user");
                }

                var nameBuffer = new StringBuilder();
                var cchName = (uint) nameBuffer.Capacity;
                var domainNameBuffer = new StringBuilder();
                var cchReferencedDomainName = (uint) domainNameBuffer.Capacity;
                NativeMethods.SID_NAME_USE sidUse;
                var sid = new byte[identity.User.BinaryLength];
                identity.User.GetBinaryForm(sid, 0);

                if (
                    !NativeMethods.LookupAccountSid(null, sid, nameBuffer, ref cchName, domainNameBuffer,
                        ref cchReferencedDomainName, out sidUse))
                {
                    var err = Marshal.GetLastWin32Error();
                    if (err == NativeMethods.ERROR_INSUFFICIENT_BUFFER)
                    {
                        nameBuffer.EnsureCapacity((int) cchName);
                        domainNameBuffer.EnsureCapacity((int) cchReferencedDomainName);
                        err = NativeMethods.NO_ERROR;
                        if (
                            !NativeMethods.LookupAccountSid(null, sid, nameBuffer, ref cchName, domainNameBuffer,
                                ref cchReferencedDomainName, out sidUse))
                            err = Marshal.GetLastWin32Error();
                    }
                    if (err != NativeMethods.NO_ERROR)
                    {
                        throw new InvalidCredentialException("Unable to extract info for provided credential");
                    }
                }

                return new CredentialInfo(nameBuffer.ToString(), domainNameBuffer.ToString());
            }
        }
    }
}
