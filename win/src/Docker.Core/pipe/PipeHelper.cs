using System.IO;
using System.IO.Pipes;
using System.Security.AccessControl;
using System.Security.Principal;

namespace Docker.Core.Pipe
{
    public static class PipeHelper
    {
        public static NamedPipeServerStream NewServerStream(string pipeName)
        {
            var sid = new SecurityIdentifier(WellKnownSidType.WorldSid, null);
            var accessRule = new PipeAccessRule(sid, PipeAccessRights.FullControl, AccessControlType.Allow);

            var pipeSecurity = new PipeSecurity();
            pipeSecurity.AddAccessRule(accessRule);

            return new NamedPipeServerStream(
                pipeName,
                PipeDirection.InOut,
                NamedPipeServerStream.MaxAllowedServerInstances,
                PipeTransmissionMode.Byte,
                PipeOptions.None,
                512,
                512,
                pipeSecurity,
                HandleInheritability.None
                );
        }

        public static NamedPipeClientStream NewClientStream(string pipeName)
        {
            return new NamedPipeClientStream(".", pipeName, PipeDirection.InOut, PipeOptions.None);
        }
        public static NamedPipeClientStream NewImpersonatedClientStream(string pipeName)
        {
            return new NamedPipeClientStream(".", pipeName, PipeDirection.InOut, PipeOptions.None, TokenImpersonationLevel.Impersonation);
        }
    }
}