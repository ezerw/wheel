import React, {useState} from 'react';
import axios from "axios";
import Button from "react-bootstrap/Button";
import BModal from "react-bootstrap/Modal";

const Modal = ({team, selected, show, handleSaved, whoIsNext}) => {
  const [saving, setSaving] = useState(false)

  const handleResponse = (o) => {
    setSaving(false)
    whoIsNext(o)
    handleSaved()
  }

  const handleSubmit = () => {
    setSaving(true)

    const options = {headers: {'Content-Type': 'application/json'}}
    const body = {person_id: selected.id}

    axios.post(`${process.env.REACT_APP_API_URL}/teams/${team.id}/turns`, body, options)
      .then(res => handleResponse(res.data.data))
      .catch(err => console.log(err.response))
  }

  return (
    <BModal
      show={show}
      onHide={handleSaved}
      backdrop="static"
      keyboard={false}
      animation={false}
      size="md"
    >
      <BModal.Body>
        <div className="flex justify-content-center">
          <h1 className="m-0 weird-font big text-center">{selected.first_name}</h1>
        </div>
      </BModal.Body>
      <BModal.Footer>
        <Button variant="outline-danger" size="sm" onClick={handleSaved}>
          Cancel
        </Button>
        <Button
          variant="success"
          size="sm"
          onClick={!saving ? handleSubmit : null}
          disabled={saving}
        >
          {`Save ${selected.first_name} as next`}
        </Button>
      </BModal.Footer>
    </BModal>
  )
}
export default Modal;