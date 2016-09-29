using System;

namespace Docker.Core.Pipe
{
    [Serializable]
    public class PipeResponse
    {
        public bool Success { get; }
        public object ReturnValue { get; }

        public PipeResponse(bool success, object returnValue)
        {
            Success = success;
            ReturnValue = returnValue;
        }
    }
}