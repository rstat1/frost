import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import * as c3 from 'c3';
import * as d3 from 'd3';

@Component({
	selector: 'app-chart',
	templateUrl: './chart.component.html',
	styleUrls: ['./chart.component.css'],
})
export class ChartComponent implements OnInit {
/*
blue: #0088ce,
red: #cc0000,
orange: #ec7a08,
green: #3f9c35,
 */
	constructor() { }

	ngOnInit() {
		setTimeout(() => {
			let chat = c3.generate({
				bindto: "#chart",
				donut: {
					label: { show: false },
					width: 11,
					title: "256MB Memory in use"
				},
				data: {
					type:"donut",
					columns: [
						["Available", 70],
						["Used", 39],
					],
					colors:{
						Used: "#3f9c35",
						Available: "#D1D1D1"
					},
				},
			// 	color: {
			// 		//pattern: ["#0088ce", "#EC7A08", "#EC7A08", "#cc0000"]
			// 		pattern: ["", "#0088ce", "#cc0000"], // the three color levels for the percentage values.
			// 		threshold: {
			// //            unit: 'value', // percentage is default
			// //            max: 200, // 100 is default
			// 			// values: [30, 60, 90, 100]
			// 		}
			// 	},

				legend: { show: false },
				size: { height: 171 }
			});
			var label = d3.select('text.c3-chart-arcs-title');
			label.html(''); // remove existant text
			label.insert('tspan').text('256').classed("graph-title-top", true).attr('dy', -10).attr('x', 0);
			label.insert('tspan').text('MB In Use').classed("graph-title", true).attr('dy', 20).attr('x', 0);
		});
	}
	getThreshold() {}
}
