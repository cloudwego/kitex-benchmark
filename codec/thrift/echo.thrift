namespace go echo

struct Request {
        1: required string action,
        2: required string msg,
}

struct Response {
        1: required string action,
        2: required string msg,
}

struct SimpleStruct {
    1: bool Bool
    2: i8 Int8
    3: i16 Int16
    4: i32 Int32
    5: i64 Int64
    6: double Double
    7: string String
    8: binary Binary
    9: list<i64> Int64List
    10: map<string,string> StringStringMap
}

struct NestedStruct {
    1: SimpleStruct Struct
    2: list<SimpleStruct> StructList
    3: map<string,SimpleStruct> StructMap
}

struct NestedRequest {
    1: NestedStruct NestedStruct
    2: Request Request
}

struct NestedResponse {
    1: NestedStruct NestedStruct
    2: Response Response
}

service EchoServer {
    Response Echo(1: Request req)
    NestedResponse NestedEcho(1: NestedRequest req)
}
