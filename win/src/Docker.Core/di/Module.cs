using System;

namespace Docker.Core.di
{
    public abstract class Module
    {
        private Binder _binder;

        protected abstract void Configure();

        internal void Configure(Binder binder)
        {
            _binder = binder;
            try
            {
                Configure();
            }
            finally
            {
                _binder = null;
            }
        }

        protected void Bind<T, TI>() where TI : T
        {
            _binder.Bind<T, TI>();
        }

        protected void Bind<T>(T value)
        {
            _binder.Bind(value);
        }

        internal void Merge(Binder other)
        {
            _binder.Merge(other);
        }

        public static ModuleOverride Override(Module module)
        {
            return new ModuleOverride(module);
        }
    }

    public class ModuleOverride
    {
        private readonly Module _module;

        internal ModuleOverride(Module module)
        {
            _module = module;
        }

        public Module With(Action<IBinder> other)
        {
            return new OverridenModule(_module, other);
        }
    }

    internal class OverridenModule : Module
    {
        private readonly Binder _firstBinder;
        private readonly Module _module;
        private readonly Action<IBinder> _other;
        private readonly Binder _secondBinder;

        public OverridenModule(Module module, Action<IBinder> other)
        {
            _module = module;
            _other = other;
            _firstBinder = new Binder();
            _secondBinder = new Binder();
        }

        protected override void Configure()
        {
            _module.Configure(_firstBinder);
            _other(_secondBinder);

            Merge(_firstBinder);
            Merge(_secondBinder);
        }
    }
}