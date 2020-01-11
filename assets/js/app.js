Vue.component('measurement', {
	template: '#measurement',
	props: {
		data: Object,
		name: String,
		val: String,
		sum: Boolean,
	},
	data: function () {
		self=this;
		let p = function(i) {
			return self.data[i] !== undefined && self.data[i] !== null && self.data[i] !== "";
		}
		let v = function(i) {
			return self.data[i];
		}
		console.log(this.val)
		console.log(this.data)
		return {
			// display: p(this.val) || p(this.val+"L1")||p(this.val+"L2")||p(this.val+"L3") || p(this.val+"S1")||p(this.val+"S2")||p(this.val+"S3"),
			display: true,
			// l1: p(this.val+"L1") || p(this.val+"S1"),
			// l2: p(this.val+"L2") || p(this.val+"S1"),
			// l3: p(this.val+"L3") || p(this.val+"S1"),
			l1: true,
			l2: true,
			l3: true,
			val1: v(this.val + "L1") || v(this.val + "S1"),
			val2: v(this.val + "L2") || v(this.val + "S2"),
			val3: v(this.val + "L3") || v(this.val + "S3"),
			valsum:
				v(this.val) ||
				v(this.val + "L1") + v(this.val + "L2") + v(this.val + "L3") ||
				v(this.val + "S1") + v(this.val + "S2") + v(this.val + "S3")
		};
	},
});


var dataapp = new Vue({
	el: '#realtime',
	delimiters: ['${', '}'],
	data: {
		meters: {},
		message: 'Loading...'
	},
	computed: {
		// return meters sorted by name
		sortedMeters: function() {
			var devs = Object.keys(this.meters);
			devs.sort();
			var res = {};
			devs.forEach(function(key) {
				res[key] = this.meters[key];
			}, this);
			return res;
		}
	},
	methods: {
		// pop returns true if it was called with any non-null argumnt
		pop: function () {
			for(var i=0; i<arguments.length; i++) {
				if (arguments[i] !== undefined && arguments[i] !== null && arguments[i] !== "") {
					return true;
				}
			}
			return false;
		},

		// val returns addable value: null, NaN and empty are converted to 0
		val: function (v) {
			v = parseFloat(v);
			return isNaN(v) ? 0 : v;
		}
	}
})

var timeapp = new Vue({
	el: '#time',
	delimiters: ['${', '}'],
	data: {
		time: 'n/a',
		date: 'n/a'
	}
})

var statusapp = new Vue({
	el: '#status',
	delimiters: ['${', '}'],
	data: {
		meters: {}
	}
})

var fixed = d3.format(".2f")
var si = d3.format(".3~s")

$().ready(function () {
	connectSocket();
});

function convertDate(unixtimestamp){
	var date = new Date(unixtimestamp);
	var day = "0" + date.getDate();
	var month = "0" + (date.getMonth() + 1);
	var year = date.getFullYear();
	return year + '/' + month.substr(-2) + '/' + day.substr(-2);
}

function convertTime(unixtimestamp){
	var date = new Date(unixtimestamp);
	var hours = date.getHours();
	var minutes = "0" + date.getMinutes();
	var seconds = "0" + date.getSeconds();
	return hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);
}

function updateTime(data) {
	timeapp.date = convertDate(data["Timestamp"])
	timeapp.time = convertTime(data["Timestamp"])
}

function updateStatus(status) {
	var id = status["Device"]
	status["Status"] = status["Online"] ? "online" : "offline"

	// update data table
	var dict = statusapp.meters[id] || {}
	dict = Object.assign(dict, status)

	// make update reactive, see
	// https://vuejs.org/v2/guide/reactivity.html#Change-Detection-Caveats
	Vue.set(statusapp.meters, id, dict)
}

function updateData(data) {
	// extract the last update
	var id = data["Device"]
	var type = data["IEC61850"]
	var value = fixed(data["Value"])

	// create or update data table
	var dict = dataapp.meters[id] || {}
	dict[type] = value

	// put into statusline
	dataapp.message = "Received " + id + " / " + type + ": " + si(value)

	// make update reactive, see
	// https://vuejs.org/v2/guide/reactivity.html#Change-Detection-Caveats
	Vue.set(dataapp.meters, id, dict)
}

function processMessage(data) {
	if (data.Meters && data.Meters.length) {
		for (var i=0; i<data.Meters.length; i++) {
			updateStatus(data.Meters[i]);
		}
	}
	else if (data.Device) {
		updateTime(data);
		updateData(data);
	}
}

function connectSocket() {
	var ws, loc = window.location;
	var protocol = loc.protocol == "https:" ? "wss:" : "ws:"

	// ws = new WebSocket(protocol + "//" + loc.hostname + (loc.port ? ":" + loc.port : "") + "/ws");
	ws = new WebSocket("ws://localhost:8081/ws");

	ws.onerror = function(evt) {
		// console.warn("Connection error");
		ws.close();
	}
	ws.onclose = function (evt) {
		// console.warn("Connection closed");
		window.setTimeout(connectSocket, 1000);
	};
	ws.onmessage = function (evt) {
		var json = JSON.parse(evt.data);
		processMessage(json);
	};
}
