using System;
using System.Runtime.Serialization;
using System.Security.Permissions;
using Docker.Core.Features;

namespace Docker.Core.backend
{
    [Serializable]
    public class BackendInstallFeatureException : Exception
    {
        public Feature[] Features { get; }

        public BackendInstallFeatureException(Feature[] features) : base("BackendInstallFeatureException")
        {
            Features = features;
        }

        [SecurityPermission(SecurityAction.Demand, SerializationFormatter = true)]
        public override void GetObjectData(SerializationInfo info, StreamingContext context)
        {
            base.GetObjectData(info, context);
            info.AddValue("Features", Features, typeof(Feature[]));
        }

        [SecurityPermission(SecurityAction.Demand, SerializationFormatter = true)]
        protected BackendInstallFeatureException(SerializationInfo info, StreamingContext context) : base(info, context)
        {
            Features = (Feature[])info.GetValue("Features", typeof(Feature[]));
        }
    }
}
