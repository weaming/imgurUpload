package remote

// http://127.0.0.1:1024/#access_token=c50d48f10a7ef703ddac416205fab680110ab6f1&expires_in=315360000&token_type=bearer&refresh_token=061cfc38d6d4e2583521470a2259337eed0087a2&account_username=gardenyuen&account_id=20248338
var indexJS = `
// First, parse the query string
var params = {}, queryString = location.hash.substring(1),
regex = /([^&=]+)=([^&]*)/g, m;
while (m = regex.exec(queryString)) {
	params[decodeURIComponent(m[1])] = decodeURIComponent(m[2]);
}

// And send the token over to the server
var req = new XMLHttpRequest();
// consider using POST so query isn't logged
req.open('GET', 'http://' + window.location.host + '/catchtoken?' + queryString, true);

req.onreadystatechange = function (e) {
	if (req.readyState == 4) {
		if(req.status == 200){
			close()
			document.write(req.response)
		}
		else if(req.status == 400) {
			document.write('There was an error processing the token.')
		}
		else {
			document.write('something else other than 200 was returned')
		}
	}
};
req.send(null);
`
