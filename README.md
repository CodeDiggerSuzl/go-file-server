refs:
- https://github.com/blankjee/file-storage-system
- https://github.com/samtake/filestore-server/blob/master/filestore-server-v0.2/static/view/signin.html


# Arch
- Front end : React + Ant D

- Back end: Golang + Gin

## Coding Progress
###  Simple file upload:
   1. [Http form action url `#`](https://developer.mozilla.org/zh-CN/docs/Learn/HTML/Forms/Sending_and_retrieving_form_data)
   2. [Http code 302](https://www.cnblogs.com/woshimrf/p/http-code-302.html)
   3. [HandleFunc `func http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))`]()
   4. `os.Create()` when the dir not exist will error.


### Keep meta data of the file
1. Use the hash of the file as it's id.
2. Generate file meta info before save files.
3. FileMeta info query interface.
   1. Single file query.
   2. Range file query.
4. Download file using fileHash.

   A better way is: when the file is uploaded, and handle these files as static source, NGINX, for example, deploy it as _Reverse Proxy_ to download the file, also we can do simple stream limiting and permission check, which will reduce the pressure of the back end.
5. Update file name by file hash.
6. Del file both in the fileMetaMap and disk.

Current work flow:
1. Using get to get the file upload page.
2. Choose the local file and, use form to upload file.
3. Back end get file upload stream, and save to local disk.
4. Update the the `metaFileMap`.

# What I need ?

1. RateLimit from stream video
2. Watch src code of golang
3. QR-Code info