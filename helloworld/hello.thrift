struct Result {
  1: string response,
}

service HelloService {
  Result HelloWorld(1:string name);
}
