namespace go echo

struct Request {
        1: required string message,
}

struct Response {
        1: required string message,
}

service EchoServer {
    Response Echo(1: Request req)
}
