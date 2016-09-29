using System;
using Docker.Core.Tracking;

namespace Docker.Core
{
    public interface INotifications
    {
        void NotifyError(string message, Exception exception);
        void NotifyError(Exception exception);
        void Notify(string message, string details, bool? useAnimation);
        void Notify(string message, string details, bool? useAnimation, AnalyticEvent mixpanelEvent);
        void SetStatus(string message, bool? useAnimation, AnalyticEvent mixpanelEvent);
        void SetStatus(string message, bool? useAnimation);
    }

    public class CliNotifications : INotifications
    {
        public void NotifyError(string message, Exception exception)
        {
            Console.Write($"{message} {exception.Message}");
            Environment.Exit(1);
        }

        public void NotifyError(Exception exception)
        {
            Console.Write(exception.Message);
            Environment.Exit(1);
        }

        public void Notify(string message, string details, bool? useAnimation)
        {
        }

        public void Notify(string message, string details, bool? useAnimation, AnalyticEvent mixpanelEvent)
        {
        }

        public void SetStatus(string message, bool? useAnimation, AnalyticEvent mixpanelEvent)
        {
        }

        public void SetStatus(string message, bool? useAnimation)
        {
        }
    }
}