using System;
using System.Runtime.Serialization;
using System.Security.Permissions;

namespace Docker.Core.backend
{
    [Serializable]
    public class BackendQuitException : Exception
    {
        public BackendQuitException() : base("")
        {
        }

        public BackendQuitException(string message) : base(message)
        {
        }

        [SecurityPermission(SecurityAction.Demand, SerializationFormatter = true)]
        protected BackendQuitException(SerializationInfo info, StreamingContext context) : base(info, context)
        {
        }
    }
}
