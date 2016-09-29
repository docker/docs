using System;
using System.Collections.Generic;
using System.Windows.Input;

namespace Docker.WPF
{
    public class Code
    {
        private readonly List<Key> _codeKeys = new List<Key>
        {
          Key.Up, Key.Up, Key.Down, Key.Down, Key.Left, Key.Right, Key.Left, Key.Right, Key.B, Key.A
        };

        private int _match;

        public void OnKeyUp(KeyEventArgs keyEvent, Action onSuccessfullCode)
        {
            if (keyEvent.Key == _codeKeys[_match])
            {
                if (++_match < _codeKeys.Count) return;
                _match = 0;
                onSuccessfullCode.Invoke();
            }
            else if (_match > 0 && keyEvent.Key != _codeKeys[_match])
            {
                _match = 0;
            }
        }
    }
}
