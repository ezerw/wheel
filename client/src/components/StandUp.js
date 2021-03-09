import React, {useState, useEffect} from 'react'
import Card from "react-bootstrap/Card";
import ListGroup from "react-bootstrap/ListGroup";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";

const StandUp = ({people}) => {
  const [sortType, setSortType] = useState('first_name');
  const [data, setData] = useState([]);

  const shuffle = (array) => {
    for (let i = array.length - 1; i > 0; i--) {
      let j = Math.floor(Math.random() * (i + 1));
      let temp = array[i];
      array[i] = array[j];
      array[j] = temp;
    }
    return array;
  }

  const randomize = () => {
    const shuffled = shuffle([...people])
    setData(shuffled)
  }

  useEffect(() => {
    const sortPeople = type => {
      const types = {
        first_name: 'first_name',
        last_name: 'last_name',
      };
      const sortProperty = types[type];
      const sorted = [...people].sort((a, b) =>
        (a[sortProperty] > b[sortProperty]) ? 1 : ((b[sortProperty] > a[sortProperty]) ? -1 : 0));
      setData(sorted)
    }
    sortPeople(sortType)
  }, [people, sortType]);

  return (
    <div className="StandUp">

      <Row>
        <Col md={4}>
          <Card.Title>Standup order</Card.Title>
        </Col>
        <Col className="text-right">
          <Form>
            <Form.Row className="justify-content-end">
              <Col xs="auto">
                <Form.Control as="select" size="sm" onChange={(e) => setSortType(e.target.value)}>
                  <option value="first_name">First name</option>
                  <option value="last_name">Last name</option>
                </Form.Control>
              </Col>
              <Col xs="auto">
                <Button variant="outline-dark" size="sm" onClick={randomize}>Random</Button>
              </Col>
            </Form.Row>
          </Form>
        </Col>
      </Row>
      <ListGroup variant="flush">
        {data.map(person => (
          <ListGroup.Item className="standup-item" key={person.id}>
            {person.first_name} {person.last_name}
          </ListGroup.Item>
        ))}
      </ListGroup>

    </div>
  );
}

export default StandUp;