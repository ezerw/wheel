import React, {useRef, useEffect, useState} from 'react'
import party from 'party-js'
import Modal from './Modal'
import spinner from '../spinner'

const Wheel = ({team, whoIsNext}) => {
  const [show, setShow] = useState(false);
  const [selected, setSelected] = useState('');

  const refWheel = useRef(null)
  const refResult = useRef(null)
  const refTrigger = useRef(null)

  const handleSaved = () => {
    setShow(false);
  }
  const throwConfetti = () => party.element(document.body, {
    color: ["#F87171", "#FCD34D", "#60A5FA", "#FBBF24", "#34D399", "#818CF8", "#F472B6", "#93C5FD"],
    count: party.variation(150, 0.5),
    size: party.minmax(6, 10),
    velocity: party.minmax(-300, -600)
  });

  useEffect(() => {
    const handleShow = (s) => {
      setSelected(s)
      throwConfetti()
      setShow(true);
    }
    spinner(team.people, refWheel.current, refResult.current, refTrigger.current, handleShow)
  }, [team])

  return (
    <>
      <div className="wheel weird-font" ref={refWheel}>
        <div className="spin">
          <a href="/#" ref={refTrigger}>Spin</a>
        </div>
      </div>

      <Modal
        team={team}
        selected={selected}
        show={show}
        handleSaved={handleSaved}
        whoIsNext={whoIsNext}/>
    </>
  );
}

export default Wheel;