using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;

namespace Docker.Core.di
{
    public class Singletons : IDisposable
    {
        private readonly IDictionary<Type, object> _instances;
        private readonly Binder _binder;

        public Singletons()
        {
            _instances = new Dictionary<Type, object>();
            _binder = new Binder();
            _binder.InstancesPerType[GetType()] = this;
        }

        public Singletons(Module module) : this()
        {
            module.Configure(_binder);
        }

        public T Get<T>() where T : class
        {
            return (T) Get(typeof(T), 0);
        }

        private object Get(Type type, int depth)
        {
            var constant = _binder.GetValueFor(type);
            if (constant != null)
            {
                return constant;
            }

            type = _binder.GetImplementationFor(type);

            object instance;

            if (!_instances.TryGetValue(type, out instance))
            {
                instance = NewInstance(type, depth);
                _instances.Add(type, instance);
            }

            return instance;
        }

        private object NewInstance(Type type, int depth)
        {
            if (depth > 100)
            {
                throw new StackOverflowException($"Class {type} has a cycle in its dependencies");
            }

            var constructors = type.GetConstructors();
            if (constructors.Length == 0)
            {
                throw new ArgumentException($"Class {type} has no public constructor");
            }
            if (constructors.Length > 1)
            {
                throw new ArgumentException($"Class {type} has too many public constructors");
            }

            var constructor = constructors[0];

            try
            {
                return constructor.Invoke(constructor.GetParameters()
                    .Select(_ => Get(_.ParameterType, depth + 1))
                    .ToArray());
            }
            catch (ArgumentException e)
            {
                throw new ArgumentException($"Unable to create an instance of {type}: {e.Message}");
            }
            catch (TargetInvocationException e)
            {
                throw new ArgumentException($"Unable to create an instance of {type}", e.InnerException);
            }
        }

        //https://msdn.microsoft.com/library/ms244737.aspx
        ~Singletons()
        {
            Dispose(false);
        }

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (!disposing) return;

            foreach (var instanceAndType in _instances)
            {
                var singleton = instanceAndType.Value;
                if (singleton != this)
                {
                    (singleton as IDisposable)?.Dispose();
                }
            }
            _instances.Clear();
        }
    }
}