
# Multi-Stack DevOps Infrastructure Automation Project

## Overview

This project offers practical experience in setting up a multi-component, containerised web application using a "Polyglot Stack" . The aim is to get familiarised with multiple languages, frameworks, and runtime environments, including Node.js, Python, .NET, PostgreSQL, Redis, Docker, Terraform, and Ansible, all deployed on AWS.



## Tech Stack

- **Frontend (Result Viewer)**: Node.js/Express application that displays real-time voting results.
- **Voting Application**: Python/Flask-based frontend for casting votes.
- **Worker Service**: .NET application that processes votes and stores them in a database.
- **Database**: PostgreSQL for persistent storage of vote data.
- **Cache/Queue**: Redis, an in-memory data structure store.
## Project Structure
- `vote/`: Python/Flask application.
- `result/`: Node.js/Express application.
- `worker/`: .NET Worker Service.
- `db/`: PostgreSQL setup.
- `redis/`: Redis configuration.

## Running Applications Locally
### Prerequisites
Ensure you have Python, Node.js, and .NET SDK installed on your machine.
### Steps
1. **Vote App**:
   - Navigate to the `vote/` directory.
   - Install dependencies: `pip install -r requirements.txt`.
   - Run the application: `python app.py`.

2. **Result App**:
   - Navigate to the `result/` directory.
   - Install dependencies: `npm ci` (utilises package-lock.json).
   - Run the application: `node server.js`.

3. **Worker Service**:
   - Navigate to the `worker/` directory.
   - Build and run the service: `dotnet build` followed by `dotnet run`.

4. **Redis**:
   - Install and run Redis locally or use Docker: `redis-server`.

5. **PostgreSQL**:
   - Install PostgreSQL and ensure it is running on the default port (5432).
   - Create a database named `votes` with the necessary credentials.
## Containerisation and Deployment
Each component is designed to be run within Docker containers to simplify deployment and scaling. The provided Dockerfiles in each directory facilitate building and running the respective services. 
## Infrastructure Management
Utilise Terraform for provisioning AWS infrastructure and Ansible for configuring EC2 instances and managing deployments.

## Collaborative Development
- Utilised Jira to manage tasks and workflows.
- Employ Agile methodologies including Scrum practices.
