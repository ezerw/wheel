import React, {useRef, useEffect, useState} from 'react'
import party from 'party-js'
import Card from 'react-bootstrap/Card';
import Modal from './Modal'
import spinner from '../spinner'

const Wheel = ({colors, team, whoIsNext}) => {
  const [show, setShow] = useState(false);
  const [selected, setSelected] = useState('');

  const refWheel = useRef(null)
  const refResult = useRef(null)
  const refTrigger = useRef(null)

  useEffect(() => {
    const handleShow = (s) => {
      setSelected(s)
      party.element(document.body, {
        color: colors,
        count: party.variation(150, 0.5),
        size: party.minmax(6, 10),
        velocity: party.minmax(-300, -600)
      });
      setShow(true);
    }

    spinner(
      colors,
      team.people,
      refWheel.current,
      refResult.current,
      refTrigger.current,
      handleShow
    )
  }, [colors, team])

  return (
    <>
      <Card className="shadow-sm">
        <Card.Body>
          <Card.Title>Wheel</Card.Title>
          <div className="wheel weird-font" ref={refWheel}>
            <div className="spin">
              <a href="/#" ref={refTrigger}>Spin</a>
            </div>
          </div>
        </Card.Body>
      </Card>

      <Modal
        team={team}
        selected={selected}
        show={show}
        handleSaved={() => setShow(false)}
        whoIsNext={whoIsNext}/>
    </>
  );
}

export default Wheel;