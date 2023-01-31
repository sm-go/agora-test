## agora
## webRTC

install adb and DroidCam
https://www.dev47apps.com/droidcam/linux/
Setting > Build number > Developer options > USB debuging

go get github.com/gin-gonic/gin
go get github.com/AgoraIO-Community/go-tokenbuilder

agora မှာ အကောင့် ဆောက်မယ်၊ 
project တစ်ခု ဖန်တီးမယ်
APP_ID ထွက်လာမယ်။
channel name တစ်ခု ဖန်တီးမယ်။
Token တစ်ခု ထွက်လာမယ်။
(https://webdemo.agora.io/basicVideoCall/index.html#)

$ go server project ကို run မယ်။
$ api ခေါ်လိုက်တဲ့ အခါ token ထွက်လာမယ်
(http://localhost:8080/rtc/testing/publisher/userAccount/123456/)
ရလာတဲ့ token ကို node project မှာ ထည့် run မယ််။
join လိုက်ရင် ရပြီ