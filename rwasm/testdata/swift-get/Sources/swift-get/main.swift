import Suborbital

class SwiftGet: Suborbital.Runnable {
    func run(input: String) -> String {
        return Suborbital.CacheGet(key: "important")
    }
}

@_cdecl("init")
func `init`() {
    Suborbital.Set(runnable: SwiftGet())
}