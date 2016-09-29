using System;
using System.Diagnostics.CodeAnalysis;
using Docker.Core.di;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class SingletonsTest
    {
        private Singletons _singletons;

        [SetUp]
        public void CreateSingletons()
        {
            _singletons = new Singletons();
        }

        [TearDown]
        public void DisposeSingletons()
        {
            _singletons?.Dispose();
        }

        [Test]
        public void CreateSingletonInstance()
        {
            var value = _singletons.Get<Value>();
            var secondValue = _singletons.Get<Value>();

            Check.That(value).IsNotNull().And.IsSameReferenceThan(secondValue);
        }

        [Test]
        public void ConstructorInjection()
        {
            var value = _singletons.Get<Value>();
            var holder = _singletons.Get<ValueHolder>();

            Check.That(holder.Value).IsSameReferenceThan(value);
        }

        [Test]
        public void OnePublicConstructor()
        {
            var value = _singletons.Get<OnePublicConstructor>();
            var secondValue = _singletons.Get<OnePublicConstructor>();

            Check.That(value).IsSameReferenceThan(secondValue);
        }

        [Test]
        public void NoPublicConstructor()
        {
            Check.ThatCode(() => _singletons.Get<NoPublicConstructor>())
                .Throws<ArgumentException>()
                .WithMessage("Class Docker.Tests.NoPublicConstructor has no public constructor");
        }

        [Test]
        public void TooManyConstructors()
        {
            Check.ThatCode(() => _singletons.Get<TooManyConstructors>())
                .Throws<ArgumentException>()
                .WithMessage("Class Docker.Tests.TooManyConstructors has too many public constructors");
        }

        [Test]
        public void Cycle()
        {
            Check.ThatCode(() => _singletons.Get<Cycle>())
                .Throws<StackOverflowException>()
                .WithMessage("Class Docker.Tests.Cycle has a cycle in its dependencies");
        }

        [Test]
        public void CreateSingletonInstanceByInterface()
        {
            using (var singletons = new Singletons(new TestModule()))
            {
                var value = singletons.Get<IInterface>();
                var secondValue = singletons.Get<IInterface>();
                var byClass = singletons.Get<Value>();

                Check.That(value).IsSameReferenceThan(secondValue).And.IsSameReferenceThan(byClass);
            }
        }

        [Test]
        public void AutoFindSingleImplementationForInterface()
        {
            var value = _singletons.Get<IInterface>();
            var secondValue = _singletons.Get<IInterface>();
            var byClass = _singletons.Get<Value>();

            Check.That(value).IsSameReferenceThan(secondValue).And.IsSameReferenceThan(byClass);
        }

        [Test]
        public void RegisterConstant()
        {
            var constantValue = new Value();

            using (var singletons = new Singletons(new TestModuleWithConstant(constantValue)))
            {
                var value = singletons.Get<IInterface>();
                var secondValue = singletons.Get<IInterface>();

                Check.That(value).IsSameReferenceThan(constantValue).And.IsSameReferenceThan(secondValue);
            }
        }

        [Test]
        public void InjectTheSingletonsItself()
        {
            using (var singletons = new Singletons())
            {
                var value = singletons.Get<WithSingleton>();

                Check.That(value.Singletons).IsSameReferenceThan(singletons);
            }
        }
    }

    internal class TestModule : Module
    {
        protected override void Configure()
        {
            Bind<IInterface, Value>();
        }
    }

    internal class TestModuleWithConstant : Module
    {
        private readonly Value _value;

        public TestModuleWithConstant(Value value)
        {
            _value = value;
        }

        protected override void Configure()
        {
            Bind<IInterface>(_value);
        }
    }

    internal class Value : IInterface
    {
    }

    internal class ValueHolder
    {
        internal readonly Value Value;

        public ValueHolder(Value value)
        {
            Value = value;
        }
    }

    internal class Cycle
    {
        // ReSharper disable once UnusedParameter.Local
        public Cycle(Cycle value)
        {
        }
    }

    [SuppressMessage("ReSharper", "UnusedParameter.Local")]
    internal class TooManyConstructors
    {
        public TooManyConstructors(Value ignored)
        {
        }

        public TooManyConstructors(Value ignored1, Value ignored2)
        {
        }
    }

    internal class NoPublicConstructor
    {
        // ReSharper disable once EmptyConstructor
        internal NoPublicConstructor()
        {
        }
    }

    internal class OnePublicConstructor
    {
        public OnePublicConstructor(): this(true)
        {
        }

        // ReSharper disable once UnusedParameter.Local
        internal OnePublicConstructor(bool ignored)
        {
        }
    }

    internal interface IInterface
    {
    }

    internal class WithSingleton
    {
        public readonly Singletons Singletons;

        public WithSingleton(Singletons singletons)
        {
            Singletons = singletons;
        }
    }
}