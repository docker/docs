using System;
using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class LogEntryTest
    {
        [Test]
        public void ToLogFileLine()
        {
            var date = new DateTime(2016, 1, 31, 13, 59, 1, 123);

            var entry = new LogEntry(date, "category", ChannelType.Info, "message");
            var line = entry.ToLine();

            Check.That(line).IsEqualTo("[13:59:01.123][category       ][Info   ] message");
        }
    }
}
