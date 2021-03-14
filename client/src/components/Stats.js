import React, {useState, useEffect} from 'react';
import {sub, format} from 'date-fns'
import DatePicker from "react-datepicker";
import {Doughnut} from 'react-chartjs-2'
import Card from "react-bootstrap/Card";
import Form from "react-bootstrap/Form";
import Col from "react-bootstrap/Col";
import Row from "react-bootstrap/Row";

import './Stats.scss'
import axios from "axios";

const Stats = ({colors, team}) => {
  const today = new Date();

  const [data, setData] = useState([])
  const [dateFrom, setDateFrom] = useState(sub(today, {months: 1}))
  const [dateTo, setDateTo] = useState(today)

  useEffect(() => {
    const formatRawData = res => {
      const count = (person, stats) => {
        return stats.reduce((c, curr) => {
          return person.id === curr.person.id ? c + 1 : c
        }, 0)
      }
      const stats = team.people.map(person => ({name: person.first_name, turns: count(person, res)}))
      setData(stats)
    }
    const fetchStats = () => {
      const f = 'yyyy-MM-dd'
      const query = `date_from=${format(dateFrom, f)}&date_to=${format(dateTo, f)}`

      axios.get(`${process.env.REACT_APP_API_URL}/teams/1/turns?${query}`)
        .then(res => formatRawData(res.data.data))
        .catch(err => console.log("Error: ", err.response))
    }
    fetchStats()
  }, [team, dateFrom, dateTo])

  const graphData = {
    labels: data.map(person => person.name),
    datasets: [
      {
        label: '# of Turns',
        data: data.map(person => person.turns),
        backgroundColor: colors,
        borderWidth: 0,
      },
    ],
  }

  return (
    <Card bg={'white'} className="shadow-sm border-0">
      <Card.Body>
        <Row className="mb-2">
          <Col>
            <Card.Title>Stats</Card.Title>
          </Col>
          <Col md={10}>
            <Form>
              <Form.Row className="justify-content-end">
                <Col xs="auto">
                  <DatePicker
                    className="form-control form-control-sm"
                    dateFormat="dd/MM/yyyy"
                    selected={dateFrom}
                    selectsStart
                    strictParsing
                    startDate={dateFrom}
                    endDate={dateTo}
                    onChange={setDateFrom}/>
                </Col>
                <Col xs="auto">
                  <DatePicker
                    className="form-control form-control-sm"
                    dateFormat="dd/MM/yyyy"
                    selected={dateTo}
                    selectsEnd
                    startDate={dateFrom}
                    endDate={dateTo}
                    minDate={dateFrom}
                    onChange={setDateTo}/>
                </Col>
                <Col xs="auto">

                </Col>
              </Form.Row>
            </Form>
          </Col>
        </Row>

        <Doughnut data={graphData}/>
      </Card.Body>
    </Card>
  )
}

export default Stats;