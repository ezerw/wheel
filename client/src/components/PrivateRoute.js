import React from 'react';
import {Route, Redirect} from 'react-router-dom';
import {useGoogleAuth} from "../google-auth";

const PrivateRoute = ({component: Component, ...rest}) => {
  const {isSignedIn} = useGoogleAuth();

  return (
    <Route {...rest} render={props => (
      isSignedIn ?
        <Component {...props} /> :
        <Redirect exact from="/dashboard" to="/"/>
    )}/>
  );
};

export default PrivateRoute;