using System;
using System.Runtime.InteropServices;
using System.Security.Principal;
using Docker.Core.Pipe;
using Docker.Core;

namespace Docker.Backend
{
    public class WindowsIdentityRecorder : IClientIdentityCallback
    {
        private readonly WindowsIdentitySafeHandle _identityToken = new WindowsIdentitySafeHandle();

        private WindowsIdentity _lastRecordedIdentity;

        public void OnClientImpersonated()
        {
            var currentIdentity = WindowsIdentity.GetCurrent();

            if (_lastRecordedIdentity != null && _lastRecordedIdentity.User?.Value == currentIdentity.User?.Value)
                return;

            _lastRecordedIdentity = currentIdentity;
            IntPtr handle;
            NativeMethods.DuplicateToken(_lastRecordedIdentity.Token, NativeMethods.SecurityImpersonationLevel.SecurityImpersonation, out handle);
            _identityToken.Handle = handle;
        }

        public WindowsImpersonationContext Impersonate()
        {
            return WindowsIdentity.Impersonate(_identityToken.Handle);
        }

        private static class NativeMethods
        {
            public enum SecurityImpersonationLevel
            {
                SecurityImpersonation = 2
            }

            [DllImport("advapi32.dll", SetLastError = true)]
            public static extern bool DuplicateToken(IntPtr existingTokenHandle, SecurityImpersonationLevel level, out IntPtr duplicateTokenHandle);
        }        
    }
}