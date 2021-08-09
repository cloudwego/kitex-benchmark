namespace go echo

struct Request {
        1: required string action,
        2: required string msg,
}

struct Response {
        1: required string action,
        2: required string msg,
}

service EchoServer {
    Response Echo(1: Request req)
}
