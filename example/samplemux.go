package main

import (
	"fmt"
	"github.com/jeyaprabhuj/ninetyoneweb/http/cookie"
	"github.com/jeyaprabhuj/ninetyoneweb/http/handler"
	//"github.com/jeyaprabhuj/ninetyoneweb/http/session"
	"github.com/jeyaprabhuj/ninetyoneweb/util/pipeline"
	"io/ioutil"
	"net/http"
)

var efsCookie *cookie.Cookie
var sessionStore *sessionMongo

var Chain1 pipeline.Chain
var Chain2 pipeline.Chain

func init() {

	realip := pipeline.Pipe{
		Flow:  logRealIP,
		Async: true,
	}

	requestinfo := pipeline.Pipe{
		Flow:  logRequestInfo,
		Async: true,
	}

	cookiePipe := pipeline.Pipe{
		Flow:  setCookie,
		Async: false,
	}

	Chain1 = pipeline.NewChain().BeforeHandlers(realip).AfterHandlers(requestinfo) // requestinfo is added in afterhandler for example
	Chain2 = pipeline.NewChain().BeforeHandlers(cookiePipe)

	efsCookie = cookie.NewCookie("9dP3rYPs0XZWNBHb0045bGAIckhayJV0")
	sessionStore = NewSessionMongo("ExternalSession")
}

func fileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w,
		`
		<html>
			<head></head>
			<body>
				<form enctype="multipart/form-data" action="/admin/filedisplay/" method="post">
					<label for="img">Select image:</label>
					<input type="file" name="updoc" />
					<input type="submit" value="upload"/>
				</form>
			</body>
		</html>
		`)
}

func filedisplay(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	filehandle := r.MultipartForm.File["updoc"][0]
	file, _ := filehandle.Open()
	b, _ := ioutil.ReadAll(file)
	fmt.Fprintln(w, string(b))
}

func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login with username password")
}

func paramsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.URL.Path)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.Header)
	fmt.Fprintln(w, r.Header["Accept"])
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	fmt.Fprintln(w, string(body))
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	//Session store mongo example is not complete

	// cookie, _ := r.Cookie("Cookie")
	// sessionManager := session.NewSessionManager(sessionStore)
	// newCookie := func() {
	// 	userSession := sessionManager.NewUserSession()
	// 	efsCookie.Set("Cookie", userSession.ID, w, r)
	// 	userSession.Store("Message", "Server side only acces to cookie")
	// }
	// if cookie == nil {
	// 	newCookie()
	// } else {
	// 	cookieValue, error := efsCookie.GetValue("Cookie", w, r)
	// 	if error != nil {
	// 		panic(error)
	// 	} else {
	// 		userSession, _ := sessionManager.GetUserSession(cookieValue)
	// 		fmt.Println("Cookie Value : " + cookieValue)
	// 		fmt.Println("userSession.Store['Message']:", userSession.Find("Message"))
	// 	}
	// }
}

func logRealIP(w http.ResponseWriter, r *http.Request) {
	fwd := r.Header.Get("X-FORWARDED-FOR")
	if fwd != "" {
		fmt.Println(fwd)
	} else {
		fmt.Println(r.RemoteAddr)
	}
}

func logRequestInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(
		"%s\t\t%s \n",
		r.Method,
		r.RequestURI,
	)
}

func checkAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Add logic to check if you are authenticated
		// Redirect to login page in case of web app
		// In case of api server send back a response code
		http.Redirect(w, r, "/login", 302)
	}
}

const port = ":8051"

func main() {
	mux := http.NewServeMux()
	mux2 := handler.NewMux("admin")
	mux.Handle(mux2.GetRootRoute(), mux2)
	mux2.GET("/router/", checkAuth(Chain1.FitFlow(paramsHandler).Process())) //different approach for example
	mux2.GET("/dynamic/", Chain2.FitFlow(getHandler).Process())              // Chain2 used here
	mux2.POST("/dynamic/", Chain1.FitFlow(postHandler).Process())            // Chain1 used here
	mux2.GET("/fileupload/", Chain1.FitFlow(fileUpload).Process())
	mux2.POST("/filedisplay/", Chain1.FitFlow(filedisplay).Process())

	mux2.GET("/login/", Chain1.FitFlow(adminLoginHandler).Process())
	//mux2.Print()
	server := http.Server{
		Addr:    port,
		Handler: mux2,
	}
	fmt.Println("Listening in " + port)
	server.ListenAndServe()

}
