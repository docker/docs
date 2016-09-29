using System;
using System.Runtime.Serialization;

namespace Docker.Core.backend
{
    [Serializable]
    public class BackendException : Exception
    {
        public BackendException(Exception ex) : base(ex.Message, ex)
        {
        }

        protected BackendException(SerializationInfo info, StreamingContext context) : base(info, context)
        {
        }
    }
}
