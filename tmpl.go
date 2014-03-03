package main

var defaultSMSTemplate = []byte(`summary
== total: {{.Total}} fail: {{.FailCount}} ==
- host: {{.Host}}
- time: {{.TimeCost}}`)

var defaultTemplate = []byte(`<!DOCTYPE html>
<html>
	<head>
        <meta charset="utf-8" />
		<title>travelexec report</title>
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<!-- Bootstrap -->
		<link rel="stylesheet" href="http://cdn.bootcss.com/twitter-bootstrap/3.0.3/css/bootstrap.min.css">
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1 class="navbar" role="navigation">
					<span class="glyphicon glyphicon-hand-right"></span> 
					<span style="color:red">Travel</span><span style="color:blue">exec</span> report<small> created by: <a href="https://github.com/shxsun/travelexec">travelexec</a>
				</h1>
				<div class="row">
					<div class="col-md-6 summary">
						<table>
							<tr><td><b>Start time</b>:</td><td>{{.StartTime}}</td></tr>
							<tr><td><b>Time cost</b>:</td><td>{{.TimeCost}}</td></tr>
							<tr><td><b>Fail count</b>:</td><td>{{.FailCount}}</td></tr>
							<tr><td><b>Total</b>:</td><td>{{.Total}}</td></tr>
							<tr><td><b>Run host</b>:</td><td>{{.Host}}</td></tr>
						</table>
					</div>
					<div class="col-md-6">
						<button id="btn-show-failed" class="btn btn-default btn-danger">Only Show Failed</button>
					</div>
					<div class="clear"></div>
				</div>
			</div>
			<hr/>
			{{range .Tasks}}
			<div class="panel {{if .Error}} panel-danger {{else}} panel-success {{end}} case-list">
				<div class="panel-heading">
					<strong>{{.Command}}</strong> - {{.TimeCost}} <u><i>{{.Error}}</i></u>
					<div class="pull-right">
						<button class="btn btn-xs btn-info view-output">view output</button>
						<button class="btn btn-xs btn-primary view-source">view source(TODO)</button>
					</div>
					<div class="clear"></div>
				</div>
				<div class="panel-body">
					<dl>
						<dt>start time: {{.StartTime}}</dt>
						<dt>time cost: {{.TimeCost}}</dt>
						<dt>command: {{.Command}}</dt>
					</dl>
					<div class="source">
						<pre>{{.Source}}</pre>
					</div>
					<div class="output">
						<pre>{{.Output}}</pre>
					</div>
				</div>
			</div>
			{{end}}
			<div class="footer">
				&copy author: <a href="https://github.com/shxsun">skyblue</a>
			</div>
		</div>
		<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
		<script src="http://cdn.bootcss.com/jquery/1.10.2/jquery.min.js"></script>
		<!-- Include all compiled plugins (below), or include individual files as needed -->
		<script src="http://cdn.bootcss.com/twitter-bootstrap/3.0.3/js/bootstrap.min.js"></script>
		<script>
		$(function(){
			$(".source").hide();
			var nextOutput = function($btn){
				return $btn.parents("div.case-list").find("div.output");
			};
			$("button.view-output").click(function(){
				$output = nextOutput($(this));
				$output.toggle();
			}).each(function(){
				nextOutput($(this)).hide();
			});
			$("div.panel-body").hide();
			$("div.panel-heading").click(function(){
				$(this).next("div.panel-body").toggle();
			});
			$("#btn-show-failed").click(function(){
				$("div.panel-success").toggle();
			});
		});
		</script>
	</body>
</html>
`)
