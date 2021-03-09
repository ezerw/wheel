import React from 'react';
import { useGoogleAuth } from '../google-auth';
import {Image} from "react-bootstrap";

const LogoutButton = () => {
  const { signOut, googleUser } = useGoogleAuth();

  return (
    <div>
      <Image roundedCircle width={40} className="mr-2" src={googleUser.profileObj.imageUrl} />
      <button className="btn btn-dark" onClick={signOut}>Logout</button>
    </div>

  );
};

export default LogoutButton;