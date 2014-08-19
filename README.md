flashlight-checker is a web app that allows checking connectivity to a
flashlight instance over the internet.

To use, simply make a GET request to http://flashlight-checker.herokuapp.com/,
for example:

```bash
curl -I http://flashlight-checker.herokuapp.com/ 
```

The app will attempt to proxy a HEAD request to http://www.google.com/humans.txt
via a flashlight running on the IP address from which the request was received,
and on port 443.

If successful, the request returns a 200 status code.  If unsuccessful, it
returns a 504.