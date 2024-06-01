import requests

BASE_URL = "http://localhost:3000"  # Replace with the actual base URL of your API
USERNAME = "testuser"
PASSWORD = "testpassword"


def register_user(username, password):
    url = f"{BASE_URL}/api/v1/auth/register"
    payload = {
        "username": username,
        "password": password
    }
    response = requests.post(url, json=payload)
    if response.status_code == 201:
        print("User registered successfully")
    else:
        print(f"User registration failed with status code: {response.status_code}")


def login_user(username, password):
    url = f"{BASE_URL}/api/v1/auth/login"
    payload = {
        "username": username,
        "password": password
    }
    response = requests.post(url, json=payload)
    if response.status_code == 200:
        token = response.json()["token"]
        print("User logged in successfully")
        return token
    else:
        print(f"User login failed with status code: {response.status_code}")
        return None


def send_protected_request(token):
    url = f"{BASE_URL}/api/v1/hello"  # Replace with the actual protected endpoint URL
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        print("Protected request successful")
        print(response.text)
    else:
        print(f"Protected request failed with status code: {response.status_code}")


if __name__ == '__main__':
    # Register a new user
    register_user(USERNAME, PASSWORD)

    # Login the user
    token = login_user(USERNAME, PASSWORD)

    # Send a protected request using the JWT token
    if token:
        send_protected_request(token)
