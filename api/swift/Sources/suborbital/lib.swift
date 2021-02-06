
@_silgen_name("return_result_swift")
func return_result(result_pointer: UnsafeRawPointer, result_size: Int32, ident: Int32)

@_silgen_name("log_msg_swift")
func log_msg(pointer: UnsafeRawPointer, size: Int32, level: Int32, ident: Int32)

@_silgen_name("fetch_url_swift")
func fetch_url(method: Int32, url_pointer: UnsafeRawPointer, url_size: Int32, body_pointer: UnsafeRawPointer, body_size: Int32, dest_pointer: UnsafeRawPointer, dest_max_size: Int32, ident: Int32) -> Int32

@_silgen_name("cache_set_swift")
func cache_set(key_pointer: UnsafeRawPointer, key_size: Int32, value_pointer: UnsafeRawPointer, value_size: Int32, ttl: Int32, ident: Int32) -> Int32
@_silgen_name("cache_get_swift")
func cache_get(key_pointer: UnsafeRawPointer, key_size: Int32, dest_pointer: UnsafeRawPointer, dest_max_size: Int32, ident: Int32) -> Int32

@_silgen_name("request_get_field_swift")
func request_get_field(field_type: Int32, key_pointer: UnsafeRawPointer, key_size: Int32, dest_pointer: UnsafeRawPointer, dest_max_size: Int32, ident: Int32) -> Int32

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
    var maxSize: Int32 = 256000
    var retVal = ""

    // loop until the returned size is within the defined max size, increasing it as needed
    var done = false
    while !done {
        toFFI(val: url, use: { (url_ptr: UnsafePointer<Int8>, url_size: Int32) in
            toFFI(val: body, use: { (body_ptr: UnsafePointer<Int8>, body_size: Int32) in
                let dest_ptr = allocate(size: Int32(maxSize))

                let resultSize = fetch_url(method: method, url_pointer: url_ptr, url_size: url_size, body_pointer: body_ptr, body_size: body_size, dest_pointer: dest_ptr, dest_max_size: maxSize, ident: CURRENT_IDENT)

                if resultSize == 0 {
                    done = true
                } else if resultSize < 0 {
                    retVal = "failed to fetch from url \(url)"
                    done = true
                } else if resultSize > maxSize {
                    maxSize = resultSize
                } else {
                    retVal = fromFFI(ptr: dest_ptr, size: resultSize)
                    done = true
                }
            })
        })
    }
    
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
    var maxSize: Int32 = 256000
    var retVal = ""

    // loop until the returned size is within the defined max size, increasing it as needed
    var done = false
    while !done {
        toFFI(val: key, use: { (keyPtr: UnsafePointer<Int8>, keySize: Int32) in
            let ptr = allocate(size: Int32(maxSize))

            let resultSize = cache_get(key_pointer: keyPtr, key_size: keySize, dest_pointer: ptr, dest_max_size: maxSize, ident: CURRENT_IDENT)

            if resultSize == 0 {
                done = true
            } else if resultSize < 0 {
                retVal = "failed to get from cache"
                done = true
            } else if resultSize > maxSize {
                maxSize = resultSize
            } else {
                retVal = fromFFI(ptr: ptr, size: resultSize)
                done = true
            }
        })
    }
    
    return retVal
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
    var maxSize: Int32 = 1024
    var retVal = ""

    // loop until the returned size is within the defined max size, increasing it as needed
    var done = false
    while !done {
        toFFI(val: key, use: { (keyPtr: UnsafePointer<Int8>, keySize: Int32) in
            let resultPtr = allocate(size: Int32(maxSize))

            let resultSize = request_get_field(field_type: fieldType, key_pointer: keyPtr, key_size: keySize, dest_pointer: resultPtr, dest_max_size: maxSize, ident: CURRENT_IDENT)
            
            if resultSize == 0 {
                done = true
            } else if resultSize < 0 {
                retVal = "failed to get request field"
                done = true
            } else if resultSize > maxSize {
                maxSize = resultSize
            } else {
                retVal = fromFFI(ptr: resultPtr, size: resultSize)
                done = true
            }
        })
    }
    
    return retVal
}

@_cdecl("run_e")
func run_e(pointer: UnsafeRawPointer, size: Int32, ident: Int32) {
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
func allocate(size: Int32) -> UnsafeMutableRawPointer {
  return UnsafeMutableRawPointer.allocate(byteCount: Int(size), alignment: MemoryLayout<UInt8>.alignment)
}

@_cdecl("deallocate")
func deallocate(pointer: UnsafeRawPointer, size: Int32) {
    let ptr: UnsafePointer<UInt8> = pointer.bindMemory(to: UInt8.self, capacity: Int(size))
    ptr.deallocate()
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
    let typed: UnsafePointer<UInt8> = ptr.bindMemory(to: UInt8.self, capacity: Int(size))
    let val = String(cString: typed)
    
    return val
}