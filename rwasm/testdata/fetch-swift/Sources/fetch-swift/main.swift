import Suborbital

class FetchSwift: Suborbital.Runnable {
    func run(input: String) -> String {

        let url = "https://postman-echo.com/post"
        let body = "hello, postman!"
        let resp = Suborbital.HttpPost(url: url, body: body)

        Suborbital.LogInfo(msg: resp)

        return Suborbital.HttpGet(url: input)
    }
}

@_cdecl("init")
func `init`() {
    Suborbital.Set(runnable: FetchSwift())
}