import React from 'react';
import {render} from 'react-dom';
import {GoogleAuthProvider} from './google-auth';

import './App.scss'
import App from './App';

render(
  <GoogleAuthProvider>
    <App/>
  </GoogleAuthProvider>,
  document.getElementById('root')
);
