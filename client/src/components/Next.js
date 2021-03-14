import React from 'react';
import Card from 'react-bootstrap/Card';
import {format, parseISO} from 'date-fns';

const Next = ({next}) => {
  const date = format(parseISO(next.date), 'EE dd, MMM y')
  return (
    <Card className="shadow-sm border-0">
      <Card.Body>
        <Card.Title>Next Turn</Card.Title>
        <div className="py-3">
          <h3 className="m-0 text-center big weird-font">{next.person.first_name}</h3>
          <h4 className="text-center">{date}</h4>
        </div>
      </Card.Body>
    </Card>
  )
}

export default Next;