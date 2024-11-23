# Golang Project Setup

This is Application For Test Software Developer Lead From Bosnet Distribution Indonesia.

## Requirements

- **Golang**: Version 1.23.2
- **Operating System**: Debian Based
- **Docker**:
  - Docker Desktop (Windows)
  - Docker (Linux/Mac)

## Getting Started

### Step 1: Configure Environment Variables

- Copy the example environment file:

```bash
cp .env.example .env
```

- Open the `.env` file and fill in your configuration

### Step 2: Deploy the Application

To deploy your application, run the appropriate script based on your operating system:

- Windows

```cmd
C:\Users\Your_User\Document\Project> .\deploy.bat
```

- Linux / Mac

```bash
$ bash deploy.sh
```

### Step 3: Choose Your Environment

During the deployment process, you will be prompted to choose between Production or Development modes. Select the one that suits your needs.

Press Y or y for **Production** Deploy, otherwise is a **Development** Deploy

### Step 4: Access API

Accessing Your App API via:

```curl
http://localhost:{PORT}
```

_{PORT}_ depends on your .env parameter **OUT_PORT**
