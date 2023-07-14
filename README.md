# Chrome DevTools Network TUI (cnTUI)

Little cli/tui tool to export HTTP requests to cURL commands from the terminal.  
Currently exports requests to your X11 clipboard using xclip. 

![image](https://github.com/fipso/cntui/assets/8930842/074541a2-10a5-426a-aaad-b6051b81e5f8)

  
### Installation
Requirements: go, xclip  
  
`git clone https://github.com/fipso/cntui.git`  
`go build .`  
`sudo cp /usr/local/bin cntui`  
  
### Usage
- Start chrome with open DevTools Server (Debug mode):  
  `google-chrome-stable --remote-debugging-port=9222`  
- Select tab you want to hack on  
- Run `cntui`  
- Select a request. Hit enter
- cURL command is now in your clipboard. Paste. Have fun

### TODO:
- [ ] Edit Mode: edit post body in terminal editor
- [ ] Replay request directly
- [ ] Request description screen

The development of this tool has been recorded on YouTube:  
https://www.youtube.com/watch?v=ywqy_tIq7xc&list=PLd-Mx7H0BuG9Cfsu-8oqUS54MxbaOxmg9
