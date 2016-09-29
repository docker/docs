using System;
using System.Runtime.Serialization;
using System.Security.Permissions;

namespace Docker.Core.backend
{
    [Serializable]
    public class BackendWarnException : Exception
    {
        public BackendWarnException() : base("")
        {
        }

        public BackendWarnException(string message) : base(message)
        {
        }

        [SecurityPermission(SecurityAction.Demand, SerializationFormatter = true)]
        protected BackendWarnException(SerializationInfo info, StreamingContext context) : base(info, context)
        {
        }
    }
}
