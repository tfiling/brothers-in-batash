import csv
import json

import requests

BASE_URL = "http://localhost:3000"  # Replace with the actual base URL of your API
USERNAME = "john_doe"
PASSWORD = "admin"


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


def read_csv_data(file_path):
    with open(file_path, 'r') as file:
        reader = csv.DictReader(file)
        data = []
        for row in reader:
            converted_row = {}
            for key, value in row.items():
                field_name, field_type = extract_field_name_and_type(key)
                converted_row[field_name] = cast_value(value, field_type)
            data.append(converted_row)
    return data


def extract_field_name_and_type(field_header):
    if '(' in field_header and field_header.endswith(')'):
        field_name, field_type = field_header[:-1].split('(')
        return field_name, field_type
    else:
        return field_header, 'string'


def cast_value(value, field_type):
    if field_type == 'int':
        return int(value)
    elif field_type == 'float':
        return float(value)
    elif field_type == 'bool':
        return bool(value)
    elif field_type == 'dict':
        return json.loads(value)
    else:
        return value


def create_soldiers(token, soldiers_data):
    url = f"{BASE_URL}/api/v1/soldiers"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    for soldier in soldiers_data:
        response = requests.post(url, json=soldier, headers=headers)
        if response.status_code == 201:
            print(f"Soldier created successfully: {soldier}")
        else:
            print(f"Soldier creation failed with status code: {response.status_code}")


def create_shifts(token, shifts_data):
    url = f"{BASE_URL}/api/v1/shifts"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    for shift in shifts_data:
        response = requests.post(url, json=shift, headers=headers)
        if response.status_code == 201:
            print(f"Shift created successfully: {shift}")
        else:
            print(f"Shift creation failed with status code: {response.status_code}")


def create_day_schedules(token, day_schedules_data):
    url = f"{BASE_URL}/api/v1/day-schedules"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    for day_schedule in day_schedules_data:
        response = requests.post(url, json=day_schedule, headers=headers)
        if response.status_code == 201:
            print(f"Day Schedule created successfully: {day_schedule}")
        else:
            print(f"Day Schedule creation failed with status code: {response.status_code}")


def create_shift_templates(token, shift_templates_data):
    url = f"{BASE_URL}/api/v1/shift-templates"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    for shift_template in shift_templates_data:
        response = requests.post(url, json=shift_template, headers=headers)
        if response.status_code == 201:
            print(f"Shift Template created successfully: {shift_template}")
        else:
            print(f"Shift Template creation failed with status code: {response.status_code}")


def get_soldiers(token):
    url = f"{BASE_URL}/api/v1/soldiers"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        soldiers = response.json()
        print("Soldiers retrieved successfully:")
        for soldier in soldiers:
            print(soldier)
    else:
        print(f"Failed to retrieve soldiers with status code: {response.status_code}")


def get_shifts(token):
    url = f"{BASE_URL}/api/v1/shifts"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        shifts = response.json()
        print("Shifts retrieved successfully:")
        for shift in shifts:
            print(shift)
    else:
        print(f"Failed to retrieve shifts with status code: {response.status_code}")


def get_day_schedules(token):
    url = f"{BASE_URL}/api/v1/day-schedules"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        day_schedules = response.json()
        print("Day Schedules retrieved successfully:")
        for day_schedule in day_schedules:
            print(day_schedule)
    else:
        print(f"Failed to retrieve day schedules with status code: {response.status_code}")


def get_shift_templates(token):
    url = f"{BASE_URL}/api/v1/shift-templates"
    headers = {
        "Authorization": f"Bearer {token}"
    }
    response = requests.get(url, headers=headers)
    if response.status_code == 200:
        shift_templates = response.json()
        print("Shift Templates retrieved successfully:")
        for shift_template in shift_templates:
            print(shift_template)
    else:
        print(f"Failed to retrieve shift templates with status code: {response.status_code}")


if __name__ == '__main__':
    # # Register a new user
    # register_user(USERNAME, PASSWORD)

    # Login the user
    token = login_user(USERNAME, PASSWORD)

    # Read fake data from CSV files
    soldiers_data = read_csv_data("fake_data/soldier.csv")
    shifts_data = read_csv_data("fake_data/shifts.csv")
    day_schedules_data = read_csv_data("fake_data/day_schedules.csv")
    shift_templates_data = read_csv_data("fake_data/shift_templates.csv")

    # Create soldiers
    create_soldiers(token, soldiers_data)

    # Create shifts
    create_shifts(token, shifts_data)

    # Create day schedules
    create_day_schedules(token, day_schedules_data)

    # Create shift templates
    create_shift_templates(token, shift_templates_data)

    get_soldiers(token)
    get_shifts(token)
    get_day_schedules(token)
    get_shift_templates(token)
