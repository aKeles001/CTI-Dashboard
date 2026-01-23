
# CTI Dashboard

A  desktop application for Cyber Threat Intelligence (CTI) analysis and visualization, built with Go and modern web technologies using the Wails framework.

## Overview

The CTI Dashboard provides a user interface to help user to interact with the intelligence data. It leverages the power of Go for high-performance backend operations and a React/TypeScript frontend for a responsive user experience.

## Features

*   **Data Visualization:** Interactive charts and graphs for visualizing threat data.
*   **Forum Management:** Easily add, and manage forums.
*   **Asses Severity** Asses severity of the threads and post of scraped forums.

## Technology Stack

*   **Backend:** [Go](https://golang.org/)
*   **Frontend:**
    *   [React](https://reactjs.org/)
    *   [TypeScript](https://www.typescriptlang.org/)
*   **Framework:** [Wails v2](https://wails.io/)
*   **UI Components:** [Shadcn/ui](https://ui.shadcn.com/)




### Prerequisites

First, you need to have Go, Node.js, and the Wails CLI installed on your system.

1.  **Go:** Install Go 1.18 or newer. You can find instructions at [golang.org](https://golang.org/doc/install).

2.  **Node.js:** Install Node.js (which includes npm). You can download it from [nodejs.org](https://nodejs.org/).

3.  **Wails CLI:** Install the Wails CLI by following the instructions on the [official Wails website](https://wails.io/docs/gettingstarted/installation).

### Installation

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/aKeles001/CTI-Dashboard.git
    cd CTI-Dashboard
    ```

2.  **Build the Docker Container:**
    ```sh
    docker-compose build
    ```

### Running the Container


```sh
docker-compose up -d
```

## Building for Production

To build and start the wails desktop application, run the following command:


```sh
docker-compose exec wails-dev bash -c "cd CTI-Dashboard; wails dev -tags webkit2_41"
```