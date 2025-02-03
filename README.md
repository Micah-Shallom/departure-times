# Departure Times

A Go service that geo-localizes a user and provides real-time departure times for public transportation.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)

## Features

- Geo-location of users.
- Real-time departure times for various public transportation.
- Supports multiple transportation methods (buses, trains, etc.).
- Easy to deploy and configure.

## Installation

1. Clone the repository:

```sh
git clone https://github.com/Micah-Shallom/departure-times.git
cd departure-times
```

2. Install dependencies:

```sh
go mod tidy
```

## Usage

1. Start the server:

```sh
go run main.go
```

2. Access the service at `http://localhost:8080`.

### Example Request

To get the departure times, you can make a GET request to the `/departures` endpoint:

```sh
curl -X GET "http://localhost:8080/departures?latitude=40.7128&longitude=-74.0060"
```

## API Endpoints

### `GET /departures`

Retrieve real-time departure times for public transportation based on user location.

#### Query Parameters

- `latitude` (required): The latitude of the user's location.
- `longitude` (required): The longitude of the user's location.

#### Response

```json
{
  "departure_times": [
    {
      "line": "Bus 42",
      "departure_in": "5 minutes"
    },
    {
      "line": "Train A",
      "departure_in": "10 minutes"
    }
  ]
}
```

## Configuration

The service can be configured using environment variables:

- `PORT`: The port on which the server runs (default: 8080).
- `API_KEY`: The API key for accessing the transportation data provider.
- `LOG_LEVEL`: The log level (e.g., DEBUG, INFO, WARN, ERROR).

Example `.env` file:

```
PORT=8080
API_KEY=your_api_key_here
LOG_LEVEL=INFO
```

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

### Code of Conduct

Please adhere to the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/0/code_of_conduct/).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
