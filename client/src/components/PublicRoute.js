import React from 'react';
import {Route, Redirect} from 'react-router-dom';
import {useGoogleAuth} from "../google-auth";

const PublicRoute = ({component: Component, ...rest}) => {

  const {isSignedIn} = useGoogleAuth();

  return (
    <Route {...rest} render={props => (
      !isSignedIn ?
        <Component {...props} /> :
        <Redirect exact to="/dashboard"/>
    )}/>
  );
};

export default PublicRoute;