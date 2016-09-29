using System;
using System.Collections.Generic;
using System.Management.Automation;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class PowerShellTests
    {
        private readonly Core.PowerShell _powershell = new Core.PowerShell();

        [Test]
        public void CaptureOutput()
        {
            var output = _powershell.Output("Write-Output \"Hello\"");

            Check.That(output).IsEqualTo("Hello");
        }

        [Test]
        public void CaptureOutputMultipleLines()
        {
            var output = _powershell.Output("Write-Output \"Line1\"\nWrite-Output \"Line2\"");

            Check.That(output).IsEqualTo($"Line1{Environment.NewLine}Line2");
        }

        [Test]
        public void CaptureLines()
        {
            var lines = new List<string>();

            _powershell.Run("Write-Output \"Line1\"\nWrite-Output \"Line2\"\n", null, line => lines.Add(line));

            Check.That(lines).ContainsExactly("Line1", "Line2");
        }

        [Test]
        public void CaptureLinesBeforeError()
        {
            var lines = new List<string>();

            var exception = Assert.Throws<RuntimeException>(() =>
            {
                _powershell.Run("Write-Output \"Line1\"\nthrow \"ERROR\"\nWrite-Output \"Line2\"", null, line => lines.Add(line));
            });

            Check.That(exception.Message).IsEqualTo("ERROR");
            Check.That(lines).ContainsExactly("Line1");
        }

        [Test]
        public void CaptureError()
        {
            var exception = Assert.Throws<RuntimeException>( () =>
                _powershell.Output("throw \"ERROR\"")
            );

            Check.That(exception.Message).IsEqualTo("ERROR");
        }
    }
}