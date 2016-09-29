using System;
using System.Collections.Generic;
using System.Linq;

namespace Docker.Cli
{
    public class Flags
    {
        private readonly IEnumerable<string> _args;

        public Flags(IEnumerable<string> args)
        {
            _args = args;
        }

        public bool Contains(string name)
        {
            return _args.Contains(name);
        }

        public void If(string name, Action action)
        {
            if (Contains(name))
            {
                action.Invoke();
            }
        }

        public void With(string name, Action<string> action)
        {
            var value = Value(name);
            if (!string.IsNullOrEmpty(value))
            {
                action.Invoke(value);
            }
        }

        private string Value(string name)
        {
            return (from _ in _args where _.StartsWith(name) select ExtractValue(_)).FirstOrDefault();
        }

        private static string ExtractValue(string flag)
        {
            return flag.Split(new[] {'='}, 2)[1];
        }
    }
}