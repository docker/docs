using System;

namespace Docker.Core.Pipe
{
    [Serializable]
    public class PipeRequest
    {
        public string Action { get; }
        public object[] Parameters { get; }

        public PipeRequest(string action, params object[] parameters)
        {
            Action = action;
            Parameters = parameters;
        }
    }
}