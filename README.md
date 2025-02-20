# BeKind, Rewind
BeKind, Rewind is a draft project for a simple web frontend UI of an offline Movie Tracker app I plan to build in the near future. This project serves as a playground for exploring [HTMX](https://htmx.org/) and Go.

Currently, the app has no backend, and movies are stored in-memory with non-persistent storage. The few movies displayed when you run the app are hardcoded. However, you can edit, delete, and add new movies to the list.

## Features
- View a hardcoded list of movies.
- Edit and delete movies.
- Add new movies to the list.

## Prerequisites
Before running the app, ensure that you have Go version 1.22 or higher installed on your system.

## Installation and running the app

### Install Go 
First, install Go if you haven't already. You can download Go from the official website: https://golang.org/dl/.

#### Building from the source

##### Clone the Repository
Clone this repository to your local machine:

```bash
git clone https://github.com/markort147/bekind.git
```

##### Build the App
Navigate to the directory where you cloned the repo, then run the following command:
```bash
go build ./cmd/bekindrewind/ -o bekindrewind
```

##### Run the app
Run the app:
```bash
./bekindrewind -config <path-to-your-config-yaml>
```

#### As a Go command

##### Install via Go
```bash
go install github.com/markort147/bekind/cmd/bekindrewind@latest
```

##### Run the app
Run the app:
```bash
~/go/bin/bekindrewind -config <path-to-your-config-yaml>
```
### Access the Web App
Open your browser and go to http://localhost:[PORT] to view the app in action.

## Development with Air
For a smoother development experience, I recommend using Air, which provides live reloading. The repository includes a .air.toml configuration file for easy integration with [Air](https://github.com/air-verse/air).

To use Air:
1. Install Air by following the instructions on the official [Air GitHub page](https://github.com/air-verse/air).
2. Run the app with Air:
```bash
air .
```
This will automatically reload the application whenever changes are made to the code, enhancing your development workflow.
