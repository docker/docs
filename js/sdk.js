(function (window) {

  'use strict';


  // Private functions

  function injectScript(src, cb) {
    var sj = document.createElement('script');
    sj.type = 'text/javascript';
    sj.async = true;
    sj.src = src;
    sj.addEventListener ? sj.addEventListener('load', cb, false) : sj.attachEvent('onload', cb);
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(sj, s);
  }


  // Process recaptcha input and inits SDK
  var verifyCallback = function(response) {
    var self = this;
    var data = encodeURIComponent('g-recaptcha-response') + '=' + encodeURIComponent(response);
    data += '&' + encodeURIComponent('session-duration') + '=' + encodeURIComponent('90m');
    sendRequest('POST', this.opts.baseUrl + '/', {headers:{'Content-type':'application/x-www-form-urlencoded'}}, data, function(resp) {
      //TODO handle errors
      if (resp.status == 200) {
        var sessionData = JSON.parse(resp.responseText);
        self.opts.baseUrl = 'http://' + sessionData.hostname;
        self.init(sessionData.session_id, self.opts);
        self.terms.forEach(function(term) {

          // Create terminals only for those elements that exist at least once in the DOM
          if (document.querySelector(term.selector)) {
            self.terminal(function() {
              //Remove captchas after initializing terminals;
              var captcha = document.querySelectorAll(term.selector + ' .captcha');
              for (var n=0; n < captcha.length; ++n) {
                captcha[n].parentNode.removeChild(captcha[n]);
              }
            });
          }

        });
      } else if (resp.status == 403) {
        // Forbidden, we need to display te captcha
        var term = window.pwd.terms[0];
        var els = document.querySelectorAll(term.selector);
        for (var n=0; n < els.length; ++n) {
          var captcha = document.createElement('div');
          captcha.className = 'captcha';
          els[n].appendChild(captcha);
          window.grecaptcha.render(captcha, {'sitekey': '6Ld8pREUAAAAAOkrGItiEeczO9Tfi99sIHoMvFA_', 'callback': verifyCallback.bind(window.pwd)});
        }
      };
    });
  };

  // register Recaptcha global onload callback
  window.onloadCallback = function() {
    //Register captcha only on the first term to avoid showing multiple times
    verifyCallback.call(window.pwd);
  }

  function registerInputHandlers(termName, instance) {
    var self = this;
    // Attach block actions
    var actions = document.querySelectorAll('code[class*="'+termName+'"]');
    for (var n=0; n < actions.length; ++n) {
      actions[n].onclick = function() {
        self.socket.emit('terminal in', instance.name, this.innerText);
      };
    }
  }

  function registerPortHandlers(termName, instance) {
    var self = this;
    // Attach block actions
    var actions = document.querySelectorAll('[data-term*="'+termName+'"]');
    for (var n=0; n < actions.length; ++n) {
      actions[n].onclick = function(evt) {
        evt.preventDefault();
        var port = this.getAttribute("data-port");
        if (port) {
          window.open('//pwd'+ instance.ip.replace(/\./g, "_") + '-' + port + '.' + self.opts.baseUrl.split('/')[2] + this.pathname);
        }
      };
    }
  }


  // PWD instance
  var pwd = function () {
    this.instances = {};
    this.instanceBuffer = {};
    return;
  };


  function setOpts(opts) {
    var opts = opts || {};
    this.opts = opts;
    this.opts.baseUrl = this.opts.baseUrl || 'http://labs.play-with-docker.com';
    this.opts.ports = this.opts.ports || [];
    this.opts.ImageName = this.opts.ImageName || '';
  }

  pwd.prototype.newSession = function(terms, opts) {
    setOpts.call(this, opts);
    terms = terms || [];
    if (terms.length > 0) {
      this.terms = terms;
      injectScript('https://www.google.com/recaptcha/api.js?onload=onloadCallback&render=explicit');
    } else {
      console.warn('No terms specified, nothing to do.');
    }
  };

  // your sdk init function
  pwd.prototype.init = function (sessionId, opts, callback) {
    var self = this;
    setOpts.call(this, opts);
    this.sessionId = sessionId;
    this.socket = io(this.opts.baseUrl, {path: '/sessions/' + sessionId + '/ws' });
    this.socket.on('terminal out', function(name ,data) {
      var instance = self.instances[name];
      if (instance && instance.terms) {
        instance.terms.forEach(function(term) {term.write(data)});
      } else {
        //Buffer the data if term is not ready
        if (self.instanceBuffer[name] == undefined) self.instanceBuffer[name] = '';
        self.instanceBuffer[name] += data;
      }
    });

    // Resize all terminals
    this.socket.on('viewport resize', function(cols, rows) {
      // Resize all terminals
      for (var name in self.instances) {
        self.instances[name].terms.forEach(function(term){
          term.resize(cols,rows);
        });
      };
    });

    // Handle window resizing
    window.onresize = function() {
      self.resize();
    };

    sendRequest('GET', this.opts.baseUrl + '/sessions/' + sessionId, undefined, undefined, function(response){
      var session = JSON.parse(response.responseText);
      for (var name in session.instances) {
        var i = session.instances[name];
        // Setup empty terms
        i.terms = [];
        self.instances[name] = i;
      }
      !callback || callback();
    });
  };


  pwd.prototype.resize = function() {
    var name = Object.keys(this.instances)[0]
    for (var n in this.instances) {
      for (var i = 0; i < this.instances[n].terms.length; i ++) {
        var term = this.instances[n].terms[i];
        var size = term.proposeGeometry();
        if (size.cols && size.rows) {
          return this.socket.emit('viewport resize', size.cols, size.rows);
        }
      }
    }
  };


  // I know, opts and data can be ommited. I'm not a JS developer =(
  // Data needs to be sent encoded appropriately
  function sendRequest(method, url, opts, data, callback) {
    var request = new XMLHttpRequest();
    request.open(method, url, true);

    if (opts && opts.headers) {
      for (var key in opts.headers) {
        request.setRequestHeader(key, opts.headers[key]);
      }
    }
    request.withCredentials = true;
    request.setRequestHeader('X-Requested-With', 'XMLHttpRequest')
    request.onload = function() {
      callback(request);
    };
    if (typeof(data) === 'object') {
      request.send(JSON.stringify(data));
    } else {
      request.send(data);
    }
  };

  pwd.prototype.createInstance = function(callback) {
    var self = this;
    //TODO handle http connection errors
    sendRequest('POST', self.opts.baseUrl + '/sessions/' + this.sessionId + '/instances', {headers:{'Content-type':'application/json'}}, {ImageName: self.opts.ImageName}, function(response) {
      if (response.status == 200) {
        var i = JSON.parse(response.responseText);
        i.terms = [];
        self.instances[i.name] = i;
        callback(undefined, i);
      } else if (response.status == 409) {
        var err = new Error();
        err.max = true;
        callback(err);
      } else {
        callback(new Error());
      }
    });
  }

  pwd.prototype.createTerminal = function(term, name) {
    var self = this;
    var i = this.instances[name];
    term.name = term.name || term.selector;

    var elements = document.querySelectorAll(term.selector);
    for (var n=0; n < elements.length; ++n) {
      var t = new Terminal({cursorBlink: false});
      t.open(elements[n]);
      t.on('data', function(d) {
        self.socket.emit('terminal in', i.name, d);
      });
      var size = t.proposeGeometry();
      self.socket.emit('viewport resize', size.cols, size.rows);
      i.terms.push(t);
    }


    registerPortHandlers.call(self, term.name, i)

    registerInputHandlers.call(self, term.name, i);



    if (self.instanceBuffer[name]) {
      //Flush buffer and clear it
      i.terms.forEach(function(term){
        term.write(self.instanceBuffer[name]);
      });
      self.instanceBuffer[name] = '';
    }

    return i.terms;
  }

  pwd.prototype.terminal = function(callback) {
    var self = this;
    this.createInstance(function(err, instance) {
      if (err && err.max) {
        !callback || callback(new Error("Max instances reached"))
        return
      } else if (err) {
        !callback || callback(new Error("Error creating instance"))
        return
      }

      var instance_number = instance.name[instance.name.length -1 ];
      self.createTerminal(self.terms[instance_number - 1], instance.name);


      !callback || callback(undefined, instance);

    });
  }



  // define your namespace myApp
  window.pwd = new pwd();

})(window, undefined);
