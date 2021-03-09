import React, {useCallback, useEffect, useState} from 'react';
import axios from "axios";
import dayjs from 'dayjs';
import { Icon, InlineIcon } from "@iconify/react";
import barChart from '@iconify-icons/emojione-v1/bar-chart';
import constructionIcon from '@iconify-icons/emojione-v1/construction';
import {Navbar, Nav, Card, Col, Row} from "react-bootstrap";

import LogoutButton from './LogoutButton';
import {useGoogleAuth} from "../google-auth";
import Wheel from "./Wheel";
import StandUp from "./StandUp";

const Dashboard = () => {
  const {isInitialized, signOut} = useGoogleAuth();
  const [team, setTeam] = useState()
  const [next, setNext] = useState()

  const handleSignOut = useCallback(() => {
    localStorage.removeItem('token')
    signOut()
  }, [signOut])


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
      axios.get(`${process.env.REACT_APP_API_URL}/teams/1/turns?limit=1&order=desc`)
        .then(res => {
          setNext(res.data.data[0])
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
              <Nav className="mr-auto">
              </Nav>
              <LogoutButton/>
            </Navbar.Collapse>
          </div>
        </Navbar>
        <div className="container">
          <Row className="my-2">
            <Col md={4}>
              <Card className="shadow-sm border-0" bg={'info'} text={'white'}>
                <Card.Body>
                  <Card.Title className="text-center">{
                    next && dayjs(next.date).format("dddd DD, MMMM YYYY")
                  }</Card.Title>
                  <div className="text-center">
                    {!next ? 'Loading...' : (
                      <h3 className="m-0 big weird-font">{next.person.first_name}</h3>
                    )}
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col md={4}>
              <Card bg={'white'} className="shadow-sm border-0">
                <Card.Body>
                  <Card.Title>
                    Something
                  </Card.Title>
                  Maybe a graph? <InlineIcon width={42} icon={barChart} />
                </Card.Body>
              </Card>
            </Col>
            <Col md={4}>
              <Card bg={'white'} className="shadow-sm border-0">
                <Card.Body>
                  <Card.Title>
                    Something else
                  </Card.Title>
                  I don't know what else <InlineIcon width={42} icon={constructionIcon} />
                </Card.Body>
              </Card>
            </Col>
          </Row>
          <Row className="my-4">
            <Col md={6}>
              <Card className="shadow-sm">
                <Card.Body>
                  <Card.Title>
                    <h4>Wheel</h4>
                  </Card.Title>
                  {team && <Wheel team={team} whoIsNext={setNext} />}
                </Card.Body>
              </Card>
            </Col>
            <Col md={6}>
              <Card className="shadow-sm">
                <Card.Body>
                  <Card.Title>
                    <h4>StandUp</h4>
                  </Card.Title>
                  {team && <StandUp people={team.people} />}
                </Card.Body>
              </Card>
            </Col>
          </Row>
        </div>
      </>
    )
  );
};

export default Dashboard;