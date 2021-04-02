
@_silgen_name("return_result_swift")
func return_result(result_pointer: UnsafeRawPointer, result_size: Int32, ident: Int32)

@_silgen_name("get_ffi_result_swift")
func get_ffi_result(result_pointer: UnsafeRawPointer, ident: Int32) -> Int32

@_silgen_name("log_msg_swift")
func log_msg(pointer: UnsafeRawPointer, size: Int32, level: Int32, ident: Int32)

@_silgen_name("fetch_url_swift")
func fetch_url(method: Int32, url_pointer: UnsafeRawPointer, url_size: Int32, body_pointer: UnsafeRawPointer, body_size: Int32, ident: Int32) -> Int32

@_silgen_name("cache_set_swift")
func cache_set(key_pointer: UnsafeRawPointer, key_size: Int32, value_pointer: UnsafeRawPointer, value_size: Int32, ttl: Int32, ident: Int32) -> Int32
@_silgen_name("cache_get_swift")
func cache_get(key_pointer: UnsafeRawPointer, key_size: Int32, ident: Int32) -> Int32

@_silgen_name("request_get_field_swift")
func request_get_field(field_type: Int32, key_pointer: UnsafeRawPointer, key_size: Int32, ident: Int32) -> Int32

@_silgen_name("get_static_file_swift")
func get_static_file(name_pointer: UnsafeRawPointer, name_size: Int32, ident: Int32) -> Int32

// keep track of the current ident
var CURRENT_IDENT: Int32 = 0

// the Runnable instance currently being used
var RUNNABLE: Runnable = defaultRunnable()

// the protocol that users conform to to make their package a Runnable
public protocol Runnable {
    func run(input: String) -> String
}

// something to hold the Runnable's place until set is called
class defaultRunnable: Runnable {
    func run(input: String) -> String {
        return ""
    }
}

public func Set(runnable: Runnable) {
    RUNNABLE = runnable
}

let httpMethodGet = Int32(1)
let httpMethodPost = Int32(2)
let httpMethodPatch = Int32(3)
let httpMethodDelete = Int32(4)

public func HttpGet(url: String) -> String {
    return fetch(method: httpMethodGet, url: url, body: "")
}

public func HttpPost(url: String, body: String) -> String {
    return fetch(method: httpMethodPost, url: url, body: body)
}

public func HttpPatch(url: String, body: String) -> String {
    return fetch(method: httpMethodPatch, url: url, body: body)
}

public func HttpDelete(url: String) -> String {
    return fetch(method: httpMethodDelete, url: url, body: "")
}

func fetch(method: Int32, url: String, body: String) -> String {
    var retVal = ""

    toFFI(val: url, use: { (url_ptr: UnsafePointer<Int8>, url_size: Int32) in
        toFFI(val: body, use: { (body_ptr: UnsafePointer<Int8>, body_size: Int32) in
            let resultSize = fetch_url(method: method, url_pointer: url_ptr, url_size: url_size, body_pointer: body_ptr, body_size: body_size, ident: CURRENT_IDENT)

            retVal = ffiResult(size: resultSize)
        })
    })
    
    return retVal
}

public func CacheSet(key: String, value: String, ttl: Int) {    

    toFFI(val: key, use: { (keyPtr: UnsafePointer<Int8>, keySize: Int32) in
        toFFI(val: value, use: { (valPtr: UnsafePointer<Int8>, valSize: Int32) in
            let _ = cache_set(key_pointer: keyPtr, key_size: keySize, value_pointer: valPtr, value_size: valSize, ttl: Int32(ttl), ident: CURRENT_IDENT)
        })
    })
}

public func CacheGet(key: String) -> String {
    var retVal = ""

    toFFI(val: key, use: { (keyPtr: UnsafePointer<Int8>, keySize: Int32) in
        let resultSize = cache_get(key_pointer: keyPtr, key_size: keySize, ident: CURRENT_IDENT)

        retVal = ffiResult(size: resultSize)
    })
    
    return retVal
}

public func LogDebug(msg: String) {
    log(msg: msg, level: 4)
}

public func LogInfo(msg: String) {
    log(msg: msg, level: 3)
}

public func LogWarn(msg: String) {
    log(msg: msg, level: 2)
}

public func LogErr(msg: String) {
    log(msg: msg, level: 1)
}

