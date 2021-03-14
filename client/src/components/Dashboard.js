import React, {useCallback, useEffect, useState} from 'react';
import axios from "axios";
import {Navbar, Nav, Col, Row} from "react-bootstrap";

import LogoutButton from './LogoutButton';
import {useGoogleAuth} from "../google-auth";
import Wheel from './Wheel';
import StandUp from './StandUp';
import Stats from './Stats'
import Next from "./Next";
import LastThree from "./LastThree";

const Dashboard = () => {
  const {isInitialized, signOut} = useGoogleAuth();
  const [team, setTeam] = useState();
  const [next, setNext] = useState();
  const [lastThree, setLastThree] = useState();

  const handleSignOut = useCallback(() => {
    localStorage.removeItem('token')
    signOut()
  }, [signOut]);

  const colors = [
    "#F87171",
    "#FCD34D",
    "#60A5FA",
    "#FBBF24",
    "#34D399",
    "#818CF8",
    "#F472B6",
    "#93C5FD"
  ];

  useEffect(() => {
    const getTeam = () => {
      axios.get(`${process.env.REACT_APP_API_URL}/teams/1`)
        .then(res => {
          setTeam(res.data.data)
        })
        .catch(err => {
          console.log("Error: ", err.response)
          if (err.response.status === 401) {
            handleSignOut()
          }
        })
    }
    getTeam()
  }, [handleSignOut])

  useEffect(() => {
    const getNextTurn = () => {
      axios.get(`${process.env.REACT_APP_API_URL}/teams/1/turns?limit=4&order=desc`)
        .then(res => {
          const data = res.data.data
          const next = data.shift()
          setNext(next)
          setLastThree(data)
        })
        .catch(err => {
          console.log("Error: ", err.response)
          if (err.response.status === 401) {
            handleSignOut()
          }
        })
    }
    getNextTurn()
  }, [handleSignOut])

  return (
    isInitialized && (
      <>
        <Navbar bg="dark" variant="dark" expand="lg" className="mb-4">
          <div className="container">
            <Navbar.Brand href="">
              <h2 className="m-0">{team && team.name}</h2>
            </Navbar.Brand>
            <Navbar.Toggle aria-controls="basic-navbar-nav"/>
            <Navbar.Collapse id="basic-navbar-nav">
              <Nav className="mr-auto"/>
              <LogoutButton/>
            </Navbar.Collapse>
          </div>
        </Navbar>
        <div className="container">
          <Row className="my-2">
            <Col md={3}>
              {next && <Next next={next}/>}
            </Col>
            <Col md={3}>
              {lastThree && <LastThree lastThree={lastThree}/>}
            </Col>
            <Col md={6}>
              {team && <Stats colors={colors} team={team}/>}
            </Col>
          </Row>
          <Row className="my-4">
            <Col md={6}>
              {team && <Wheel colors={colors} team={team} whoIsNext={setNext}/>}
            </Col>
            <Col md={6}>
              {team && <StandUp people={team.people}/>}
            </Col>
          </Row>
        </div>
      </>
    )
  );
};

export default Dashboard;