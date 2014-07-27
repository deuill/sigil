package http

type errorCode struct {
	Code  string
	Desc  string
	Title string
	Body  string
}

var errorCodes = map[int]errorCode{
	100: errorCode{"100", "Continue", "", ""},
	101: errorCode{"101", "Switching Protocols", "", ""},
	200: errorCode{"200", "OK", "", ""},
	201: errorCode{"201", "Created", "", ""},
	202: errorCode{"202", "Accepted", "", ""},
	203: errorCode{"203", "Non-Authoritative Information", "", ""},
	204: errorCode{"204", "No Content", "", ""},
	205: errorCode{"205", "Reset Content", "", ""},
	206: errorCode{"206", "Partial Content", "", ""},
	300: errorCode{"300", "Multiple Choices", "", ""},
	301: errorCode{"301", "Moved Permanently", "", ""},
	302: errorCode{"302", "Moved Temporarily", "", ""},
	303: errorCode{"303", "See Other", "", ""},
	304: errorCode{"304", "Not Modified", "", ""},
	305: errorCode{"305", "Use Proxy", "", ""},
	400: errorCode{"400", "Bad Request", "", ""},
	401: errorCode{"401", "Unauthorized", "", ""},
	402: errorCode{"402", "Payment Required", "", ""},
	403: errorCode{"403", "Forbidden", "", ""},
	404: errorCode{"404", "Not Found", "", "Sorry, but the page you were trying to view does not exist"},
	405: errorCode{"405", "Method Not Allowed", "", ""},
	406: errorCode{"406", "Not Acceptable", "", ""},
	407: errorCode{"407", "Proxy Authentication Required", "", ""},
	408: errorCode{"408", "Request Time-out", "", ""},
	409: errorCode{"409", "Conflict", "", ""},
	410: errorCode{"410", "Gone", "", ""},
	411: errorCode{"411", "Length Required", "", ""},
	412: errorCode{"412", "Precondition Failed", "", ""},
	413: errorCode{"413", "Request Entity Too Large", "", ""},
	414: errorCode{"414", "Request-URI Too Large", "", ""},
	415: errorCode{"415", "Unsupported Media Type", "", ""},
	500: errorCode{"500", "Internal Server Error", "An error has occurred", "The server failed to complete the request"},
	501: errorCode{"501", "Not Implemented", "", ""},
	502: errorCode{"502", "Bad Gateway", "", ""},
	503: errorCode{"503", "Service Unavailable", "", ""},
	504: errorCode{"504", "Gateway Timeout", "", ""},
	505: errorCode{"505", "HTTP Version not supported", "", ""},
}

var errorTemplate = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>{{or .Title .Desc}} :(</title>
		<style>
			::-moz-selection {
				background: #b3d4fc;
				text-shadow: none;
			}

			::selection {
				background: #b3d4fc;
				text-shadow: none;
			}

			html {
				height: 100%;
				font-size: 20px;
				line-height: 1.4;
				color: #737373;
				background: #f0f0f0 url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAMAAAAp4XiDAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAEtQTFRFjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mjo+Mju7hoQAAABl0Uk5TAAECAwQFBgcICQoLDA0OEBESExQVFhcZGtS5Hb4AAAFXSURBVEjH1dTBboQgEIDhfwZQEEVd123f/0l7cN2wVUzbg7GTzOm/AAkf6RHZnWLg4/MBQwfdAOEmuLuidyfFQDtHqB24GowH9YIE3QlhCWfNUEOKEBPUA9hR0dEig8+DvAK0DhoPvgHXgukE6cxBOGm6CqIHH6HqwCRBk3mG8B4kGaB/Xr9JUPdgB0UHi5TCWRMs+AoqDzaANoI0ehC4eehbaHvwo+AmRSeH3Dz03TacNU7BGjAO1IFUglRyEJhC9pOm5SeZuUKmAEPKwqyY2QEqIOsq8NpS+MsYAdVlxSyLkYPAPcCYoBvzI1fIfb3LtwB2fRgLakEc4OQgIKUTF0Oa47PK+xbDQl+fCzdl9PU/oy9ckb71i48WeQ+79MUSffE69MmLPn9EX7wifXUuXMjo24RwTF+4CH22RJ/dpS8d0bfhQjPd9unTf0Pf5mFW4Yrh9/R9AWqxGIl05IuAAAAAAElFTkSuQmCC) repeat 0 0;
				-webkit-text-size-adjust: 100%;
				-ms-text-size-adjust: 100%;
			}

			@media (min-width: 768px) and (max-width: 979px) {
				html {
					font-size: 18px;
				}
			}

			@media (max-width: 767px) {
				html {
					font-size: 16px;
				}
			}

			body {
				margin: 0;
				padding: 0;
				height: 100%;
			}

			html,
			input {
				font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
			}

			h1 {
				font-size: 3em;
				margin: 0 0 10px 0;
				position: relative;
			}

			h1 span {
				color: #bbb;
			}

			h3 {
				margin: 1.5em 0 0.5em;
			}

			p {
				margin: 1em 0;
			}

			ul {
				padding: 0 0 0 40px;
				margin: 1em 0;
			}

			a {
				color: #444;
			}

			.container {
				margin-top: -150px;
				position: relative;
				text-align: center;
				top: 50%;
			}

			.response {
				color: #888;
				font-size: 0.9em;
			}

			.button {
				display: inline-block;
				width: 120px;
			}

			.button a {
				color: #737373;
				outline: none;
				text-decoration: none;
			}

			.button .icon {
				cursor: pointer;
				display: block;
				font-size: 3em;
				font-style: normal;
				margin: 0 auto -15px;
				-webkit-transition: -webkit-transform .2s ease;
				-moz-transition: -moz-transform .2s ease;
				-ms-transition: -ms-transform .2s ease;
				-o-transition: -o-transform .2s ease;
				transition: transform .2s ease;
			}

			.button .descr {
				font-size: 0.8em;
				opacity: 0;
				position: relative;
				bottom: 10px;
				white-space: nowrap;
				-webkit-transition: opacity .3s ease, bottom .3s ease;
				-moz-transition: opacity .3s ease, bottom .3s ease;
				-ms-transition: opacity .3s ease, bottom .3s ease;
				-o-transition: opacity .3s ease, bottom .3s ease;
				transition: opacity .3s ease, bottom .3s ease;
			}

			.button .icon:hover + .descr {
				opacity: 1;
				bottom: 0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>{{or .Title .Desc}} <span>:(</span></h1>
			<p class="response">({{.Code}} {{or .Title .Desc}})</p>
			{{if .Body}}<p>{{.Body}}</p>{{end}}
			
			<span class="button back">
				<a href="/" class="icon">↶</a>
				<span class="descr">Home</span>
			</span>

			<span class="button back">
				<i class="icon" onclick="history.back(-1)">↖</i>
				<span class="descr">Go Back</span>
			</span>

			<span class="button refresh">
				<i class="icon" onclick="document.location.reload(true)">↺</i>
				<span class="descr">Try Again</span>
			</span>
		</div>
	</body>
</html>`
