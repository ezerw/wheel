import React from 'react';
import {useGoogleAuth} from '../google-auth';
import gLogo from '../Google__G__Logo.svg'

const LoginButton = () => {

  const {signIn} = useGoogleAuth();

  const handleSignIn = async () => {
    const res = await signIn()
    localStorage.setItem('token', res.accessToken)
  }

  return (
    <button className="btn btn-outline-primary" onClick={handleSignIn}>
      <img src={gLogo} width="20px" className="mr-2" alt=""/>Login with Google
    </button>
  );
};

export default LoginButton;