using System;
using System.Collections.Generic;
using System.Linq;

namespace Docker.Core.di
{
    public interface IBinder
    {
        void Bind<T, TI>() where TI : T;
        void Bind<T>(T value);
    }

    internal class Binder : IBinder
    {
        internal readonly IDictionary<Type, object> InstancesPerType;
        internal readonly IDictionary<Type, Type> TypesPerInterface;

        private static Type[] _allTypes;

        internal Binder()
        {
            TypesPerInterface = new Dictionary<Type, Type>();
            InstancesPerType = new Dictionary<Type, object>();
        }

        public void Bind<T, TI>() where TI : T
        {
            TypesPerInterface.Add(typeof(T), typeof(TI));
        }

        public void Bind<T>(T value)
        {
            var type = typeof(T);
            try
            {
                InstancesPerType.Add(type, value);
            }
            catch (ArgumentException)
            {
                throw new ArgumentException($"A binding already exists for {type}");
            }
        }

        internal void Merge(Binder otherBinder)
        {
            foreach (var entry in otherBinder.TypesPerInterface)
            {
                TypesPerInterface[entry.Key] = entry.Value;
            }

            foreach (var entry in otherBinder.InstancesPerType)
            {
                InstancesPerType[entry.Key] = entry.Value;
            }
        }

        internal Type GetImplementationFor(Type type)
        {
            if (!type.IsInterface)
            {
                return type;
            }

            Type implementationType;
            if (TypesPerInterface.TryGetValue(type, out implementationType))
            {
                return implementationType;
            }

            if (_allTypes == null)
            {
                _allTypes = AppDomain.CurrentDomain.GetAssemblies()
                    .SelectMany(s => s.GetTypes())
                    .ToArray();
            }

            var types = _allTypes
                .Where(mytype => mytype.GetInterfaces().Contains(type))
                .ToArray();
            if (types.Length == 1)
            {
                return types[0];
            }

            throw new ArgumentException($"Unknown implementation for {type}");
        }

        internal object GetValueFor(Type type)
        {
            object value;

            return InstancesPerType.TryGetValue(type, out value) ? value : null;
        }
    }
}