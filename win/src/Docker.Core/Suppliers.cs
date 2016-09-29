using System;

namespace Docker.Core
{
    public interface ISupplier<out T>
    {
        T Get();
    }

    public static class Suppliers<T>
    {
        public static ISupplier<T> Memoize(Func<T> factory)
        {
            return new MemoizeSupplier<T>(factory);
        }
    }

    internal class MemoizeSupplier<T> : ISupplier<T>
    {
        private readonly Func<T> _factory;
        private readonly object _lock;

        private T _instance;

        internal MemoizeSupplier(Func<T> factory)
        {
            _factory = factory;
            _lock = new object();
        }

        public T Get()
        {
            lock (_lock)
            {
                if (_instance == null)
                {
                    _instance = _factory.Invoke();
                }
            }

            return _instance;
        }
    }
}
