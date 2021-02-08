import Suborbital

class GetStaticSwift: Suborbital.Runnable {
    func run(input: String) -> String {
        return Suborbital.GetStaticFile(name: "important.md")
    }
}

@_cdecl("init")
func `init`() {
    Suborbital.Set(runnable: GetStaticSwift())
}