using System;
using System.Runtime.Serialization;

namespace Docker.Backend
{
    [Serializable]
    public class HyperVException : Exception
    {
        public HyperVException(string message) : base(message)
        {
        }

        public HyperVException(string message, Exception ex) : base(message, ex)
        {
        }

        protected HyperVException(SerializationInfo info, StreamingContext context) : base(info, context)
        {
        }
    }
}