func log(msg: String, level: Int32) {
    toFFI(val: msg, use: { (ptr: UnsafePointer<Int8>, size: Int32) in
        log_msg(pointer: ptr, size: size, level: level, ident: CURRENT_IDENT)
    })
}

let fieldTypeMeta = Int32(0)
let fieldTypeBody = Int32(1)
let fieldTypeHeader = Int32(2)
let fieldTypeParams = Int32(3)
let fieldTypeState = Int32(4)

public func ReqMethod() -> String {
    return requestGetField(fieldType: fieldTypeMeta, key: "method")
}

public func ReqURL() -> String {
    return requestGetField(fieldType: fieldTypeMeta, key: "url")
}

public func ReqID() -> String {
    return requestGetField(fieldType: fieldTypeMeta, key: "id")
}

public func ReqBodyRaw() -> String {
    return requestGetField(fieldType: fieldTypeMeta, key: "body")
}

public func ReqBodyField(key: String) -> String {
    return requestGetField(fieldType: fieldTypeBody, key: key)
}

public func ReqHeader(key: String) -> String {
    return requestGetField(fieldType: fieldTypeHeader, key: key)
}

public func ReqParam(key: String) -> String {
    return requestGetField(fieldType: fieldTypeParams, key: key)
}

public func State(key: String) -> String {
    return requestGetField(fieldType: fieldTypeState, key: key)
}

func requestGetField(fieldType: Int32, key: String) -> String {
    var retVal = ""

    toFFI(val: key, use: { (keyPtr: UnsafePointer<Int8>, keySize: Int32) in
        let resultSize = request_get_field(field_type: fieldType, key_pointer: keyPtr, key_size: keySize, ident: CURRENT_IDENT)
        
        retVal = ffiResult(size: resultSize)
    })
    
    return retVal
}

public func GetStaticFile(name: String) -> String {
    var retVal = ""

    toFFI(val: name, use: { (namePtr: UnsafePointer<Int8>, nameSize: Int32) in
        let resultSize = get_static_file(name_pointer: namePtr, name_size: nameSize, ident: CURRENT_IDENT)

        retVal = ffiResult(size: resultSize)
    })
    
    return retVal
}

@_cdecl("run_e")
func run_e(pointer: UnsafeMutablePointer<Int8>, size: Int32, ident: Int32) {
    CURRENT_IDENT = ident
    
    let inString = fromFFI(ptr: pointer, size: size)
    
    // call the user-provided run function
    let retString = RUNNABLE.run(input: inString)

    // convert the output to a usable pointer/size combo
    toFFI(val: retString, use: { (ptr: UnsafePointer<Int8>, size: Int32) in
        return_result(result_pointer: ptr, result_size: size, ident: ident)
    })
}

@_cdecl("allocate")
func allocate(size: Int32) -> UnsafeMutablePointer<Int8> {
  return UnsafeMutablePointer<Int8>.allocate(capacity: Int(size) + 1)
}

@_cdecl("deallocate")
func deallocate(ptr: UnsafeRawPointer, size: Int32) {
    let ptr: UnsafeMutablePointer<Int8> = UnsafeMutablePointer(mutating: ptr.bindMemory(to: Int8.self, capacity: Int(size) + 1))
    ptr.deinitialize(count: Int(size) + 1)
    ptr.deallocate()
}

func ffiResult(size: Int32) -> String {
    if size < 0 {
        LogErr(msg: "an error was returned")
        return ""
    }
    
    let resultPtr = allocate(size: size)
    
    let code = get_ffi_result(result_pointer: resultPtr, ident: CURRENT_IDENT)
    
    if code != Int32(0) {
        LogErr(msg: "an error was returned")
        return ""
    }
    
    return fromFFI(ptr: resultPtr, size: size)
}

func toFFI(val: String, use: (UnsafePointer<Int8>, Int32) -> Void) {
    // create a nil (optional) pointer
    let size = Int32(val.utf8.count)

    // grab the pointer in a closure and give the optional a real value
    let _ = val.withCString({ (valPtr) -> UInt in
        use(valPtr, size)
        return 0
    })
}

func fromFFI(ptr: UnsafeRawPointer, size: Int32) -> String {
    let typed: UnsafeMutablePointer<Int8> = UnsafeMutablePointer(mutating: ptr.bindMemory(to: Int8.self, capacity: Int(size) + 1))
    let term = typed + Int(size)
    term.pointee = 0
    
    let val = String(cString: typed)
    
    return val
}
