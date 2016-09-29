using System;
using System.Management;

namespace Docker.Backend.Wmi
{
    public class WmiObject
    {
        internal WmiObject(ManagementObject target)
        {
            Target = target;
        }

        public ManagementObject Target { get; }
        public ManagementPath ClassPath => Target.ClassPath;

        public void Invoke(string methodName, params object[] arguments)
        {
            Target.InvokeMethod(methodName, arguments);
        }

        public WmiObject Set(string key, object value)
        {
            Target.SetPropertyValue(key, value);
            return this;
        }

        public string Get(string key)
        {
            try
            {
                return Target.Properties[key].Value.ToString();
            }
            catch (Exception)
            {
                return null;
            }
        }

        public string Text()
        {
            return Target.GetText(TextFormat.CimDtd20);
        }
    }
}