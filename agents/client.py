import sys
import os
import grpc

# Adjust sys.path to find sibling imports
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

import performance_pb2
import performance_pb2_grpc
from config import GRPC_SERVER_ADDR
from logger import logger

class GrpcClient:
    def __init__(self, address=GRPC_SERVER_ADDR):
        self.address = address
        self.channel = grpc.insecure_channel(address)
        self.stub = performance_pb2_grpc.PerformanceServiceStub(self.channel)

    def fetch_user_information(self, email_or_id):
        logger.info(f"gRPC Outgoing Request --> GetUserInformation for ID/Email: {email_or_id}")
        request = performance_pb2.UserRequest(email_or_id=email_or_id)
        try:
            response = self.stub.GetUserInformation(request)
            if response.error:
                logger.error(f"gRPC Incoming Response <-- GetUserInformation Error: {response.error}")
                raise Exception(response.error)
            
            logger.info(f"gRPC Incoming Response <-- GetUserInformation SUCCESS")
            return response.json_result
        except Exception as e:
            logger.error(f"gRPC Incoming Response <-- GetUserInformation FAILED: {e}")
            raise e
