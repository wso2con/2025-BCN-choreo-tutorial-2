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
2. Clone the repository to your local machine
3. Install the Choreo CLI - We will use a pre-release version of the Choreo CLI for this tutorial to demostrate upcoming features.
```
bash <(curl -s https://cli.choreo.dev/install.sh) v1.2.82503121000
```


## Part 1: Platform Engineer's Perspective

#### 1. Sign Up to Choreo
1. Navigate to the [Choreo console](https://choreo.dev/) and sign up using GitHub, Google, Microsoft, or email options.
1. Complete the organization creation process by verifying the emails and OTP received.
1. Select the "Platform Engineer (PE)" perspective from the Choreo Console.

#### 2. Explore Platform Engineering Perspective Overview
1. Review the Data Plane Management section showing data plane status and regional distribution.
1. Examine the pre-configured environments (Development, Production) on the Cloud Data Plane.
1. Inspect the user management dashboard displaying users, groups, and roles.
1. Verify the CD pipeline configurations that define the promotion paths between environments.
1. Review network controls for managing ingress/egress traffic and security boundaries.

#### 3. Create a New Environment
1. Navigate to `Infrastructure → Environments` and review the existing `Development` and `Production` environments.
1. Click + Create Environment and enter the following configuration:
  1. Name: `Staging`
  1. Data Plane: `Choreo Cloud US Dataplane`
  1. DNS Prefix: `staging`
  1. Production Environment: `Checked`
1. Review the auto-generated DNS URL: `{org-uuid}-staging.e1-us-east-azure.choreoapis.dev`
1. Click Create and verify the Staging environment in your environment list.

#### 4. Update the Default CD Pipeline
1. Navigate to `DevOps → CD Pipelines`.
1. Click the edit icon next to the `Default US Deployment Pipeline` pipeline.
1. Click the `+` symbol between `Development` and `Production` environments.
1. Select `Staging` from the environment dropdown.
1. Review the pipeline sequence: `Development → Staging → Production`.
1. Click Update to save your changes.
1. Verify the pipeline visualization now shows the three-environment flow with arrows between each stage.

#### 5. Manage Team Access
1. Navigate to `User Management → Users`.
1. Click `Invite Users` and add user "Joe" with email `joseph@wso2.com`.
1. Assign groups by checking `Admin`, `Project Admin` and `Developer` checkboxes.
1. Click `Invite` to send the invitation email.

#### 6. Enable Environment Promotion Workflow
1. Navigate to `Governance → Workflows`.
1. Find the `Environment Promotion` workflow in the list.
1. Toggle the `Status` switch from `Off` to `On`.
1. Under Roles, check `Admin`, `Project Admin`, and `Choreo Platform Engineer`.
1. Click `Save` to apply the workflow configuration.
1. Verify the workflow status shows Enabled with a green indicator.

#### 7. Provision OpenAI Egress API
1. Navigate to `DB & Services → GenAI Services`.
1. Click `Register Service` and select `OpenAI` as the provider.
1. Configure the service:
    1. Name: `OpenAI`,
     **Note**: Use `OpenAI` name as it is to avoid additional steps in the lab session.
    1. Version: `v1`
1. Get OpenAI API Demo Key
    1. Click on link([Get OpenAI Demo Key](https://308f5c04-e1a0-49e1-817b-e5a8d5422c26-dev.e1-us-east-azure.choreoapis.dev/demo-api-keys/openai-key-service/v1.0/key))
    1. Copy `api_key` value
1. Click Next and add API Key configuration:
    1. Key Name: `OPENAI_API_KEY`
    1. Value: [Value copied from the previous step]
1. Click `Register` to complete setup.
1. Next, Click on `Add to Marketplace` button.

#### 8. Configure Network Policies for Egress Control
1. Navigate to `Governance → Egress Control`.
1. Click `+Create` to configure Egress policies at organizational level.
1. Select `Deny All` as the Egress Control Type.
1. Add an allow rule with:
    1. Rule Name: `Allow Open AI`
    1. Rule: `api.openai.com`
1. Click `Add Rule` and verify the rule appears in the policy list.

## Part 2: Developer's Perspective

Lets develop an application with Choreo. For this tutorial, we will develop a simple webapplication to manage users accounts. Using this web application, user will be able to record their expences by either uploading a receipt image or by entering the expence details manually.

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
    1. Select correct Organization, Repository and Branch
    1. Select project directory as `expense-tracker`
    1. Click `Create`

As nextstep we will create the dependent components for our application. Tipically in an organization when building a new application you would consume existing APIs and databases. For this tutorial we will create the components that we will use in our application.


### 2. Create the Accounts API

1. Create the Accounts API
    1. Go in to `Expense Tracker` project
    1. Select `Create Component`
    1. Select `Service Type`
    1. Select `Authorize with Github`
    1. Go though the github flow and authorize the application to access [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Refresh the repository list and select the `2025-BCN-choreo-tutorial-2` repository.
    1. Select `main` branch
    1. Click Edit in Component Directory and select `expense-tracker/accounts`
    1. Select `go` as the build pack type
    1. Select language version as `1.x`
    1. Check if component name set to `accounts`
    1. Click Create
    1. You will be taken to the build page automatically

### 3. Create the Receipts API

1. Create the Receipts API
    1. Go in to `Expense Tracker` project
    1. Select `Create Component`
    1. Select `Service` Type
    1. Select `Authorize with Github`
    1. Go though the github flow and authorize the application to access [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Refresh the repository list and select the `2025-BCN-choreo-tutorial-2` repository.
    1. Select `main` branch
    1. Click Edit in Component Directory and select `expense-tracker/receipts`
    1. Select `Docker` as the build pack type
    1. Select the docker file path as `receipts/Dockerfile`
    1. Check if component name set to `receipts`
    1. Click Create
    1. You will be taken to the build page automatically

### 4. Connect the Receipts API with OpenAI

  1. Check if you are in `receipts` component
  1. Go to `Dependencies > Connections`
  1. Select `Service` as the connection type
  1. Select `OpenAI` service
  1. And provide the connection name as `Open AI Connection`
  1. Click `Create`
  1. You will be taken to the connection page automatically
  1. You can see a guide on how to add the connection to the component but for this tutorial connection configuration is already done.

### 5. Deploy the Accounts and Receipts APIs to the Dev Environment

1. Deploy Accounts API
    1. Go to `Accounts` component
    1. Click `Deploy`
    1. Click `Configure & Deploy`
    1. Click `Next` still 3rd step
    1. Click `Deploy`
    1. Now the component will be deployed to dev environment

2. Deploy Receipts API
    1. Go to `Receipts` component
    1. Click `Deploy`
    1. Click `Configure & Deploy`
    1. Click `Next` still 3rd step
    1. Click `Deploy`
    1. Now the component will be deployed to dev environment

### 6. Create the BFF API

In this step ideally we should build the BFF API from scratch. For this tutorial we will use the pre-built BFF API fround in `expense-tracker/bffapi` directory. Lets configure it to be used as a component in our project.

We use a configuration file called `component.yaml` to provide aditional component details to Choreo.

1. Create a nodejs express service which will be used as a BFF API
    1. This step is already done. Use the provided `bffapi` directory as a starting point.

1. Configure BFF API Endpoint
    1. Create a new file called `component.yaml` in `expense-tracker/bffapi/.choreo` directory.
    1. Add the following content to the file
        ```
        schemaVersion: 1.1

        endpoints:
          - name: bff-api
            displayName: BFF API Endpoint
            service:
              basePath: /api
              port: 9090
            type: REST
            networkVisibilities:
              - Public
            schemaFilePath: ./openapi.yaml
        ```
    1. Commit and push the changes to the repository

1. Create a new component
    1. Go in to `Expense Tracker` project
    1. Select `Create Component`
    1. Select `Service` Type
    1. Select `Authorize with Github`
    1. Select [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Select `main` branch
    1. Select project directory as `expense-tracker/bffapi`
    1. Select `NodeJs` as the build pack type
    1. Selevct language version as `20`
    1. Click `Create`
    1. You will be taken to the build page automatically
    1. Click on `Auto Build on Commit`

2. Connect to the accounts components
    1. Go to `Dependancies > Connections`
    1. Select `Service` as the connection type
    1. Select `Accounts` service
    1. Provide the connection name as `Accounts`
    1. Click `Create`
    1. You will be taken to the connection page automatically
    1. Copy the connection string to `components.yaml`
    1. Update the source code with the correct environment variables

3. Connect to the receipts components
    1. Go to `Dependancies > Connections`
    1. Select `Service` as the connection type
    1. Select `Receipts` service
    1. Provide the connection name as `Receipts`
    1. Click `Create`
    1. You will be taken to the connection page automatically
    1. Copy the connection string to `components.yaml`
    1. Update the source code with the correct environment variables

4. Push the changes to Choreo
    1. Commit and push the changes to the repository
    1. Choreo will automatically detect the changes and trigger a build.
    1. Deploy the build to dev environment by going to `Deployments` page and clicking `Configure & Deploy` button

### 7. Test locally

1. Testing BFF API locally
    1. Run `choreo login` to login to Choreo from the CLI
    1. Run `choreo change-org` to change to correct organization
    1. Run `choreo connect -c bffapi` to connect to the development environment for the bffapi component.
    1. Retrive the accounts component url from the CLI with `choreo describe component`
    1. Run `curl -X GET "http://<accounts-component-url>/bills"` to test the accounts API
    1. Lets test bffapi locally with the dependacies in dev environment
        1. In the same terminal run `npm run start`
        1. In a new terminal run `curl -X GET "http://<bffapi-component-url>/api/bills"`
        1. You should see the bills from the accounts component

### 8. Create the web application

1. Create a webapp component
    1. Go in to `Expense Tracker` project
    1. Select `Create Component`
    1. Select `Web Application`
    1. Select `Authorize with Github`
    1. Go though the github flow and authorize the application to access [2025-BCN-choreo-tutorial-2](https://github.com/hevayo/2025-BCN-choreo-tutorial-2) repository.
    1. Select project directory as `expense-tracker/webapp`
    1. Select `React` as the build pack type
    1. Enter `npm run build` as the Build Command
    1. Enter Build Path as `/build`
    1. Node version as `20`
    1. Click `Create`
2. Connect to the BFF API with webapp
    1. Disable Authentication for bffapi
        1. Go to deployment page of bffapi
        1. Clieck `Endpoint Configuration`
        1. Uncheck `OAuth2`
        1. Click `Deploy`
    2. Connect to the bffapi with webapp
        1. Go to bffapi component overview page
        2. Go to deployment page of webapp
        3. Click `Configure & Deploy`
        4. Past the bffapi url as apiUrl in `window.configs`
        ```
        window.configs = {
            apiUrl: 'http://<bffapi-component-url>/',
        };
        ```
        1. Disable `Managed Authentication with Choreo`
        5. Click `Deploy`

### 9. Testing the application

1. Testing the application
    1. Go to `Webapp` component
    2. Click the public endpoint to open the application


## Part 3: Operational Excellence
#### 1. Promote the Application to higher environment
1. Upon receiving the promotion requests, approve them
1. Promote to production in the following order
1. Accounts, BillParser, Bff API and Web Application

#### 2. Monitor System Performance & Observability
1. Navigate to `Insights → Operational`.
  1. Set Environment: `Development`.
  1. Set Time Period: `1D`.
  1. Analyze CPU & Memory Usage across components.
1. Navigate to `Observability → Metrics`.
  1. Set Environment: `Development`.
  1. Hover over component links to view request counts, latencies, and error percentages.
  1. Check deployment logs for further insights.

#### 3. Configure Latency Alert Rule
1. Navigate to `Observability → Alerts`.
1. Click `Create Alert Rule` to add a new alert rule.
  1. Select Alert Type: `Latency`.
  1. Choose Metric: `99th Percentile`.
  1. Set Environment: `Development`.
  1. Select Deployment Track: `master`.
  1. Define Threshold (ms): `1000`.
  1. Add Emails for notifications: `manjular@wso2.com`.
1. Click `Create` to activate the alert.