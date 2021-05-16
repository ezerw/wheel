import React from 'react';
import { useGoogleAuth } from '../google-auth';
import Image from "react-bootstrap/Image";
import NavDropdown from 'react-bootstrap/NavDropdown'

const LogoutButton = () => {
  const { signOut, googleUser } = useGoogleAuth();
  console.log(googleUser);
  return (
    <>
    <Image roundedCircle width={40} className="mr-2" src={googleUser.profileObj.imageUrl} />
    <NavDropdown title={googleUser.profileObj.givenName} id="basic-nav-dropdown">
      <NavDropdown.Item href="#action/3.1">Token Wallet</NavDropdown.Item>
      <NavDropdown.Divider />
      <NavDropdown.Item href="#" onClick={signOut}>Log out</NavDropdown.Item>
    </NavDropdown>
    </>
  );
};

export default LogoutButton;