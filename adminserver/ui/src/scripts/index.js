'use strict';

import './base.css';
import ReactDOM from 'react-dom';
import { Root } from './routes.js';

// Set the CSRF token globally before using the DTR JS SDK, ensuring that each
// request also sends along the token in the request headers.
let part = document.cookie.split(';').find(i => i.indexOf('csrf') >= 0) || 'a=null';
window.DTR_CSRF_TOKEN = part.split('=')[1];

ReactDOM.render(Root, document.getElementById('content'));
