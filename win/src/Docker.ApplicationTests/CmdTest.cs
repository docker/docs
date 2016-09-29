using System;
using Docker.Core;
using NFluent;
using NUnit.Framework;

namespace Docker.Tests
{
    [TestFixture]
    public class CmdTest
    {
        private readonly Cmd _cmd = new Cmd();

        [Test]
        public void AdvancedErrorExitCode()
        {
            var result = _cmd.Run("xcopy", null, 0);

            Check.That(result.ExitCode).IsNotZero();
        }

        [Test]
        public void AdvancedTimeout()
        {
            var result = _cmd.Run("powershell.exe", "-ExecutionPolicy UNRESTRICTED -NoProfile -NonInteractive Sleep 2", 1000);

            Check.That(result.TimedOut).IsTrue();
        }

        [Test]
        public void AdvancedStdout()
        {
            var result = _cmd.Run("powershell.exe", "-ExecutionPolicy UNRESTRICTED -NoProfile -NonInteractive Write-Output \"line1\";Write-Output \"line2\"", 0);

            Check.That(result.StandardOutput.TrimEnd()).IsEqualTo("line1" + Environment.NewLine + "line2");
            Check.That(result.CombinedOutput.TrimEnd()).IsEqualTo("line1" + Environment.NewLine + "line2");
            Check.That(result.ErrorOutput).IsEmpty();
            Check.That(result.ExitCode).IsZero();
            Check.That(result.TimedOut).IsFalse();
        }

        [Test]
        public void AdvancedStderr()
        {
            var result = _cmd.Run("xcopy", null, 0);

            Check.That(result.ExitCode).IsEqualTo(4);
            Check.That(result.ErrorOutput).Contains("Invalid number of parameters");
            Check.That(result.TimedOut).IsFalse();
            Check.That(result.StandardOutput).Contains("0 File(s) copied");
        }

        [Test]
        public void AdvancedMix()
        {
            var result = _cmd.Run("cmd", "/c\"echo stdout && echo stderr 1>&2\"", 0);

            Check.That(result.StandardOutput.TrimEnd()).IsEqualTo("stdout");
            Check.That(result.ErrorOutput.TrimEnd()).IsEqualTo("stderr");
            Check.That(result.TimedOut).IsFalse();
            Check.That(result.ExitCode).IsZero();
        }
    }
}
