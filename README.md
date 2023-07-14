# Chrome DevTools Network TUI (cnTUI)

Little cli/tui tool to export HTTP requests to cURL commands from the terminal. 
Currently exports requests to your X11 clipboard using xclip. 
  
### Installtion
Requirements: go, xlcip  
  
`git clone https://github.com/fipso/cntui.git`
`go build .`
`sudo cp /usr/local/bin cntui`

### Usage
- Start chrome with open DevTools Server (Debug mode): `google-chrome-stable --remote-debugging-port=9222`  
- Select tab you want to hack on  
- Run `cntui`  
- Select a request. Hit enter
- cURL command is now in you clipboard. Paste. Have fun

### TODO:
- [ ] Edit Mode: edit post body in terminal editor
- [ ] Replay request directly
- [ ] Request description screen

The development of this tool has been recorded on youtube:
https://www.youtube.com/watch?v=ywqy_tIq7xc&list=PLd-Mx7H0BuG9Cfsu-8oqUS54MxbaOxmg9
