using System;
using System.Management;

namespace Docker.Backend.Wmi
{
    public class Wmi
    {
        private readonly ManagementScope _root;

        public Wmi(string root = "")
        {
            _root = new ManagementScope(root);
        }

        public WmiObject Create(ManagementPath path, string type)
        {
            var instance = new ManagementClass($"\\\\{path.Server}\\{path.NamespacePath}:{type}").CreateInstance();
            if (instance == null)
            {
                throw new ArgumentException($"Unable to create a new instance of {type}");
            }

            return new WmiObject(instance);
        }

        public WmiObject GetOrFail(string className, Predicate<ManagementBaseObject> filter = null)
        {
            var value = GetOrNull(className, filter);
            if (value != null)
            {
                return value;
            }

            throw new ArgumentException($"Not found {_root}/{className}");
        }

        public WmiObject GetOrNull(string className, Predicate<ManagementBaseObject> filter = null)
        {
            using (var searcher = new ManagementObjectSearcher(_root, new SelectQuery(className)))
            {
                foreach (var element in searcher.Get())
                {
                    if ((filter == null) || filter.Invoke(element))
                    {
                        return new WmiObject((ManagementObject)element);
                    }
                }
            }

            return null;
        }
    }
}