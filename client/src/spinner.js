import * as d3 from 'd3'

const getRandomInt = (min, max) => {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min)) + min;
}

const spinner = (people, refChart, refResult, refTrigger, handleShow) => {
  const padding = {top: 20, right: 20, bottom: 20, left: 20};
  const width = 500 - padding.left - padding.right;
  const height = 500 - padding.top - padding.bottom;
  const radius = Math.min(width, height) / 2;
  const spins = 3;
  const color = d3.scaleOrdinal([
    "#F87171",
    "#FCD34D",
    "#60A5FA",
    "#FBBF24",
    "#34D399",
    "#818CF8",
    "#F472B6",
    "#93C5FD"
  ]);

  let svg = d3.select(refChart).selectAll('svg').data([null]);
  svg = svg
    .enter().append('svg')
    .merge(svg)
    .data([people])
    .attr('width', 500)
    .attr('height', 500);

  const container = svg.append('g')
    .attr('class', 'chartcontainer')
    .attr('transform', `translate(${width / 2 + padding.left},${height / 2 + padding.top})`);

  const wheel = container.append('g')
    .attr('class', 'wheel');

  const pie = d3.pie().sort(null).value(function (d) {
    return 1;
  });

  const arc = d3.arc()
    .innerRadius(0)
    .outerRadius(radius);

  const arcs = wheel.selectAll('g.slice')
    .data(pie)
    .enter()
    .append('g')
    .attr('class', 'slice');

  arcs.append('path')
    .attr('fill', function (d, i) {
      return color(i);
    })
    .attr('d', function (d) {
      return arc(d);
    });

  arcs.append("text").attr("transform", function (d) {
    d.innerRadius = 0;
    d.outerRadius = radius;
    d.angle = (d.startAngle + d.endAngle) / 2;
    return `rotate(${(d.angle * 180 / Math.PI - 90)})translate(${d.outerRadius - 10})`;
  })
    .attr("text-anchor", "end")
    .text(function (d, i) {
      return people[i].first_name;
    })
    .style('font-size', '30px');

  // arrow
  svg.append('g')
    .attr('class', 'arrow')
    .attr('transform', `translate(${(width + padding.left + padding.right) / 2 - 15}, 12)`)
    .append('path')
    .attr('d', `M0 0 H30 L 15 ${Math.sqrt(3) / 2 * 35}Z`)
    .style('fill', '#000809');

  const degrees = spins * 360;

  const spin = (e) => {
    e.preventDefault()
    const piedegree = 360 / people.length;
    const randomAssetIndex = getRandomInt(0, people.length);
    const randomPieMovement = getRandomInt(1, piedegree);

    let rotation = (people.length - randomAssetIndex) * piedegree - randomPieMovement + degrees;

    const rotTween = () => {
      let i = d3.interpolate(0, rotation);
      return function (t) {
        return `rotate(${i(t)})`;
      };
    }

    wheel.transition()
      .duration(2000)
      .attrTween('transform', rotTween)
      .ease(d3.easeCircleOut)
      .on('end', () => handleShow(people[randomAssetIndex]))

  }

  const trigger = d3.select(refTrigger);
  trigger.on('click', (e) => spin(e));
}

export default spinner