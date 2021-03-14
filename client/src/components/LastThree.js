import React from 'react';
import Card from "react-bootstrap/Card";
import {format, parseISO} from "date-fns";

const LastThree = ({lastThree}) => {
  console.log(lastThree)
  const listItems = lastThree.map(turn =>
    <li className="list-group-item d-flex justify-content-between align-items-center" key={turn.id}>
        {turn.person.first_name}
      <span className="badge badge-light badge-pill">
        {format(parseISO(turn.date), 'dd/MM/y')}
      </span>
    </li>
  );
  return (
    <Card className="shadow-sm border-0">
      <Card.Body>
        <Card.Title>Prev Three</Card.Title>
        <ul className="list-group list-group-flush">{listItems}</ul>
      </Card.Body>
    </Card>
  )
}

export default LastThree;