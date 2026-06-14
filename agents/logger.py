import logging
import os

def setup_logger():
    logger = logging.getLogger("AI_Agent")
    logger.setLevel(logging.INFO)
    
    # Avoid duplicate handlers if already initialized
    if logger.handlers:
        return logger

    # Formatter with timestamp
    formatter = logging.Formatter('[%(asctime)s] [%(levelname)s] - %(message)s')

    # Console Handler
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)

    # File Handler
    log_dir = os.path.dirname(os.path.abspath(__file__))
    log_file = os.path.join(log_dir, "agent.log")
    file_handler = logging.FileHandler(log_file)
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)

    return logger

logger = setup_logger()
