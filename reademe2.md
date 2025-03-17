# Building Applications with Choreo: An Internal Developer Platform Lab Session

## How an Internal Developer Platform Lets Developers Focus on Code

This hands-on lab session demonstrates how Choreo functions as an Internal Developer Platform (IDP), creating a clear separation of concerns between platform engineers and application developers. By implementing this separation, organizations can enable developers to focus entirely on writing code while platform engineers handle infrastructure, security, and governance concerns.

## Lab Outline

### Part 1: Platform Engineering Setup
Platform Engineers establish the foundational elements necessary for application development, including environments, pipelines, databases, and security policies.

### Part 2: Developer Experience
Developers leverage these pre-configured resources to build, test, and deploy applications with minimal infrastructure concerns.

### Part 3: Operational Excellence
Platform Engineers implement monitoring, alerting, and observability to ensure application reliability and performance.

### Part 4: Recap and Path Forward
Summary of achievements, benefits realized, and exploration of advanced capabilities for future implementation.

---
## Prerequisites
1. Fork the [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository


<details>
<summary><h2>Part 1: Platform Engineer's Perspective</h2></summary>

#### 1. Sign Up to Choreo
- Navigate to the [Choreo console](https://console.choreo.dev/) and sign up using GitHub, Google, Microsoft, or email options.
- Complete the organization creation process by verifying the emails and OTP received.
- Select the "Platform Engineer (PE)" perspective from the Choreo Console.

#### 2. Explore Platform Engineering Perspective Overview
- Review the Data Plane Management section showing data plane status and regional distribution.
- Examine the pre-configured environments (Development, Production) with their specific runtime characteristics.
- Inspect the operational metrics dashboard showing resource utilization, deployment frequency, and system health.
- Verify the CD pipeline configurations that define the promotion paths between environments.
- Review network controls for managing ingress/egress traffic and security boundaries.

#### 3. Create a New Environment
- Navigate to Infrastructure → Environments and review the existing Development and Production environments.
- Click + Create Environment and enter the following configuration:
  - Name: "Staging"
  - Data Plane: "Choreo Cloud US Dataplane"
  - DNS Prefix: "staging"
  - Production Environment: Unchecked
- Review the auto-generated DNS URL: "staging.{org-name}.choreoapis.dev".
- Click Create and verify the Staging environment appears with "Active" status in your environment list.

#### 4. Update the Default CD Pipeline
- Navigate to DevOps → CD Pipelines and select the "Default US Deployment Pipeline".
- Click the edit (✏️) icon next to the pipeline.
- Click the "+" symbol between Development and Production environments.
- Select "Staging" from the environment dropdown.
- Review the pipeline sequence: Development → Staging → Production.
- Click Update to save your changes.
- Verify the pipeline visualization now shows the three-environment flow with arrows between each stage.

#### 5. Manage Team Access
- Navigate to Project Settings → Team.
- Click "Invite Member" and add user "Joe" with email "joe@example.com".
- Assign roles by checking "Project Admin" and "Developer" checkboxes.
- Set environment access permissions for Development, Staging, and Production.
- Click Send to issue the invitation.

#### 6. Enable Environment Promotion Workflow
- Navigate to Governance → Workflows.
- Find the "Environment Promotion" workflow in the list.
- Toggle the Status switch from Off to On.
- Under Roles, check "Admin", "Project Admin", and "Choreo DevOps".
- Click Save to apply the workflow configuration.
- Verify the workflow status shows Enabled with a green indicator.

#### 7. Provision Resources
- Navigate to DB & Services → Databases.
- Click + Create and select "PostgreSQL" as the database type.
- Configure the database:
  - Service name: "customer-portal-db"
  - Cloud provider: "Azure" or "AWS"
  - Region: "US East" or region closest to your location
  - Service plan: "Hobbyist" (1vCPU, 2GB RAM)
- Click Create and wait for status to change from "Creating" to "Active".
- Copy connection parameters:
  - Host: customer-portal-db.postgres.database.azure.com
  - Port: 5432
  - Default User: postgres
  - Default Database: postgres
  - Password: [auto-generated password]

#### 8. Configure Network Policies for Egress Control
- Navigate to Governance → Egress Control.
- Click + Create to establish a new policy.
- Select "Deny All" as the default rule.
- Add an allow rule with:
  - Name: "DB-Access"
  - Type: "Egress"
  - Target: "customer-portal-db.postgres.database.azure.com"
  - Port: 5432
  - Protocol: TCP
- Click Add Rule and verify the rule appears in the policy list.
- Confirm the policy status shows "Active" and test connectivity to confirm only database traffic is allowed.


</details>

<details>
<summary><h2>Part 2: Developer's Perspective</h2></summary>

Lets develop an application with Choreo. For this tutorial, we will develop a simple webapplication to manage users accounts. Using this web application, user will be able to record their expences by either uploading a receipt image or by entering the expence details manually.

#### User Interface


#### Architecture

The application follows a microservices architecture with the following components:

![Application Architecture](./docs/images/architecture.png)

1. **Web Application (Webapp)**: React-based frontend application that provides the user interface for expense management.

2. **API Gateway (API GW)**: Entry point for all client-side requests that handles routing and authentication.

3. **Backend for Frontend (BFF)**: Orchestration layer that optimizes and aggregates backend service calls for the frontend.

4. **Accounts Service**: Core service responsible for managing user accounts and expense records.

5. **Receipt Service**: Specialized service that processes receipt images and extracts data using OpenAI integration.

6. **Egress Gateway**: Security control point that manages and secures all outbound traffic to external services.

All components are deployed and managed through Choreo, ensuring secure communication and monitoring.


### 1. Create a new project

Choreo project allows you to group related components together. It also creates a newtwok boundry around the components and allows you to manage incoming and outgoing traffic.

1. Create a new project
    1. project display name: `Expense Tracker`
    1. Project name: `expense-tracker`
    1. Project description: `Personal Expense Tracking Application`
    1. Select `Authorize with Github`
    1. Go though the github flow and authorize the application to access [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Select project directory as `expense-tracker`
    1. Click `Create`


### 2. Creating dependent components

At this step we will create the dependent components for our application. Tipically in an organization when building a new application you would consume existing APIs and databases. For this tutorial we will create the components that we will use in our application.

1. Create Accounts API
    1. Go in to `Choreo-Tutorial-2` project
    1. Select `Create Component`
    1. Select `Service Type`
    1. Select `Authorize with Github`
    1. Go though the github flow and authorize the application to access [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Refresh the repository list and select the `2025-BCN-choreo-tutorial-2` repository.
    1. Select `main` branch
    1. Click Edit in Component Directory and select `accounts`
    1. Select `go` as the build pack type
    1. Select language version as `1.x`
    1. Endpoint details configuration
        1. Enter port as `8080`
        1. API Type as `REST` 
        1. Base Path as `/`
        1. Click Edit under Schema Path and select `accounts/openapi.yaml`
    1. Click Create
1. Test Accounts API

1. Create Receipts API
    1. Go in to `Choreo-Tutorial-2` project
    1. Select `Create Component`
    1. Select `Service Type`
    1. Select `Authorize with Github`
    1. Go though the github flow and authorize the application to access [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Refresh the repository list and select the `2025-BCN-choreo-tutorial-2` repository.
    1. Select `main` branch
    1. Click Edit in Component Directory and select `accounts`
    1. Select `go` as the build pack type
    1. Select language version as `1.x`
    1. Endpoint details configuration
        1. Enter port as `8080`
        1. API Type as `REST` 
        1. Base Path as `/`
        1. Click Edit under Schema Path and select `accounts/openapi.yaml`
    1. Click Create

1. Test Receipts API
    1.  
    1. Use the following curl command to test the API replace the `{account_id}` with the account id you created in the previous step
        ```
        curl -X POST http://receipts.staging.choreoapis.dev/receipts \
        -H "Content-Type: application/json" \
        -d '{"name": "John Doe", "email": "john.doe@example.com"}'
        ```

