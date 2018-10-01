import grpc
import service_pb2
import service_pb2_grpc

def main():
    channel = grpc.insecure_channel('localhost:6000')
    stub = service_pb2_grpc.UserServiceStub(channel)
    metadata = [('authorization', 'nuveosummit')]
    response = stub.CreateUser(service_pb2.CreateUserRequest(
        user = service_pb2.User(
            username = 'felipe',
            role = 'speaker'
        )
    ),
    metadata = metadata)
    if response:
        print('User created', response)
    response = stub.GetUser(service_pb2.GetUserRequest(
        username = 'felipe'
    ),
    metadata = metadata)
    if response:
        print('User', response)
    

if __name__ == '__main__':
  main()
