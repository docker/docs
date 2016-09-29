using System;

namespace Docker.Cli
{
    public interface IConsoleWriter
    {
        void WriteLine(string line);
    }

    public class ConsoleWriter : IConsoleWriter
    {
        public void WriteLine(string line)
        {
            Console.WriteLine(line);
        }
    }
}