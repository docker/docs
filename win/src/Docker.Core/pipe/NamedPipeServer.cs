using System;
using System.Collections.Generic;
using System.IO;
using System.IO.Pipes;
using System.Runtime.Serialization.Formatters.Binary;
using System.Threading;
using Docker.Core.backend;

namespace Docker.Core.Pipe
{
    public interface IClientIdentityCallback
    {
        void OnClientImpersonated();
    }

    public class NamedPipeServer
    {
        private const string StopMessage = "<<STOP>>";

        private readonly Logger _logger;
        private readonly string _pipeName;
        private readonly IDictionary<string, Func<object[], object>> _actions;
        private readonly IClientIdentityCallback _clientIdentityCallback;
        private bool _stopped;

        public NamedPipeServer(string pipeName, IClientIdentityCallback clientIdentityCallback = null)
        {
            _logger = new Logger(GetType());
            _pipeName = pipeName;
            _actions = new Dictionary<string, Func<object[], object>>();
            _clientIdentityCallback = clientIdentityCallback;
        }

        public void Register(string action, Func<object[], object> func)
        {
            _actions.Add(action, func);
        }

        public void Register(string action, Action<object[]> func)
        {
            _actions.Add(action, parameters =>
            {
                func(parameters);
                return "";
            });
        }

        public void Run()
        {
            while (!_stopped)
            {
                NamedPipeServerStream pipeServer;
                try
                {
                    pipeServer = PipeHelper.NewServerStream(_pipeName);
                }
                catch (Exception e)
                {
                    _logger.Error($"Unable to create a pipe: {e.Message} {e.StackTrace}");
                    continue;
                }

                try
                {
                    pipeServer.WaitForConnection();
                }
                catch (Exception e)
                {
                    _logger.Error($"Unable to connect: {e.Message} {e.StackTrace}");
                    continue;
                }

                new Thread(() =>
                {
                    try
                    {
                        using (var server = pipeServer)
                        {
                            var bf = new BinaryFormatter();
                            var request = (PipeRequest)bf.Deserialize(server);

                            if (request.Action == StopMessage)
                            {
                                _stopped = true;
                                return;
                            }

                            if (_clientIdentityCallback != null)
                            {
                                server.RunAsClient(() => _clientIdentityCallback.OnClientImpersonated());
                            }

                            var response = RunAction(request.Action, request.Parameters);
                            bf.Serialize(server, response);
                            server.Flush();
                            server.WaitForPipeDrain();
                        }
                    }
                    catch (Exception e)
                    {
                        _logger.Error($"Pipe failure: {e.Message} {e.StackTrace}");
                    }
                }).Start();
            }
        }

        private PipeResponse RunAction(string action, object[] parameters)
        {
            _logger.Info($"{action}({string.Join(", ", parameters)})");

            if (!_actions.ContainsKey(action))
            {
                _logger.Error($"Unknown request. Don't know what to do with {action}");

                return ErrorResponse(new DockerException($"Unknown action {action}"));
            }

            try
            {
                var answer = _actions[action](parameters);

                _logger.Info($"{action} done.");

                return new PipeResponse(true, answer);
            }
            catch (BackendQuitException exception)
            {
                return ErrorResponse(exception);
            }
            catch (BackendWarnException exception)
            {
                return ErrorResponse(exception);
            }
            catch (BackendInstallFeatureException exception)
            {
                return ErrorResponse(exception);
            }
            catch (Exception exception)
            {
                _logger.Error($"Unable to execute {action}: {exception.Message} {exception.StackTrace}");

                return ErrorResponse(new BackendException(exception));
            }
        }

        public void Stop()
        {
            try
            {
                var request = new PipeRequest(StopMessage);

                using (var client = PipeHelper.NewClientStream(_pipeName))
                {
                    client.Connect();
                    _logger.Info("Sending Stop Message...");
                    var bf = new BinaryFormatter();
                    bf.Serialize(client, request);
                    client.Flush();
                    client.WaitForPipeDrain();
                }
            }
            catch (Exception ex)
            {
                _logger.Error(ex.Message);
                throw;
            }
        }

        private static PipeResponse ErrorResponse(Exception exception)
        {
            return new PipeResponse(false, CanBeDerialized(exception) ? exception : new DockerException(exception.Message));
        }

        private static bool CanBeDerialized(Exception exception)
        {
            var binary = new BinaryFormatter();
            var stream = new MemoryStream();

            try
            {
                binary.Serialize(stream, exception);
                binary.Deserialize(new MemoryStream(stream.ToArray()));

                return true;
            }
            catch
            {
                return false;
            }
        }
    }
}