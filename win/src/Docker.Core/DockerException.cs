using System;
using System.Runtime.Serialization;

namespace Docker.Core
{
    [Serializable]
    public class DockerException : Exception
    {
        public DockerException(string message) : base(message)
        {
        }

        public DockerException(string message, Exception ex) : base(message, ex)
        {
        }

        protected DockerException(SerializationInfo info, StreamingContext context) : base(info, context)
        {
        }
    }
}