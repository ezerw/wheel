import React from 'react';
import LoginButton from './LoginButton';
import {useGoogleAuth} from "../google-auth";

const Home = () => {
  const {isInitialized} = useGoogleAuth();
  return (
    isInitialized && (
      <div className="h-100 d-flex justify-content-center align-items-center">
        <div className="w-25 d-flex flex-column align-items-center">
          <h2 className="weird-font big">Wheel</h2>
          <LoginButton/>
        </div>
      </div>
    )
  );
};

export default Home;