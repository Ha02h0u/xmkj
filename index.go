package main

import "net/http"

const indexWeb = `<html>
<head>
<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
 <input type="file" name="file" />
 <input type="submit" value="upload" />
</form>
</body>
</html>`

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexWeb))
}
